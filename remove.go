package dl

import (
	"bytes"
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

var _ cmd = (*removeCmd)(nil)

type removeCmd struct{}

func newRemoveCmd() *removeCmd {
	return &removeCmd{}
}

// Remove removes environment settings for dl commands.
func (r *removeCmd) Run(ctx context.Context, baseDir string) error {
	if err := r.removeScript(ctx, filepath.Join(baseDir, ".git", "hooks", "pre-commit"), preCommitScript); err != nil {
		return err
	}
	if err := r.removeScript(ctx, filepath.Join(baseDir, ".git", "hooks", "post-commit"), postCommitScript); err != nil {
		return err
	}

	return nil
}

func (r *removeCmd) removeScript(ctx context.Context, path string, script string) error {
	buf, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
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
