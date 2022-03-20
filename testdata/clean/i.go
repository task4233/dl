package main

import "github.com/task4233/dl"

func RangeStmt() {
	s := []int{1, 3, 5}
	f := func() []int {
		dl.Info(s)
		return s
	}
	for _, elem := range f() {
		_ = elem
	}
}
