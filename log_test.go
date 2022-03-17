package dl_test

import (
	"bytes"
	"testing"

	"github.com/task4233/dl"
)

func TestFPrintf(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		message string
		args    []any
		want    string
	}{
		"success with string": {
			message: "hoge",
			args:    nil,
			want:    "hoge",
		},
		"success with format string": {
			message: "hoge: %s",
			args:    []any{"fuga"},
			want:    "hoge: [fuga]",
		},
		"success with empty string": {
			message: "",
			args:    nil,
			want:    "",
		},
		"success with nil": {
			message: "",
			args:    nil,
			want:    "",
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := new(bytes.Buffer)
			if tt.args == nil {
				dl.Fprintf(out, tt.message)
			} else {
				dl.Fprintf(out, tt.message, tt.args)
			}

			if tt.want != out.String() {
				t.Fatalf("failed TestPrintf, want=%s, got=%s", tt.want, out.String())
			}
		})
	}
}

func TestFPrintln(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args []any
		want string
	}{
		"success with string": {
			args: nil,
			want: "<nil>\n",
		},
		"success with format string": {
			args: []any{"fuga"},
			want: "[fuga]\n",
		},
	}

	for name, tt := range tests {
		tt := tt

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			out := new(bytes.Buffer)
			if tt.args == nil {
				dl.Fprintln(out, nil)
			} else {
				dl.Fprintln(out, tt.args)
			}

			if tt.want != out.String() {
				t.Fatalf("failed TestPrintf, want=%s, got=%s", tt.want, out.String())
			}
		})
	}
}
