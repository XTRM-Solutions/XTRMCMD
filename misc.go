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

func getFlagBool(val string) bool {
	rc, err := nFlags.GetBool(val)
	if nil != err {
		if nil != xLog {
			xLog.Printf("error fetching value for boolean flag [ %s ]\n", val)
		}
		return false
	}
	return rc
}

func getFlagString(val string) string {
	rc, err := nFlags.GetString(val)
	if nil != err {
		if nil != xLog {
			xLog.Printf("error fetching value for string flag [ %s ]\n", val)
		}
		return ""
	}
	return rc
}
