package main

import "github.com/task4233/dl/v2"

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

func SwitchStmt() {
	var v int = 4
	switch v {
	case 1:
		dl.Info(typ)
	default:
		dl.Info(typ)
	}
}
