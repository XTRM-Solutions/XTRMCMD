package main

import (
	"fmt"
	"os"
)

func DeferError(f func() error) {
	err := f()
	if nil != err {
		_, _ = fmt.Fprintln(os.Stderr, "deferred error: "+err.Error())
	}
}
