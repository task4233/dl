package main

import (
	"fmt"

	d "github.com/task4233/dl/v2"
)

type T struct {
	Name string
}

const (
	message = "message"
)

func main() {
	fmt.Println(message)

	const message = "localMessage"

	d.Println(message)
}
