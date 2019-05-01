// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package main

import (
	"fmt"
	"os"
	"strings"
)

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

import (
	"github.com/aws/aws-sdk-go/aws"
	//"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

import (
	"github.com/spatialcurrent/go-reader-writer/grw"
)

const (
	flagAWSProfile         string = "aws-profile"
	flagAWSDefaultRegion   string = "aws-default-region"
	flagAWSRegion          string = "aws-region"
	flagAWSAccessKeyID     string = "aws-access-key-id"
	flagAWSSecretAccessKey string = "aws-secret-access-key"
	flagAWSSessionToken    string = "aws-session-token"

	flagBufferSize     string = "buffer-size"
	flagAppendNewlines string = "append-new-lines"
)

func initFlags(flag *pflag.FlagSet) {
	flag.String(flagAWSProfile, "", "AWS Profile")
	flag.String(flagAWSDefaultRegion, "", "AWS Default Region")
	flag.StringP(flagAWSRegion, "", "", "AWS Region (overrides default region)")
	flag.StringP(flagAWSAccessKeyID, "", "", "AWS Access Key ID")
	flag.StringP(flagAWSSecretAccessKey, "", "", "AWS Secret Access Key")
	flag.StringP(flagAWSSessionToken, "", "", "AWS Session Token")
	flag.IntP(flagBufferSize, "b", 4096, "buffer size for file reader")
	flag.BoolP(flagAppendNewlines, "a", false, "append new lines to files that do not end in new lines characters")
}

func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()
	err := v.BindPFlags(cmd.Flags())
	if err != nil {
		return v, errors.Wrap(err, "error binding flag set to viper")
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

func main() {
	cmd := &cobra.Command{
		Use:   "gocat [-|stdin|FILE|URI]...",
		Short: "gocat",
		Long:  `gocat is a super simple utility to concatenate files (local, remote, or on AWS S3) provided as positional arguments.  Supports stdin (aka "-"), local files (path/to/file or file://path/to/file), remote files (http://path/to/file), or files on AWS S3 (s3://path/to/file).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return errors.Wrap(err, "error initializing viper")
			}

			if len(args) == 0 {
				return cmd.Usage()
			}

			if errConfig := checkConfig(v); errConfig != nil {
				return errConfig
			}

			bufferSize := v.GetInt(flagBufferSize)

			stdinBytes := make([]byte, 0)

			var session *awssession.Session

			var s3Client *s3.S3

			for _, uri := range args {

				if uri == "-" {
					uri = "stdin"
				}

				if strings.HasPrefix(uri, "s3://") {
					if session == nil {
						accessKeyID := v.GetString(flagAWSAccessKeyID)
						secretAccessKey := v.GetString(flagAWSSecretAccessKey)
						sessionToken := v.GetString(flagAWSSessionToken)

						region := v.GetString(flagAWSRegion)
						if len(region) == 0 {
							if defaultRegion := v.GetString(flagAWSDefaultRegion); len(defaultRegion) > 0 {
								region = defaultRegion
							}
						}

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

						session = awssession.Must(awssession.NewSessionWithOptions(awssession.Options{
							Config: config,
						}))
					}
					if s3Client == nil {
						s3Client = s3.New(session)
					}
				}

				inputBytes := make([]byte, 0)

				// if not reading from stdin or stdin hasn't been read yet.
				if uri != "stdin" || len(stdinBytes) == 0 {
					inputReader, _, inputError := grw.ReadFromResource(uri, "none", bufferSize, false, s3Client)
					if inputError != nil {
						return errors.Wrap(inputError, fmt.Sprintf("error reading from uri %q", uri))
					}
					b, inputError := inputReader.ReadAllAndClose()
					if inputError != nil {
						return errors.Wrap(inputError, fmt.Sprintf("error reading from uri %q", uri))
					}
					inputBytes = b
				}

				// if reading from stdin and stdin bytes are already cached
				if uri == "stdin" && len(inputBytes) == 0 && len(stdinBytes) > 0 {
					inputBytes = stdinBytes
				}

				if len(inputBytes) > 0 {

					_, err = os.Stdout.Write(inputBytes)
					if err != nil {
						return errors.Wrap(err, fmt.Sprintf("error writing bytes from uri %q", uri))
					}

					if v.GetBool(flagAppendNewlines) {
						if inputBytes[len(inputBytes)-1] != '\n' {
							_, err := os.Stdout.Write([]byte("\n"))
							if err != nil {
								return errors.Wrap(err, fmt.Sprintf("error writing new line from uri %q", uri))
							}
						}
					}

					if uri == "stdin" {
						stdinBytes = inputBytes
					}

				}

			}
			return nil
		},
	}
	initFlags(cmd.Flags())

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
