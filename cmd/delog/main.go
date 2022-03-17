package main

import (
	"context"
	"fmt"
	"os"

	"github.com/task4233/delog"
)

var version string

func main() {
	cli := delog.New()
	if err := cli.Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: \n%v\n", err)
		os.Exit(1)
	}
}
