package dl

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	once.Do(extractZip)
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		targetPath string
		args       []string
		wantErr    bool
	}{
		"success with clean": {
			args:    []string{"clean", "testdata/run"},
			wantErr: false,
		},
		"success with init": {
			args:    []string{"init", "testdata/run"},
			wantErr: false,
		},
		"no effect with invalid file extension": {
			args:    []string{"clean", "testdata/a.txt"},
			wantErr: false,
		},
		"failed with unknown command": {
			args:    []string{"hoge"},
			wantErr: true,
		},
		"failed with no arg": {
			args:    []string{},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			cli := New()

			err := cli.Run(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error, wantError=%v, got=%v", tt.wantErr, err)
			}

			if tt.targetPath != "" {
				data, err := os.ReadFile(tt.targetPath)
				if err != nil {
					t.Fatalf("failed ReadFile: %s", err.Error())
				}
				if bytes.Contains(data, []byte("dl")) {
					t.Fatalf("failed to delete dl from data: \n%s", string(data))
				}
			}
		})
	}
}

func TestInit(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		baseDir string
		wantErr bool
	}{
		"success with": {
			baseDir: "testdata/init",
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

			if err := New().Init(context.Background(), tt.baseDir); err != nil {
				if (err != nil) != tt.wantErr {
					t.Fatalf("failed Init: err=%v", err)
				}
				return
			}

			precommitFilePath := filepath.Join(tt.baseDir, ".git", "hooks", "pre-commit")

			if _, err := os.Stat(precommitFilePath); os.IsNotExist(err) {
				t.Fatalf("failed to create pre-commit script")
			}

			data, err := os.ReadFile(precommitFilePath)
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			if !bytes.Contains(data, []byte("dl clean")) {
				t.Fatalf("pre-commit script is not installed")
			}

		})
	}
}
