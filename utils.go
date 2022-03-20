package dl

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// copyFile is an utility function to copy file.
func copyFile(ctx context.Context, dstFilePath string, srcFilePath string) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

func walkDirWithValidation(ctx context.Context, baseDir string, fn func(path string, info fs.DirEntry) error) error {
	// might be good running concurrently? TODO(#7)
	return filepath.WalkDir(baseDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walkDir: %w", err)
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		return fn(path, info)
	})
}
