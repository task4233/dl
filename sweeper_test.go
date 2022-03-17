package dl

import (
	"bytes"
	"context"
	"os"
	"testing"
)

func init() {
	once.Do(extractZip)
}

// TestClean invokes cli.Clean which invokes sweeper.Sweep inside
func TestClean(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		targetPath string
		args       []string
		wantErr    bool
	}{
		"success": {
			targetPath: "testdata/clean/a.go",
			args:       []string{"clean", "testdata/clean/a.go"},
			wantErr:    false,
		},
		"success with packager alias": {
			targetPath: "testdata/clean/b.go",
			args:       []string{"clean", "testdata/clean/b.go"},
			wantErr:    false,
		},
		"success with only dl package": {
			targetPath: "testdata/clean/c.go",
			args:       []string{"clean", "testdata/clean/c.go"},
			wantErr:    false,
		},
		"success with only dl package alias": {
			targetPath: "testdata/clean/d.go",
			args:       []string{"clean", "testdata/clean/d.go"},
			wantErr:    false,
		},
		"success with dl package oneliner": {
			targetPath: "testdata/clean/e.go",
			args:       []string{"clean", "testdata/clean/e.go"},
			wantErr:    false,
		},
		"success with dl package alias oneliner": {
			targetPath: "testdata/clean/f.go",
			args:       []string{"clean", "testdata/clean/f.go"},
			wantErr:    false,
		},
		// fault cases are tested in TestRun
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
