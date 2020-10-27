package main

import "fmt"

func DeferError(f func() error) {
	err := f()
	if nil != err {
		fmt.Println("deferred error: " + err.Error())
	}
}
