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

// An intHeap is a min-heap of ints.
type intHeap struct {
	s  *[]int
	mu *sync.Mutex
}

// newintHeap is a factory function for intHeap.
// heap.Init is done in this function.
func newintHeap(s []int) *intHeap {
	if s == nil {
		s = []int{}
	}
	h := &intHeap{
		s:  &s,
		mu: new(sync.Mutex),
	}
	heap.Init(h)
	return h
}

func (h intHeap) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.s
	return len(*s)
}

// greater order
func (h intHeap) Less(i, j int) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.s
	return (*s)[i] > (*s)[j]
}
func (h intHeap) Swap(i, j int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.s
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

func (h *intHeap) Push(x interface{}) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	s := h.s
	*s = append(*s, x.(int))
}

func (h *intHeap) Pop() interface{} {
	h.mu.Lock()
	defer h.mu.Unlock()
	old := h.s
	n := len(*old)
	x := (*old)[n-1]
	*old = (*old)[:n-1]
	return x
}
