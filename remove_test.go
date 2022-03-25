package dl_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/task4233/dl"
)

func init() {
	once.Do(extractZip)
}

func TestRemove(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		baseDir string
		wantErr bool
	}{
		"success with inited directory": {
			baseDir: "testdata/remove",
			wantErr: false,
		},
		"success with removed directory": {
			baseDir: "testdata/removed",
			wantErr: false,
		},
		"success with pre-commit-script added other commands before": {
			baseDir: "testdata/remove-before",
			wantErr: false,
		},
		"success with pre-commit-script added other commands after": {
			baseDir: "testdata/remove-after",
			wantErr: false,
		},
		"success with pre-commit-script added other commands both before and after": {
			baseDir: "testdata/remove-both",
			wantErr: false,
		},
		"success with pre-commit-script unrelated to dl": {
			baseDir: "testdata/remove-from-unrelated",
			wantErr: false,
		},
		"failed because .git directory does not exists": {
			baseDir: "testdata/clean", // there is not a .git in ./testdata/clean directory
			wantErr: true,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if err := dl.ExportNewRemoveCmd().Run(context.Background(), tt.baseDir); err != nil {
				if (err != nil) != tt.wantErr {
					t.Fatalf("failed Init: err=%v", err)
				}
				return
			}
			precommitFilePath := filepath.Join(tt.baseDir, ".git", "hooks", "pre-commit")
			if _, err := os.Stat(precommitFilePath); os.IsNotExist(err) {
				return
			}
			data, err := os.ReadFile(precommitFilePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}
			if bytes.Contains(data, []byte("dl clean")) {
				t.Fatalf("failed to remove pre-commit script:\n%s", string(data))
			}
		})
	}
}
