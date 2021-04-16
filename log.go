package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var xLog *log.Logger
var xLogFile *os.File

func InitLog() {
	var err error

	xLogFile, err = os.OpenFile("xtrmcmd.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		644)
	if nil != err {
		fmt.Fprintln(os.Stderr, "Could not open logging file xtrmcmd.log because "+err.Error())
	}

	xbf := io.MultiWriter(xLogFile, os.Stderr)
	xLog = log.New(xbf, "xtrmcmd: ", log.Lshortfile|log.LUTC)

}

func SetQuietLog() {
	xLog.SetOutput(xLogFile)
}
