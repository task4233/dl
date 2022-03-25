package dl_test

import (
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/task4233/dl"
	"go.uber.org/zap"
)

func Example() {
	num := 123
	name := "dl"

	dl.Fprintln(os.Stdout, "num: ", num)
	dl.Println("num: ", num)
	dl.Fprintf(os.Stdout, "name: %s\n", name)
	dl.Printf("name: %s", name)

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	var log logr.Logger = zapr.NewLogger(zapLog)
	dlr := dl.NewLogger(&log)
	dlr.Info("Info: ", "num", num)

	// Output:
	// num:  123
	// name: dl
}

func ExampleFprintf() {
	type person struct {
		name string
		age  int
	}
	alice := person{
		name: "alice",
		age:  15,
	}

	dl.Fprintf(os.Stdout, "name: %s", alice.name)
	// Output: name: alice
}

func ExampleFprintln() {
	type person struct {
		name string
		age  int
	}
	alice := person{
		name: "alice",
		age:  15,
	}

	dl.Fprintln(os.Stdout, "name:", alice.name)
	// Output: name: alice
	//
}

func ExamplePrintf() {
	type person struct {
		name string
		age  int
	}
	alice := person{
		name: "alice",
		age:  15,
	}

	// dl.Printf prints to sandard error.
	// name: alice
	dl.Printf("name: %s", alice.name)
}

func ExamplePrintln() {
	type person struct {
		name string
		age  int
	}
	alice := person{
		name: "alice",
		age:  15,
	}

	// dl.Printf prints to sandard error.
	// name: alice
	//
	dl.Println("name:", alice.name)
}

func ExampleNewLogger() {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	var log logr.Logger = zapr.NewLogger(zapLog)

	// You can use your logr.Logger as it is.
	dlr := dl.NewLogger(&log)

	num := 57
	dlr.Info("Info: ", "num", num)
}
