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
