package dl

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
)

var _ Cmd = (*Remove)(nil)

type Remove struct{}

func NewRemove() *Remove {
	return &Remove{}
}

// Remove removes environment settings for dl commands.
func (r *Remove) Run(ctx context.Context, baseDir string) error {
	if err := r.removeScript(ctx, filepath.Join(baseDir, ".git", "hooks", "pre-commit"), preCommitScript); err != nil {
		return err
	}
	if err := r.removeScript(ctx, filepath.Join(baseDir, ".git", "hooks", "post-commit"), postCommitScript); err != nil {
		return err
	}

	return nil
}

func (r *Remove) removeScript(ctx context.Context, path string, script string) error {
	buf, err := r.readFile(ctx, path)
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
func (*Remove) readFile(ctx context.Context, path string) ([]byte, error) {
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
