package main

import (
	"encoding/json"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
)

func main() {

	s, err := homedir.Dir()
	if nil != err {
		fmt.Fprintf(os.Stdout, "error locating homedir: %s \n", err.Error())
	} else {
		fmt.Fprintf(os.Stdout, "found homedir: %s\n", s)
		os.Exit(0)
	}

	InitLog()
	defer DeferError(xLogFile.Close)

	InitFlags()

	if GetFlagBool("quiet") {
		setQuietLog()
	}

	InitConfig()

	if "" != GetFlagString("payfile") {
		MockPayments()
		os.Exit(0)
	}

	xAuthorize(
		"POST",
		xData["apiAuthorizeUrl"],
		xData["xClient"],
		xData["xSecret"],
	)

	if GetFlagBool("debug") {
		xLog.Print("Received access token: " +
			xData["TokenType"] + " " +
			xData["AccessToken"])
	}

	var sendMoney TransferFundToDynamicAccountUserStruct
	sm := &sendMoney.TransferFundToDynamicAccountUser.Request
	sm.IssuerAccountNumber = xData["xIssuerID"]
	sm.FromAccountNumber = xData["xIssuerID"]
	sm.FromWalletID = xData["xDefaultWallet"]
	sm.Amount = GetFlagString("amount")
	sm.Currency = GetFlagString("currency")
	sm.RecipientFirstName = GetFlagString("firstname")
	sm.RecipientLastName = GetFlagString("lastname")
	sm.Description = GetFlagString("description")
	sm.RecipientEmail = GetFlagString("payee")

	tResp, err := xTransferDynamic(sendMoney)

	if nil != err {
		log.Fatal("TransferFundToDynamicAccountUser failed because: " + err.Error())
	}

	tr := &tResp.TransferFundToDynamicAccountUserResponse.TransferFundToDynamicAccountUserResult

	if !tr.OperationStatus.Success {
		xLog.Println("TransferFundToDynamicAccountUser failed because:" +
			tr.OperationStatus.Errors)
	}

	if GetFlagBool("debug") {
		jsonData, err := json.MarshalIndent(tResp, "", "  ")
		if nil != err {
			log.Fatalf("could not unmarshal JSON response because %s\n", err.Error())
		}
		log.Printf("jsondata response: \n%s\n", string(jsonData))
	}

	if !GetFlagBool("quiet") {
		fmt.Printf("\nSuccess! TransactionID %s for %s (%s transferred, %s fee) %s to recipient %s (%s %s %s)\n",
			tr.TransactionID, tr.TotalAmount, tr.Amount,
			tr.Fee, tr.Currency, tr.RecipientAccountNumber,
			sm.RecipientFirstName, sm.RecipientLastName, sm.RecipientEmail)
	}

}
