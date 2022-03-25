package dl_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/task4233/dl/v2"
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
			err := dl.New().Run(context.Background(), "v0.0.0", tt.args)
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
