dl(Debug x Log) - The instant logger package for debug
======

[![Go Reference](https://pkg.go.dev/badge/github.com/task4233/dl.svg)](https://pkg.go.dev/github.com/task4233/dl)
[![.github/workflows/ci.yml](https://github.com/task4233/dl/actions/workflows/ci.yml/badge.svg)](https://github.com/task4233/dl/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/task4233/dl)](https://goreportcard.com/report/github.com/task4233/dl)
[![codecov](https://codecov.io/gh/task4233/delog/branch/main/graph/badge.svg?token=93KXZTJGGL)](https://codecov.io/gh/task4233/delog)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

## Description

dl provides a package for logging instantly such as dubugging and one command to remove them all.

## Installation
### Go1.18

```bash
$ go install github.com/task4233/dl/cmd/dl@latest
```

## Use Case
### Debug

[Playground](https://go.dev/play/p/DW6BBg2Wd9a)
```go
package main

import "github.com/task4233/dl"

type T []int

func (t T) append(v int) {
	dl.Printf("Type: %T, v: %v\n", v, v) // This statement can be removed by `$ dl clean main.go`
	t = append(t, v)
}

func (t T) change(v int) {
	dl.Printf("Type: %T, v: %v\n", v, v) // This statement can be removed by `$ dl clean main.go`
	t[0] = v
}

func main() {
	var t T = []int{1, 3}
	t.append(5)
	t.change(5)
}

// Output:
// Type: int, v: 5
// Type: int, v: 5
```

### Adds dl into pre-commit of Git
1. Please run commands below to install dl in your Git repository.

```bash
$ cat > ./.git/hooks/pre-commit << EOF
#!/bin/sh
$(echo $GOBIN)/dl clean .
git add .
EOF
$ chmod +x ./.git/hooks/pre-commit
```

2. Just commit

```bash
$ cat main.go 
package main

import (
	"fmt"
	
	"github.com/task4233/dl"
)

func SayHi[T any](v T) {
	dl.Printf("Type: %T, v: %v\n", v, v) // This statement can be removed by `$ dl clean main.go`
	fmt.Println("Hi, ", v)
}

func main() {
    SayHi("hoge")
}
$ go run main.go
Type: string, v: hoge
Hi,  hoge
$ git add main.go
$ git commit -m "feat: add main.go"
remove dl from main.go # automatically removed
[master 975ecf9] feat: add main.go
 1 file changed, 12 insertions(+), 21 deletions(-)
 rewrite main.go (91%)
$ git diff HEAD^
diff --git a/main.go b/main.go
index 90a78bd..0e28e8a 100644
--- a/main.go
+++ b/main.go
@@ -1,21 +1,12 @@
 package main

+import (
+       "fmt"
+)
 
+func SayHi[T any](v T) {
+       fmt.Println("Hi, ", v)
+}

 func main() {
+       SayHi("hoge")
 }
```

## Author

[task4233](https://task4233.dev)
