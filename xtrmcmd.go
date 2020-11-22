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

	xAuthorize("POST",
		xData["apiAuthorizeUrl"],
		xData["xClient"],
		xData["xSecret"])

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
		log.Fatal(err.Error())
	}

	// this is just to lessen the amount of typing to access the object fields
	transferResult := &tResp.TransferFundToDynamicAccountUserResponse.TransferFundToDynamicAccountUserResult

	if !transferResult.OperationStatus.Success {
		xLog.Fatal("TransferFundToDynamicAccountUser failed because:" +
			transferResult.OperationStatus.Errors)
	}

	if getFlagBool("debug") {
		jsonData, err := json.MarshalIndent(tResp, "", "    ")
		if nil != err {
			log.Fatalf("could not unmarshal JSON response because %s\n", err.Error())
		}
		log.Printf("jsondata response: \n%s\n", string(jsonData))
	}

	if !getFlagBool("quiet") {

		fmt.Printf("\nSuccess! TransactionID %s for %s (%s transferred, %s fee) %s to recipient %s (%s %s %s)\n",
			transferResult.TransactionID, transferResult.TotalAmount, transferResult.Amount,
			transferResult.Fee, transferResult.Currency, transferResult.RecipientAccountNumber,
			sm.RecipientFirstName, sm.RecipientLastName,
			sm.RecipientEmail)
	}

}
