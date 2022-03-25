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
	type MyInt int
	var myNum MyInt = 123

	dl.Fprintln(os.Stdout, "num: ", num)
	dl.Println("num: ", num)
	dl.Fprintf(os.Stdout, "name: %s\n", name)
	dl.Printf("name: %s", name)
	dl.FInfo(os.Stdout, myNum)
	dl.Info(myNum)

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
	//[DeLog] info: 123 (dl_test.MyInt) log_example_test.go:23
}

func ExampleFInfo() {
	type person struct {
		name string
		age  int
	}
	alice := person{
		name: "alice",
		age:  15,
	}

	_, _ = dl.FInfo(os.Stdout, alice)
	// Output: [DeLog] info: dl_test.person{name:"alice", age:15} (dl_test.person) log_example_test.go:50
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

func ExampleInfo() {
	type person struct {
		name string
		age  int
	}
	alice := person{
		name: "alice",
		age:  15,
	}

	// dl.Info prints to sandard error.
	// [DeLog] info: dl_test.person{name:"alice", age:15} (dl_test.person) log_example_test.go:96
	_, _ = dl.Info(alice)
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
