package main

import (
	"fmt"
	"os"
)

// for those functions that should be deferred,
// yet return errors, this handles & reports the error
// cannot assume Log is {yet, still} active, so do not log
func DeferError(f func() error) {
	err := f()
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr,
			"%s%s\n",
			"(may be harmless) error in deferred function: ",
			err.Error())
	}
}
