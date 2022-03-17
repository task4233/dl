package main

import (
	"context"
	"fmt"
	"os"

	"github.com/task4233/dl"
)

func main() {
	cli := dl.New()
	if err := cli.Run(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: \n%v\n", err)
		os.Exit(1)
	}
}
