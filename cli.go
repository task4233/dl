package delog

import (
	"context"
	"errors"
	"fmt"
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
			args = append(args, "sandbox/hello.go") // TODO(#4): fix here
		}
		return c.Clean(ctx, args[1])
	default:
		return fmt.Errorf("command %s is not implemented", args[0])
	}
}

// Clean deletes all methods related to delog in ".go" files under the given directory path
func (c *Cli) Clean(ctx context.Context, baseDir string) error {
	// TODO(#4): check recursively all files if baseDir is a directory
	// just in case, run with an assumption: baseDir == a single file
	targetFilePath := baseDir

	if err := c.sweeper.Clean(ctx, targetFilePath); err != nil {
		return err
	}

	return nil
}
