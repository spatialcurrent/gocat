[![CircleCI](https://circleci.com/gh/spatialcurrent/gocat/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/gocat/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/gocat)](https://goreportcard.com/report/spatialcurrent/gocat)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/gocat?status.svg)](https://godoc.org/github.com/spatialcurrent/gocat) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/gocat/blob/master/LICENSE)

# gocat

# Description

**gocat** is a super simple command line program for concatenating files.  **gocat** supports the following operating systems and architectures.

| GOOS | GOARCH |
| ---- | ------ |
| darwin | amd64 |
| linux | amd64 |
| windows | amd64 |
| linux | arm64 |

# Installation

No installation is required.  Just grab a [release](https://github.com/spatialcurrent/gocat/releases).  You might want to rename your binary to just `gocat` (or `cat`) for convenience.

If you do have go already installed, you can just run using `go run main.go` or install with `bash scripts/install.sh`

# Usage

See the few examples below.

# Examples

**Shell**

```shell
gocat "~/.bash_aliases" "~/.bashrc"
```

**Remote File**

`gocat` arguments can point to remote files

```shell
gocat https://spatialcurrent.io
```

**All Files With Extension**

`gocat` will print every file given as a positional argument, which can be useful when looking through thousands of files.

```shell
gocat $( find . -print | grep -i '.*[.]md')
```

# Building

You can build all the released artifacts using `make build` or run the make target for a specific operating system and architecture.

# Testing

Run test using `bash scripts/test.sh` or `make test`, which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [ineffassign](https://github.com/gordonklaus/ineffassign), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/gocat/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
