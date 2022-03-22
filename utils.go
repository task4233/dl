package dl

import (
	"container/heap"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
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

func walkDirWithValidation(ctx context.Context, baseDir string, fn func(ctx context.Context, path string, info fs.DirEntry) error) error {
	eg, ctx := errgroup.WithContext(ctx)

	err := filepath.WalkDir(baseDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walkDir: %w", err)
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		eg.Go(func() error {
			return fn(ctx, path, info)
		})
		return nil
	})
	if err != nil {
		return err
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// An IntHeap is a min-heap of ints.
type IntHeap struct {
	s  []int
	mu *sync.Mutex
}

// NeaIntHeap is a factory function for IntHeap.
// heap.Init is done in this function.
func NewIntHeap(s []int) *IntHeap {
	if s == nil {
		s = []int{}
	}
	h := &IntHeap{
		s:  s,
		mu: new(sync.Mutex),
	}
	heap.Init(h)
	return h
}

func (h IntHeap) Len() int { return len(h.s) }

// greater order
func (h IntHeap) Less(i, j int) bool { return h.s[i] > h.s[j] }
func (h IntHeap) Swap(i, j int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.s[i], h.s[j] = h.s[j], h.s[i]
}

func (h *IntHeap) Push(x any) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	h.s = append(h.s, x.(int))
}

func (h *IntHeap) Pop() any {
	h.mu.Lock()
	defer h.mu.Unlock()
	old := h.s
	n := len(old)
	x := old[n-1]
	h.s = old[0 : n-1]
	return x
}
