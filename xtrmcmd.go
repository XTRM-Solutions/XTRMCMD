package xtrmcmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	// "github.com/spf13/cobra" // not yet
)

var xLog *log.Logger
var xLogFile *os.File

// var xLogWriter *bufio.Writer

func InitLog() {
	var err error
	xLogFile, err = os.OpenFile("xtrmcmd.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		644)
	if nil != err {
		fmt.Println("Could not open logging file xtrmcmd.log because " + err.Error())
	}
	// xLogWriter = bufio.NewWriter(xLogFile)
	xbf := io.MultiWriter(xLogFile, os.Stderr)
	xLog = log.New(xbf, "xtrmcmd: ", log.Lshortfile)
	xLog.Print("\nStarted logging for XTRMCMD")
}

func main() {
	InitLog()
	defer DeferError(xLogFile.Close)

	InitFlags()
	InitConfig()

	xAuthorize("POST",
		xData["apiAuthorizeUrl"],
		xData["xClient"],
		xData["xSecret"])

	if *Flags.Debug {
		xLog.Print("Received access token: " +
			xData["TokenType"] + " " +
			xData["AccessToken"])

	}

	var sendMoney TransferFundToDynamicAccountUserStruct
	sm := &sendMoney.TransferFundToDynamicAccountUser.Request
	sm.IssuerAccountNumber = xData["xIssuerID"]
	sm.FromAccountNumber = xData["xIssuerID"]
	sm.FromWalletID = xData["xDefaultWallet"]
	sm.Amount = *Flags.Amount
	sm.Currency = *Flags.Currency
	sm.RecipientFirstName = *Flags.FirstName
	sm.RecipientLastName = *Flags.LastName
	sm.Description = *Flags.Description
	sm.RecipientEmail = *Flags.PayeeEmail

	tResp, err := xTransferDynamic(sendMoney)

	if nil != err {
		log.Fatal(err.Error())

	}

	if *Flags.Debug {
		var s string

		jsonData, err := json.MarshalIndent(tResp, "", "    ")
		if nil != err {
			s = "(could not unmarshal response) error: " + err.Error()
		} else {
			s = string(jsonData)
		}
		log.Fatal(s)
	}

}
