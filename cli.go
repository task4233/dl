package dl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// DeLog structs
type DeLog struct {
	*Sweeper
}

// New for running dl package with CLI
func New() *DeLog {
	return &DeLog{
		NewSweeper(),
	}
}

// Run executes each method for dl package
func (d *DeLog) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("no argument")
	}
	if len(args) == 1 {
		args = append(args, ".")
	}

	switch args[0] {
	case "clean":
		return d.Clean(ctx, args[1])
	case "init":
		return d.Init(ctx, args[1])
	case "remove":
		return d.Remove(ctx, args[1])
	default:
		return d.usage(args[0])
	}
}

func (d *DeLog) usage(invalidCmd string) error {
	msg := "%s is not implemented.\n"
	fmt.Fprintf(os.Stderr, msg+
		`Usage: dl [command]
Commands:
init <dir>                  add dl command into pre-commit.
clean <dir>                 deletes logs used this package.
`, invalidCmd)
	return fmt.Errorf(msg, invalidCmd)
}

// Clean deletes all methods related to dl in ".go" files under the given directory path
func (d *DeLog) Clean(ctx context.Context, baseDir string) error {
	return filepath.Walk(baseDir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") {
			fmt.Fprintf(os.Stderr, "remove dl from %s\n", path)
			// might be good running concurrently?
			return d.Sweeper.Sweep(ctx, path)
		}
		return nil
	})
}

// Remove deletes `$ dl clean` command from pre-commit script
func (d *DeLog) Remove(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, ".git", "hooks", "pre-commit")
	FInfo(os.Stderr, path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := []byte{}
	if _, err := f.Read(buf); err != nil {
		return err
	}

	idx := bytes.Index(buf, []byte(precommitScript))
	if idx < 0 {
		return nil
	}
	if len(buf) == len(precommitScript) {
		return os.Remove(path)
	}

	if _, err := f.Write(append(buf[:idx], buf[idx+len(precommitScript)+1:]...)); err != nil {
		return err
	}

	return nil
}

// Although these values will be casted to []byte, they are declared as constants
// because that is't happened frequently.
const (
	precommitScript = `#!/bin/sh
dl clean .
git add .
`
	cleanCmd = "dl clean"
)

// Init inserts dl command into git pre-commit hook
func (d *DeLog) Init(ctx context.Context, baseDir string) error {
	if err := d.addGitHookScript(ctx, baseDir); err != nil {
		return err
	}

	// TODO(#17): add feature of commented out
	return nil
}

func (d *DeLog) addGitHookScript(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, ".git", "hooks")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	path = filepath.Join(path, "pre-commit")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	// It checks if `$ dl clean` has been installed or not.
	// If so, not inserting codes.
	buf := []byte{}
	if _, err := f.Read(buf); err != nil {
		return err
	}
	if bytes.Contains(buf, []byte(cleanCmd)) {
		return nil
	}

	if _, err := fmt.Fprintf(f, precommitScript); err != nil {
		return err
	}

	return os.Chmod(path, 0755)
}
