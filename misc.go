package xtrmcmd

import (
	"fmt"
	"os"
)

// for those irritating functions that should be deferred,
// yet return errors, this handles & reports the error
// cannot assume Log is {yet, still} active, so do not log
func DeferError(f func() error) {
	err := f()
	if nil != err {
		_, _ = fmt.Fprintln(
			os.Stderr,
			"(may be harmless) error in deferred function: "+err.Error(),
		)
	}
}
