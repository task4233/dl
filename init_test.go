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

func TestInit(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		baseDir string
		wantErr bool
	}{
		"success": {
			baseDir: "testdata/init",
			wantErr: false,
		},
		"success in inited directory": {
			baseDir: "testdata/inited",
			wantErr: false,
		},
		"failed when .dl, which is file, exists": {
			baseDir: "testdata/inited-with-dl-file",
			wantErr: true,
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
			if err := dl.NewInit().Run(context.Background(), tt.baseDir); err != nil {
				if (err != nil) != tt.wantErr {
					t.Fatalf("failed Init: err=%v", err)
				}
				return
			}
			preCommitFilePath := filepath.Join(tt.baseDir, ".git", "hooks", "pre-commit")
			if _, err := os.Stat(preCommitFilePath); os.IsNotExist(err) {
				t.Fatalf("failed to create pre-commit script: %v", err)
			}
			data, err := os.ReadFile(preCommitFilePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}
			if !bytes.Contains(data, []byte("dl clean")) {
				t.Fatalf("pre-commit script is not installed")
			}

			dlDirPath := filepath.Join(tt.baseDir, ".dl")
			if _, err := os.Stat(dlDirPath); os.IsNotExist(err) {
				t.Fatalf("failed to create .dl dir: %v", err)
			}

			postCommitFilePath := filepath.Join(tt.baseDir, ".git", "hooks", "post-commit")
			data, err = os.ReadFile(postCommitFilePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}
			if !bytes.Contains(data, []byte("dl restore")) {
				t.Fatalf("post-commit script is not installed")
			}

			gitignoreFilePath := filepath.Join(tt.baseDir, ".gitignore")
			data, err = os.ReadFile(gitignoreFilePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}
			if !bytes.Contains(data, []byte(".dl")) {
				t.Fatalf("gitignore is not configured")
			}
		})
	}
}
