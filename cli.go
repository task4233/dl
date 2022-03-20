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
func (d *DeLog) Run(ctx context.Context, version string, args []string) error {
	if len(args) == 0 {
		d.usage(version, "")
		return errors.New("no command is given.")
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
	case "restore":
		return d.Restore(ctx, args[1])
	default:
		return d.usage(version, args[0])
	}
}

var (
	excludedFiles = []string{dlDirPath, ".git"}
)

// Clean deletes all methods related to dl in ".go" files under the given directory path
func (d *DeLog) Clean(ctx context.Context, baseDir string) error {
	dlDirPath := filepath.Join(baseDir, dlDirPath)
	if _, err := os.Stat(dlDirPath); os.IsNotExist(err) {
		return fmt.Errorf(".dl directory doesn't exist. Please execute $ dl init .: %s", dlDirPath)
	}

	return filepath.WalkDir(baseDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walkDir: %w", err)
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		for _, file := range excludedFiles {
			if strings.Contains(path, file) {
				return nil
			}
		}
		if err := d.Sweeper.Evacuate(ctx, baseDir, path); err != nil {
			return fmt.Errorf("failed to evacuate %s, %s", path, err.Error())
		}

		// might be good running concurrently? TODO(#7)
		fmt.Fprintf(os.Stderr, "remove dl from %s\n", path)
		return d.Sweeper.Sweep(ctx, path)
	})
}

// Although these values will be casted to []byte, they are declared as constants
// because that is't happened frequently.
const (
	preCommitScript = `#!/bin/sh
dl clean .
git add .
`
	postCommitScript = `#!/bin/sh
dl restore .
`
	cleanCmd   = "dl clean"
	restoreCmd = "dl restore"
	dlDirPath  = ".dl"
)

// Init inserts dl command into git pre-commit hook
func (d *DeLog) Init(ctx context.Context, baseDir string) error {
	if err := d.addGitPreHookScript(ctx, baseDir); err != nil {
		return err
	}
	if err := d.addGitPostHookScript(ctx, baseDir); err != nil {
		return err
	}
	if err := d.createDlDirIfNotExist(ctx, baseDir); err != nil {
		return err
	}
	if err := d.addDlIntoGitIgnore(ctx, baseDir); err != nil {
		return err
	}

	return nil
}

func (d *DeLog) addGitPreHookScript(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, ".git", "hooks")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	preCommitPath := filepath.Join(path, "pre-commit")
	f, err := os.OpenFile(preCommitPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
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

	if _, err := fmt.Fprintf(f, preCommitScript); err != nil {
		return err
	}

	return os.Chmod(preCommitPath, 0755)
}

func (d *DeLog) addGitPostHookScript(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, ".git", "hooks")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	postCommitPath := filepath.Join(path, "post-commit")
	f, err := os.OpenFile(postCommitPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// It checks if `$ dl restore` has been installed or not.
	// If so, not inserting codes.
	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if bytes.Contains(buf, []byte(restoreCmd)) {
		return nil
	}

	if _, err := fmt.Fprintf(f, postCommitScript); err != nil {
		return err
	}

	return os.Chmod(postCommitPath, 0755)
}

func (d *DeLog) addDlIntoGitIgnore(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, ".gitignore")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if bytes.Contains(buf, []byte(dlDirPath)) {
		return nil
	}

	if _, err := fmt.Fprintf(f, "\n%s\n", dlDirPath); err != nil {
		return err
	}

	return nil
}

func (d *DeLog) createDlDirIfNotExist(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, dlDirPath)
	if stat, err := os.Stat(path); err == nil {
		if stat.IsDir() {
			return nil
		}
		return fmt.Errorf("%s has been already existed as file. Please rename or delete it.", path)
	}

	return os.Mkdir(path, 0755)
}

// Remove deletes `$ dl clean` command from pre-commit script
func (d *DeLog) Remove(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, ".git", "hooks")

	if err := removeScript(ctx, filepath.Join(path, "pre-commit"), preCommitScript); err != nil {
		return err
	}

	if err := removeScript(ctx, filepath.Join(path, "post-commit"), postCommitScript); err != nil {
		return err
	}

	return nil
}

func removeScript(ctx context.Context, path string, script string) error {
	buf, err := readFile(ctx, path)
	if err != nil {
		return err
	}
	if buf == nil {
		return nil
	}

	idx := bytes.Index(buf, []byte(script))
	if idx < 0 {
		return nil
	}
	if len(buf) == len(script) {
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

	// 0 <= idx && idx < len(buf)-1	=> append(buf[:idx], buf[idx+len(script):]
	// idx == len(buf)-len(script)	=> buf[:idx]
	if idx == len(buf)-len(script) {
		if _, err := f.Write(buf[:idx]); err != nil {
			return err
		}
	} else {
		if _, err := f.Write(append(buf[:idx], buf[idx+len(script):]...)); err != nil {
			return err
		}
	}

	return nil
}

// readFile is created for calling f.Close() surely
func readFile(ctx context.Context, path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if a pre-commit file does not exist, this method has no effect.
		return nil, nil
	}

	f, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

// Restore restores raw files from .dl directory
func (d *DeLog) Restore(ctx context.Context, baseDir string) error {
	dlDirPath := filepath.Join(baseDir, dlDirPath)
	if _, err := os.Stat(dlDirPath); os.IsNotExist(err) {
		return fmt.Errorf(".dl directory doesn't exist. Please execute $ dl init .: %s", dlDirPath)
	}

	return filepath.WalkDir(baseDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walkDir: %w", err)
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// might be good running concurrently? TODO(#7)
		idx := strings.Index(path, dlDirPath)
		if idx < 0 {
			return nil
		}
		return d.Sweeper.Restore(ctx, path, path[:idx]+path[idx+len(dlDirPath)+1:])
	})
}

const msg = "dl %s: The instant logger package for debug.\n"

func (d *DeLog) usage(version, invalidCmd string) error {
	fmt.Fprintf(os.Stderr, msg+`Usage: dl [command]
Commands:
clean <dir>                 deletes logs used this package.
init <dir>                  add dl command into pre-commit script.
remove <dir>                remove dl command from pre-commit script.
restore <dir>               restore removed logs.
`, version)
	return fmt.Errorf("%s is not implemented.", invalidCmd)
}
