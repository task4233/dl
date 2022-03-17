dl(Debug x Log) - The instant logger package for debug
======

[![.github/workflows/ci.yaml](https://github.com/task4233/dl/actions/workflows/ci.yaml/badge.svg)](https://github.com/task4233/dl/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/task4233/dl)](https://goreportcard.com/report/github.com/task4233/dl)
[![codecov](https://codecov.io/gh/task4233/delog/branch/main/graph/badge.svg?token=93KXZTJGGL)](https://codecov.io/gh/task4233/delog)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

## Description

dl provides a package for logging instantly such as dubugging and one command to remove them all.

## Installation
### Go

```bash
$ go install github.com/task4233/dl/cmd/dl@latest
```

## Use Case
- Writes dl function on Go codes for debugging and sweeps them all with `$ dl clean .`
- Adds dl into pre-commit of Git

## Author

[task4233](https://task4233.dev)
