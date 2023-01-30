package dl

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/go-logr/logr"
)

// Logger is a struct for preserving *logr.Logger.
type Logger struct {
	*logr.Logger
}

// NewLogger wraps logr.Logger.
func NewLogger(l *logr.Logger) *Logger {
	return &Logger{
		l,
	}
}

// FInfo gives a val, a type, a file name, a line number and writes to w..
func FInfo[T any](w io.Writer, v T) (int, error) {
	return finfo(w, v, 2)
}

func finfo[T any](w io.Writer, v T, depth int) (int, error) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		return Fprintf(w, "failed FInfo: %v", v)
	}

	file = file[strings.LastIndex(file, "/")+1:]
	if any(v) == nil {
		return Fprintf(w, "[DeLog] info: nil %s:%d\n", file, line)
	}
	return Fprintf(w, "[DeLog] info: %#v (%T) %s:%d\n", v, v, file, line)
}

// Info gives a val, a type, a file name, a line number to print to the standard logger.
func Info[T any](v T) (int, error) {
	return finfo(os.Stderr, v, 2)
}

// Fprintf formats according to a format specifier and writes to w.
// Arguments are handled in the manner of fmt.Fprintf.
func Fprintf(w io.Writer, format string, v ...any) (int, error) {
	return fmt.Fprintf(w, format, v...)
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// Arguments are handled in the manner of fmt.Fprintln.
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
