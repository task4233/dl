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
