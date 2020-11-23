package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {

	InitLog()
	defer DeferError(xLogFile.Close)

	InitFlags()

	if getFlagBool("quiet") {
		setQuietLog()
	}

	InitConfig()

	xAuthorize(
		"POST",
		xData["apiAuthorizeUrl"],
		xData["xClient"],
		xData["xSecret"],
	)

	if getFlagBool("debug") {
		xLog.Print("Received access token: " +
			xData["TokenType"] + " " +
			xData["AccessToken"])

	}

	var sendMoney TransferFundToDynamicAccountUserStruct
	sm := &sendMoney.TransferFundToDynamicAccountUser.Request
	sm.IssuerAccountNumber = xData["xIssuerID"]
	sm.FromAccountNumber = xData["xIssuerID"]
	sm.FromWalletID = xData["xDefaultWallet"]
	sm.Amount = getFlagString("amount")
	sm.Currency = getFlagString("currency")
	sm.RecipientFirstName = getFlagString("firstname")
	sm.RecipientLastName = getFlagString("lastname")
	sm.Description = getFlagString("description")
	sm.RecipientEmail = getFlagString("payee")

	tResp, err := xTransferDynamic(sendMoney)

	if nil != err {
		log.Fatal("TransferFundToDynamicAccountUser failed because: " + err.Error())
	}

	tr := &tResp.TransferFundToDynamicAccountUserResponse.TransferFundToDynamicAccountUserResult

	if !tr.OperationStatus.Success {
		xLog.Println("TransferFundToDynamicAccountUser failed because:" +
			tr.OperationStatus.Errors)
	}

	if getFlagBool("debug") {
		jsonData, err := json.MarshalIndent(tResp, "", "  ")
		if nil != err {
			log.Fatalf("could not unmarshal JSON response because %s\n", err.Error())
		}
		log.Printf("jsondata response: \n%s\n", string(jsonData))
	}

	if !getFlagBool("quiet") {
		fmt.Printf("\nSuccess! TransactionID %s for %s (%s transferred, %s fee) %s to recipient %s (%s %s %s)\n",
			tr.TransactionID, tr.TotalAmount, tr.Amount,
			tr.Fee, tr.Currency, tr.RecipientAccountNumber,
			sm.RecipientFirstName, sm.RecipientLastName, sm.RecipientEmail)
	}

}
