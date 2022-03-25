package dl

import (
	"container/heap"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCopyFile(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		srcFilePath string
		dstFilePath string
		wantErr     bool
	}{
		"failed with unexisted src file path": {
			srcFilePath: "unexisted/filepath",
			wantErr:     true,
		},
		"failed with irregal dst file path": {
			srcFilePath: "testdata/restore/test.go",
			dstFilePath: "/root",
			wantErr:     true,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := copyFile(context.Background(), tt.dstFilePath, tt.srcFilePath)
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Fatalf("unexpected error: want=%v, got=%v", tt.wantErr, err)
				}
				return
			}
		})
	}
}

func TestintHeap(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args []int
		want []int
	}{
		"3, 1, 5, 7, 9": {
			args: []int{3, 1, 5, 7, 9},
			want: []int{9, 7, 5, 3, 1},
		},
		"9, 5, 7, 3, 1": {
			args: []int{9, 5, 7, 3, 1},
			want: []int{9, 7, 5, 3, 1},
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			h := newintHeap(tt.args)

			got := make([]int, 0, h.Len())

			for h.Len() > 0 {
				got = append(got, heap.Pop(h).(int))
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("-want,+got\n%s", diff)
			}
		})

	}

}
