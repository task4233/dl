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

// An IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int { return len(h) }

// greater order
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
