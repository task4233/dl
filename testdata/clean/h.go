package main

import (
	"time"

	"github.com/task4233/dl/v2"
)

const hoge = "hoge"

const (
	fuga = "fuga"
)

func AssignStmt() {
	_, _ = dl.Info()
}

func BlockStmt() {
	{
		dl.Info()
	}
}

func ForStmt() {
	for i := 0; i < 3; i++ {
		dl.Info(i)
	}
}

func IfStmt() {
	if hoge == fuga {
		dl.Info(i)
	}
}

func SelectStmt() {
	var ch chan int
	go func() {
		time.Sleep(1)
		ch <- 1
	}()

	select {
	case v := <-ch:
		dl.Info(v)
	}
}

func TypeSwitchStmt() {
	var interf interface{}
	switch typ := interf.(type) {
	case int:
		dl.Info(typ)
	default:
		dl.Info(typ)
	}
}

func ReturnStmt() (int, error) {
	return dl.Info(fuga)
}
