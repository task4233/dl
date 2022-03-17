Delog(debug x log) - The instant logger package for debug
======

[![.github/workflows/ci.yml](https://github.com/task4233/delog/actions/workflows/ci.yml/badge.svg)](https://github.com/task4233/delog/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/task4233/delog)](https://goreportcard.com/report/github.com/task4233/delog)
[![codecov](https://codecov.io/gh/task4233/delog/branch/main/graph/badge.svg?token=93KXZTJGGL)](https://codecov.io/gh/task4233/delog)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

## Description

Delog provides a package for logging instantly such as dubugging and one command to remove them all.

## Installation
### Go

```bash
$ go install github.com/task4233/delog/cmd/delog@latest
```

## Use Case
- Writes delog function on Go codes for debugging and sweeps them all with `$ delog clean .`
- Adds delog into pre-commit of Git

## Author

[task4233](https://task4233.dev)
