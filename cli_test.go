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
			args:    []string{"clean", "testdata/run/clean"},
			wantErr: false,
		},
		"success with init": {
			args:    []string{"init", "testdata/run"},
			wantErr: false,
		},
		"success with remove": {
			args:    []string{"remove", "testdata/run/remove"},
			wantErr: false,
		},
		"success restore": {
			args:    []string{"restore", "testdata/run/restore"},
			wantErr: false,
		},
		"failed restore with uninited directory": {
			args:    []string{"restore", "testdata/run/clean-with-uninited"},
			wantErr: true,
		},
		"failed clean with uninited directory": {
			args:    []string{"clean", "testdata/run/clean-with-uninited"},
			wantErr: true,
		},
		"failed clean with excluded directory": {
			args:    []string{"clean", "testdata/run/clean/.dl"},
			wantErr: true,
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
			err := cli.Run(context.Background(), "v0.0.0", tt.args)
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
			dlDirPath := filepath.Join(tt.baseDir, ".dl")
			if _, err := os.Stat(dlDirPath); os.IsNotExist(err) {
				t.Fatalf("failed to create .dl dir: %v", err)
			}
		})
	}
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
			if err := New().Remove(context.Background(), tt.baseDir); err != nil {
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
