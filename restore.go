package dl

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var _ Cmd = (*Restore)(nil)

type Restore struct{}

func NewRestore() *Restore {
	return &Restore{}
}

// Restore restores raw files from .dl directory
func (r *Restore) Run(ctx context.Context, baseDir string) error {
	dlDirPath := filepath.Join(baseDir, dlDir)
	if _, err := os.Stat(dlDirPath); os.IsNotExist(err) {
		return fmt.Errorf(".dl directory doesn't exist. Please execute $ dl init .: %s", dlDirPath)
	}

	// check files under .dl directory recursively
	return walkDirWithValidation(ctx, baseDir, func(ctx context.Context, path string, info fs.DirEntry) error {
		idx := strings.Index(path, dlDir)
		if idx < 0 {
			return nil
		}

		// copies ".go" files to raw places.
		dstFilePath := path[:idx] + path[idx+len(dlDir)+1:]
		srcFilePath := path
		return copyFile(ctx, dstFilePath, srcFilePath)
	})
}
