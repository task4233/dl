package dl

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func init() {
	once.Do(extractZip)
}

// TestSweep invokes cli.Clean which invokes sweeper.Sweep inside
func TestSweep(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		targetPath string
		wantErr    bool
	}{
		"success": {
			targetPath: "testdata/clean/a.go",
			wantErr:    false,
		},
		"success with packager alias": {
			targetPath: "testdata/clean/b.go",
			wantErr:    false,
		},
		"success with only dl package": {
			targetPath: "testdata/clean/c.go",
			wantErr:    false,
		},
		"success with only dl package alias": {
			targetPath: "testdata/clean/d.go",
			wantErr:    false,
		},
		"success with dl package oneliner": {
			targetPath: "testdata/clean/e.go",
			wantErr:    false,
		},
		"success with dl package alias oneliner": {
			targetPath: "testdata/clean/f.go",
			wantErr:    false,
		},
		"success with non dl package alias oneliner": {
			targetPath: "testdata/clean/g.go",
			wantErr:    false,
		},
		"failed with invalid extension": {
			targetPath: "README.md",
			wantErr:    true,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sweeper := NewSweeper()

			err := sweeper.Sweep(context.Background(), tt.targetPath)
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Fatalf("unexpected error, wantError=%v, got=%v", tt.wantErr, err)
				}
				return
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

// TestEvacuate assumes `.dl` directory is created.
func TestEvacuate(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		baseDirPath    string
		targetFilePath string
		dlFilePath     string
		wantErr        bool
	}{
		"success with a go file": {
			baseDirPath:    "testdata/evacuate",
			targetFilePath: "testdata/evacuate/a.go",
			dlFilePath:     "testdata/evacuate/.dl/a.go",
			wantErr:        false,
		},
		"failed with unexisted file": {
			baseDirPath:    "testdata/evacuate",
			targetFilePath: "testdata/evacuate/hoge.txt",
			wantErr:        true,
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := NewSweeper().Evacuate(context.Background(), tt.baseDirPath, tt.targetFilePath)

			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Fatalf("unexpected error: want=%v, got=%v", tt.wantErr, err)
				}
				return
			}

			// check if existed
			if _, err := os.Stat(tt.dlFilePath); os.IsNotExist(err) {
				t.Fatalf("%s is not created", tt.dlFilePath)
			}
			// check if same file
			raw, err := os.ReadFile(tt.targetFilePath)
			if err != nil {
				t.Fatalf("failed to readFile from %s, err=%v", tt.targetFilePath, err)
			}
			copied, err := os.ReadFile(tt.dlFilePath)
			if err != nil {
				t.Fatalf("failed to readFile from %s, err=%v", tt.dlFilePath, err)
			}
			if diff := cmp.Diff(raw, copied); diff != "" {
				t.Fatalf("copied file is not same to raw file; -raw+copied\n%s", diff)
			}
		})
	}

}
