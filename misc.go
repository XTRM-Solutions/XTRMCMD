package main

import (
	"fmt"
	"os"
)

// all of the XTRM api requests use this status structure
type OperationStatus struct {
	Success bool   `json:"Success"`
	Errors  string `json:"Errors"`
} // `json:"OperationStatus"`

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

// GetFlagBool
// Get the true/false of a boolean program flag
func GetFlagBool(key string) (val bool) {
	val, err := nFlags.GetBool(key)
	if nil != err {
		if nil != xLog {
			xLog.Printf("error fetching value for boolean flag [ %s ]\n", key)
		}
		return false
	}
	return val
}

// GetFlagBool
// Get the value of a boolean program flag
func GetFlagString(key string) (val string) {
	val, err := nFlags.GetString(key)
	if nil != err {
		if nil != xLog {
			xLog.Printf("error fetching value for string flag [ %s ]\n", key)
		}
		return ""
	}
	return val
}
