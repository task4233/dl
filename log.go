package delog

import (
	"fmt"
	"io"
	"os"
)

func Fprintf(w io.Writer, format string, v ...any) {
	fmt.Fprintf(w, format, v...)
}

func Fprintln(w io.Writer, v ...any) {
	fmt.Fprintln(w, v...)
}

func Printf(format string, v ...any) {
	Fprintf(os.Stderr, format, v...)
}
func Println(v ...any) {
	Fprintln(os.Stderr, v...)
}
