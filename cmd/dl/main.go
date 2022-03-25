package main

import (
	"context"
	"fmt"
	"os"

	"github.com/task4233/dl/v2"
)

var version string

func main() {
	if err := dl.New().Run(context.Background(), version, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: \n%v\n", err)
		os.Exit(1)
	}
}
