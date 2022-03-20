package dl

import (
	"context"
	"testing"
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
