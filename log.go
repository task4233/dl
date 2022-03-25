package dl

import (
	"fmt"
	"io"
	"os"

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

// Fprintf formats according to a format specifier and writes to w.
// Arguments are handled in the manner of fmt.FPrintf.
func Fprintf(w io.Writer, format string, v ...interface{}) (int, error) {
	return fmt.Fprintf(w, format, v...)
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// Arguments are handled in the manner of fmt.FPrintln.
func Fprintln(w io.Writer, v ...interface{}) (int, error) {
	return fmt.Fprintln(w, v...)
}

// Printf calls Fprintf to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) (int, error) {
	return Fprintf(os.Stderr, format, v...)
}

// Println calls Fprintln to print to the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func Println(v ...interface{}) (int, error) {
	return Fprintln(os.Stderr, v...)
}
