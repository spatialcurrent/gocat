[![CircleCI](https://circleci.com/gh/spatialcurrent/gocat/tree/main.svg?style=svg)](https://circleci.com/gh/spatialcurrent/gocat/tree/main)
[![Go Report Card](https://goreportcard.com/badge/spatialcurrent/gocat?style=flat-square)](https://goreportcard.com/report/github.com/spatialcurrent/gocat)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/spatialcurrent/gocat)](https://pkg.go.dev/github.com/spatialcurrent/gocat)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/gocat/blob/master/LICENSE)

# gocat

# Description

**gocat** is a super simple command line program for concatenating files.  **gocat** lazily loads files to reduce memory usage and the number of active file pointers.  **gocat** can read from local files, HTTP(S) endpoints, and files on S3.

## Platforms

The following platforms are supported.  Pull requests to support other platforms are welcome!

| GOOS | 386 | amd64 | arm | arm64 |
| ---- | --- | ----- | --- | ----- |
| darwin | - | ✓ | - | - |
| freebsd | ✓ | ✓ | ✓ | - |
| linux | ✓ | ✓ | ✓ | ✓ |
| openbsd | ✓ | ✓ | - | - |
| solaris | - | ✓ | - | - |
| windows | ✓ | ✓ | - | - |

## Releases

Find releases for the supported platforms at [https://github.com/spatialcurrent/gocat/releases](https://github.com/spatialcurrent/gocat/releases).  See the **Building** section below to build for another platform from source.

## Usage

See the usage below or the following examples.

```shell
gocat is a super simple utility to concatenate files (local, remote, or on AWS S3) provided as positional arguments.
Supports stdin (aka "-"), local files (path/to/file or file://path/to/file), remote files (http://path/to/file), or files on AWS S3 (s3://path/to/file).
Supports the following compression algorithms: bzip2, flate, gzip, none, snappy, zip, zlib.

Usage:
  gocat [flags] [-|stdin|FILE|URI]...

Flags:
  -a, --append-new-lines               append new lines to files
      --aws-access-key-id string       AWS Access Key ID
      --aws-default-region string      AWS Default Region
      --aws-profile string             AWS Profile
      --aws-region string              AWS Region (overrides default region)
      --aws-secret-access-key string   AWS Secret Access Key
      --aws-session-token string       AWS Session Token
  -b, --buffer-size int                buffer size for file reader (default 4096)
  -h, --help                           help for gocat
      --version                        show version
```

## Examples

### Shell

```shell
gocat "~/.bash_aliases" "~/.bashrc"
```

### Remote File

`gocat` arguments can point to remote files

```shell
gocat https://spatialcurrent.io
```

### All Files With Extension

`gocat` will print every file given as a positional argument, which can be useful when looking through thousands of files.

```shell
gocat $( find . -print | grep -i '.*[.]md')
```

## Building

**gocat** is written in pure Go.  The only dependency needed to compile the program is [Go](https://go.dev/).

This project supports the use of [direnv](https://direnv.net/) to manage environment variables.

If using `macOS`, follow the `macOS` instructions below.

Use `make bin/gocat` to build for your local operating system.  Use `make build_release` to build for a release.  Alternatively, you can call `go build` directly for your specific scenario.

### macOS

To install `go` on macOS with homebrew use `brew install go`.

To install `direnv` on macOS with homebrew use `brew install direnv`.  If using bash, then add `eval \"$(direnv hook bash)\"` to the `~/.bash_profile` file .  If using zsh, then add `eval \"$(direnv hook zsh)\"` to the `~/.zshrc` file.

## Testing

**CLI**

To run CLI testes use `make test_cli`, which uses [shUnit2](https://github.com/kward/shunit2).  If you recive a `shunit2:FATAL Please declare TMPDIR with path on partition with exec permission.` error, you can modify the `TMPDIR` environment variable in line or with `export TMPDIR=<YOUR TEMP DIRECTORY HERE>`. For example:

```
TMPDIR="/usr/local/tmp" make test_cli
```

**Go**

To run Go tests use `make test_go` or (`bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

## Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/gocat/blob/master/CONTRIBUTING.md) for how to get started.

## License

This work is distributed under the **MIT License**.  See **LICENSE** file.
