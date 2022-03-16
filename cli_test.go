package delog

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func init() {

	unzip("testdata.zip", "./")
}

func TestRun(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args    []string
		wantErr bool
	}{
		// TODO: write tests in an other test for "clean"
		"success with clean": {
			args:    []string{"clean", "testdata/a.go"},
			wantErr: false,
		},
		"success with clean 2": {
			args:    []string{"clean", "testdata/b.go"},
			wantErr: false,
		},
		"success with no arg": {
			args:    []string{"clean"},
			wantErr: false,
		},
		"failed with invalid file extention": {
			args:    []string{"clean", "testdata/a.txt"},
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

	cli := New()

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := cli.Run(context.Background(), tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unexpected error, wantError=%v, got=%v", tt.wantErr, err)
			}
		})
	}
}

func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
