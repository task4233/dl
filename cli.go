package delog

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type Cli struct {
	*sweeper
}

func New() *Cli {
	return &Cli{
		NewDelog(),
	}
}

func (c *Cli) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("no argument")
	}
	switch args[0] {
	case "clean":
		if len(args) == 1 {
			args = append(args, ".")
		}
		return c.Clean(ctx, args[1])
	default:
		return fmt.Errorf("command %s is not implemented", args[0])
	}
}

// Clean deletes all methods related to delog in ".go" files under the given directory path
func (c *Cli) Clean(ctx context.Context, baseDir string) error {
	return filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		fmt.Printf("path: %s\n", path)

		if strings.HasSuffix(path, ".go") {
			// might be good running concurrently?
			return c.sweeper.Clean(ctx, path)
		}
		return nil
	})
}
