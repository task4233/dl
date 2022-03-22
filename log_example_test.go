package dl_test

import (
	"os"

	"github.com/task4233/dl"
)

func Example() {
	num := 123
	name := "dl"
	type MyInt int
	var myNum MyInt = 123
	var names []string = []string{"alice", "bob", "Charlie"}

	dl.Fprintln(os.Stdout, "num: ", num)
	dl.Println("num: ", num)

	dl.Fprintf(os.Stdout, "name: %s\n", name)
	dl.Printf("name: %s", name)

	dl.FInfo(os.Stdout, myNum)
	dl.Info(myNum)

	dl.FInfo(os.Stdout, names)
	dl.Info(names)

	// Output:
	// num:  123
	// name: dl
	//[DeLog] info: 123 (dl_test.MyInt) log_example_test.go:22
	//[DeLog] info: []string{"alice", "bob", "Charlie"} ([]string) log_example_test.go:25

}
