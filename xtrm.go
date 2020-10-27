package main

import (
	"encoding/json"
	"fmt"
	"os"
	// "github.com/spf13/cobra" // not yet
)

func main() {

	InitFlags()
	InitConfig()

	xAuthorize("POST",
		xData["apiAuthorizeUrl"],
		xData["xClient"],
		xData["xSecret"])

	if *Flags.Debug {
		fmt.Println("Received access token: " +
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
	sm.RecipientEmail = *Flags.Payee

	tResp, err := xTransferDynamic(sendMoney)

	if nil != err {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	if *Flags.Debug {
		var s string

		jsonData, err := json.MarshalIndent(tResp, "", "    ")
		if nil != err {
			s = "(could not unmarshal response) error: " + err.Error()
		} else {
			s = string(jsonData)
		}
		fmt.Println(s)
	}

}
