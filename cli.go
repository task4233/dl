package dl

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// dl structs
type dl struct {
	*sweeper
}

// New for running dl package with CLI
func New() *dl {
	return &dl{
		newdl(),
	}
}

// Run executes each method for dl package
func (c *dl) Run(ctx context.Context, args []string) error {
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

// Clean deletes all methods related to dl in ".go" files under the given directory path
func (c *dl) Clean(ctx context.Context, baseDir string) error {
	return filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			fmt.Fprintf(os.Stderr, "remove dl from %s\n", path)
			// might be good running concurrently?
			return c.sweeper.Sweep(ctx, path)
		}
		return nil
	})
}
