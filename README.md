# *NOTICE*
**If you use Go 1.18 or later, Please use [github.com/task4233/dl/v2](https://pkg.go.dev/github.com/task4233/dl/v2), not github.com/task4233/dl**.

dl - The logger not committed to Git for debug
======

[![Go Reference](https://pkg.go.dev/badge/github.com/task4233/dl.svg)](https://pkg.go.dev/github.com/task4233/dl)
[![.github/workflows/ci.yml](https://github.com/task4233/dl/actions/workflows/ci.yml/badge.svg)](https://github.com/task4233/dl/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/task4233/dl)](https://goreportcard.com/report/github.com/task4233/dl)
[![codecov](https://codecov.io/gh/task4233/dl/branch/main/graph/badge.svg?token=xrhysp4Tzf)](https://codecov.io/gh/task4233/dl)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)


![delog.gif](https://user-images.githubusercontent.com/29667656/159164178-f72aede7-f825-438a-add6-aa3deedf8c4c.gif)

## Description
Who doesn't write wrong codes? No one.  
Then, programs don't work well, and developers write logs for debug to understand what happens.

However, some developers forget to delete their logs after resolving the problem and push their codes. In the worse case, the logs might be released.

**dl** is developed to resolve their problems.

## Features
- **Logging package for debug in Go**
  - [`dl` provides wrapping function for **logr.Logger**](https://pkg.go.dev/github.com/task4233/dl#NewLogger).
- **Command for parallel Sweeping all functions of this package**
- **Command for installing git hooks and .gitignore**

## Installation
### Go 1.17 or earlier
It doesn't contain a generics feature.

```bash
$ go install github.com/task4233/dl/cmd/dl@v1
```

### Go 1.18

```bash
$ go install github.com/task4233/dl/cmd/dl@main
```

## Usage

1. debug your codes with `dl` package

[Playground](https://go.dev/play/p/PWJAD4tvLLS)
```go
package main

import (
	"github.com/task4233/dl"
)

type T []int

func (t T) append(v int) {
	t = append(t, v)
	dl.Printf("t: %#v, v: %d\n", t, v)
}

func (t T) change(v int) {
	t[0] = v
	dl.Printf("t: %#v, v: %d\n", t, v)
}

func main() {
	var t T = []int{1, 3}
	t.append(5)
	t.change(5)
}

// Output:
// t: main.T{1, 3, 5}, v: 5
// t: main.T{5, 3}, v: 5
```

2. Install dl

```bash
$ dl init .
```

3. Just commit

- `delog` is used in the file.

```bash
$ cat main.go 
package main

import (
	"github.com/task4233/dl"
)

type T []int

func (t T) append(v int) {
	t = append(t, v)
	dl.Printf("t: %#v, v: %d\n", t, v)
}

func (t T) change(v int) {
	t[0] = v
	dl.Printf("t: %#v, v: %d\n", t, v)
}

func main() {
	var t T = []int{1, 3}
	t.append(5)
	t.change(5)
}

```

- invoke `$ git commit`

```bash
$ git add main.go
$ git commit -m "feat: add main.go"
remove dl from main.go # automatically removed
[master 975ecf9] feat: add main.go
 1 file changed, 12 insertions(+), 21 deletions(-)
 rewrite main.go (91%)
```

- `delog` is removed automatically

```bash
$ git show HEAD
index 90a78bd..0e28e8a 100644
--- a/main.go
+++ b/main.go
@@ -1,21 +1,12 @@
package main

type T []int

func (t T) append(v int) {
	t = append(t, v)

}

func (t T) change(v int) {
	t[0] = v

}

func main() {
	var t T = []int{1, 3}
	t.append(5)
	t.change(5)
}
```

- removed `delog` codes are restored(not commited)

```bash
$ cat main.go 
package main

import (
	"github.com/task4233/dl"
)

type T []int

func (t T) append(v int) {
	t = append(t, v)
	dl.Printf("t: %#v, v: %d\n", t, v)
}

func (t T) change(v int) {
	t[0] = v
	dl.Printf("t: %#v, v: %d\n", t, v)
}

func main() {
	var t T = []int{1, 3}
	t.append(5)
	t.change(5)
}
```

### Remove dl from GitHooks

```bash
$ dl remove .
```

## Contribution
Please feel free to make [issues](https://github.com/task4233/dl/issues/new/choose) and pull requests.

## Author
[task4233](https://github.com/task4233)
