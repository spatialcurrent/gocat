// =================================================================
//
// Copyright (C) 2022 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/spatialcurrent/go-lazy/pkg/lazy"
	"github.com/spatialcurrent/go-reader-writer/pkg/grw"
	"github.com/spatialcurrent/go-reader-writer/pkg/splitter"
)

const (
	GoCatVersion = "v2.1.0"
)

const (
	flagAWSProfile         = "aws-profile"
	flagAWSDefaultRegion   = "aws-default-region"
	flagAWSRegion          = "aws-region"
	flagAWSAccessKeyID     = "aws-access-key-id"
	flagAWSSecretAccessKey = "aws-secret-access-key"
	flagAWSSessionToken    = "aws-session-token"

	flagBufferSize     = "buffer-size"
	flagAppendNewlines = "append-new-lines"

	flagVersion = "version"
)

func initFlags(flag *pflag.FlagSet) {
	flag.String(flagAWSProfile, "", "AWS Profile")
	flag.String(flagAWSDefaultRegion, "", "AWS Default Region")
	flag.StringP(flagAWSRegion, "", "", "AWS Region (overrides default region)")
	flag.StringP(flagAWSAccessKeyID, "", "", "AWS Access Key ID")
	flag.StringP(flagAWSSecretAccessKey, "", "", "AWS Secret Access Key")
	flag.StringP(flagAWSSessionToken, "", "", "AWS Session Token")
	flag.IntP(flagBufferSize, "b", 4096, "buffer size for file reader")
	flag.BoolP(flagAppendNewlines, "a", false, "append new lines to files")
	flag.Bool(flagVersion, false, "show version")
}

func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()
	err := v.BindPFlags(cmd.Flags())
	if err != nil {
		return v, fmt.Errorf("error binding flag set to viper: %w", err)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv() // set environment variables to overwrite config
	return v, nil
}

func checkConfig(v *viper.Viper) error {
	bufferSize := v.GetInt(flagBufferSize)
	if bufferSize <= 0 {
		return fmt.Errorf("buffer size must be greater than 0")
	}
	return nil
}

func initAWSSession(accessKeyID string, secretAccessKey string, sessionToken string, region string) (*awssession.Session, error) {
	config := aws.Config{
		MaxRetries: aws.Int(3),
		Region:     aws.String(region),
	}

	if len(accessKeyID) > 0 && len(secretAccessKey) > 0 {
		config.Credentials = credentials.NewStaticCredentials(
			accessKeyID,
			secretAccessKey,
			sessionToken)
	}

	s, err := awssession.NewSessionWithOptions(awssession.Options{
		Config: config,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating new AWS session for region %q: %w", region, err)
	}
	return s, nil
}

func main() {
	cmd := &cobra.Command{
		Use:                   `gocat [flags] [-|stdin|FILE|URI]...`,
		DisableFlagsInUseLine: true,
		Short:                 "gocat is a super simple utility to concatenate files (local, remote, or on AWS S3) provided as positional arguments.",
		Long: `gocat is a super simple utility to concatenate files (local, remote, or on AWS S3) provided as positional arguments.
Supports stdin (aka "-"), local files (path/to/file or file://path/to/file), remote files (http://path/to/file), or files on AWS S3 (s3://path/to/file).
Supports the following compression algorithms: ` + strings.Join(grw.Algorithms, ", ") + `.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return fmt.Errorf("error initializing viper: %w", err)
			}

			if v.GetBool(flagVersion) {
				fmt.Println(GoCatVersion)
				return nil
			}

			if len(args) == 0 {
				return cmd.Usage()
			}

			if errConfig := checkConfig(v); errConfig != nil {
				return errConfig
			}

			bufferSize := v.GetInt(flagBufferSize)
			appendNewlines := v.GetBool(flagAppendNewlines)

			// if reading from AWS, use this region to create a session
			region := v.GetString(flagAWSRegion)
			if len(region) == 0 {
				if defaultRegion := v.GetString(flagAWSDefaultRegion); len(defaultRegion) > 0 {
					region = defaultRegion
				}
			}

			accessKeyID := v.GetString(flagAWSAccessKeyID)
			secretAccessKey := v.GetString(flagAWSSecretAccessKey)
			sessionToken := v.GetString(flagAWSSessionToken)

			sessions := map[string]*awssession.Session{}
			clients := map[string]*s3.S3{}

			ctx := context.Background()

			inputReaders := make([]io.Reader, 0)
			for _, uri := range args {
				uri := uri

				if uri == "-" {
					uri = "stdin"
				}

				if strings.HasPrefix(uri, "s3://") {
					if _, ok := sessions[region]; !ok {
						s, initAWSSessionError := initAWSSession(accessKeyID, secretAccessKey, sessionToken, region)
						if initAWSSessionError != nil {
							return fmt.Errorf("error creating new AWS session for uri %q: %w", uri, initAWSSessionError)
						}
						sessions[region] = s
					}
				}

				if uri == "stdin" {
					inputReaders = append(inputReaders, os.Stdin)
				} else {
					inputReaders = append(inputReaders, lazy.NewLazyReader(func() (io.Reader, error) {
						var bucketClient *s3.S3
						if strings.HasPrefix(uri, "s3://") {
							_, fullpath := splitter.SplitURI(uri)
							bucketName := fullpath[0:strings.Index(fullpath, "/")]
							bucketRegion, getBucketRegionError := s3manager.GetBucketRegion(ctx, sessions[region], bucketName, region)
							if getBucketRegionError != nil {
								if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
									return nil, fmt.Errorf("error getting region for uri %q: %w", uri, aerr)
								}
								return nil, fmt.Errorf("error getting region for uri %q: %w", uri, getBucketRegionError)
							}
							if _, ok := sessions[bucketRegion]; !ok {
								sessions[bucketRegion] = sessions[region].Copy(&aws.Config{Region: aws.String(bucketRegion)})
							}
							if _, ok := clients[bucketRegion]; !ok {
								clients[bucketRegion] = s3.New(sessions[bucketRegion])
							}
							bucketClient = clients[bucketRegion]
						}
						r, err := grw.ReadFromResource(&grw.ReadFromResourceInput{
							URI:        uri,
							Alg:        "none",
							Dict:       grw.NoDict,
							BufferSize: bufferSize,
							S3Client:   bucketClient,
						})
						if err != nil {
							return nil, fmt.Errorf("error reading from URI %q: %w", uri, err)
						}
						return r.Reader, nil
					}))
				}
				// if appendNewlines is true, then append a new line character at the end of the file
				if appendNewlines {
					inputReaders = append(inputReaders, bytes.NewReader([]byte("\n")))
				}
			}

			if _, err := io.Copy(os.Stdout, io.MultiReader(inputReaders...)); err != nil {
				return fmt.Errorf("error copying input: %w", err)
			}

			return nil
		},
	}
	initFlags(cmd.Flags())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "gocat: "+err.Error())
		fmt.Fprintln(os.Stderr, "Try gocat --help for more information.")
		os.Exit(1)
	}
}
