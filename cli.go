package dl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
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
	buf, err := io.ReadAll(f)
	if err != nil {
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

// Remove deletes `$ dl clean` command from pre-commit script
func (d *DeLog) Remove(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, ".git", "hooks", "pre-commit")
	buf, err := readFile(ctx, path)
	if err != nil {
		return err
	}
	if buf == nil {
		return nil
	}

	return removePrecommitScript(ctx, path, buf)
}

// readFile is created for calling f.Close() surely
func readFile(ctx context.Context, path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if a pre-commit file does not exist, this method has no effect.
		return nil, nil
	}

	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

func removePrecommitScript(ctx context.Context, path string, buf []byte) error {
	idx := bytes.Index(buf, []byte(precommitScript))
	if idx < 0 {
		return nil
	}
	if len(buf) == len(precommitScript) {
		if err := os.Remove(path); err != nil {
			return err
		}
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// 0 <= idx && idx < len(buf)-1	=> append(buf[:idx], buf[idx+len(precommitScript):]
	// idx == len(buf)-len(preCommitScript)	=> buf[:idx]
	if idx == len(buf)-len(precommitScript) {
		if _, err := f.Write(buf[:idx]); err != nil {
			return err
		}
	} else {
		if _, err := f.Write(append(buf[:idx], buf[idx+len(precommitScript):]...)); err != nil {
			return err
		}
	}

	return nil
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
