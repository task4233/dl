package dl

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// Info gives a val, a type, a file name, a line number and writes to w..
func FInfo[T any](w io.Writer, v T) (int, error) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return Fprintf(w, "failed FInfo: %v", v)
	}

	file = file[strings.LastIndex(file, "/")+1:]
	if any(v) == nil {
		return Fprintf(w, "[DeLog] info: nil\n%s:%d\n", file, line)
	}
	return Fprintf(w, "[DeLog] info: %#v (%T)\n%s:%d\n", v, v, file, line)
}

// Info gives a val, a type, a file name, a line number to print to the standard logger.
func Info[T any](v T) (int, error) {
	return FInfo(os.Stderr, v)
}

// Fprintf formats according to a format specifier and writes to w.
// Arguments are handled in the manner of fmt.FPrintf.
func Fprintf(w io.Writer, format string, v ...any) (int, error) {
	return fmt.Fprintf(w, format, v...)
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// Arguments are handled in the manner of fmt.FPrintln.
func Fprintln(w io.Writer, v ...any) (int, error) {
	return fmt.Fprintln(w, v...)
}

// Printf calls Fprintf to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...any) (int, error) {
	return Fprintf(os.Stderr, format, v...)
}

// Println calls Fprintln to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Println(v ...any) (int, error) {
	return Fprintln(os.Stderr, v...)
}
