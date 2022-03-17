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

	dl.Fprintln(os.Stdout, "num: ", num)
	dl.Fprintf(os.Stdout, "name: %s\n", name)
	dl.FInfo(os.Stdout, myNum)

	// Output:
	// num:  123
	// name: dl
	//[DeLog] info: 123 (dl_test.MyInt) log_example_test.go:17
}
