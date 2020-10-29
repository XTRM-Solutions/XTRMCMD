package main

import (
	"bytes"
	"encoding/json"
)

// need to stack this stuff to marshall / unmarshall correctly
// JSON unmarshaller uses the tag information to map structs to fields by reflection
type Request struct {
	IssuerAccountNumber string `json:"IssuerAccountNumber"`
	FromAccountNumber   string `json:"FromAccountNumber"`
	FromWalletID        string `json:"FromWalletID"`
	RecipientFirstName  string `json:"RecipientFirstName"`
	RecipientLastName   string `json:"RecipientLastName"`
	RecipientEmail      string `json:"RecipientEmail"`
	Currency            string `json:"Currency"`
	Amount              string `json:"Amount"`
	Description         string `json:"Description"`
}

type TransferFundToDynamicAccountUser struct {
	Request Request `json:"Request"`
}

type TransferFundToDynamicAccountUserStruct struct {
	TransferFundToDynamicAccountUser TransferFundToDynamicAccountUser `json:"TransferFundToDynamicAccountUser"`
}

type OperationStatus struct {
	Success bool        `json:"Success"`
	Errors  interface{} `json:"Errors"`
} // `json:"OperationStatus"`

type TransferFundToDynamicAccountUserResult struct {
	TransactionID          string          `json:"TransactionID"`
	Amount                 string          `json:"Amount"`
	Fee                    string          `json:"Fee"`
	TotalAmount            string          `json:"TotalAmount"`
	Currency               string          `json:"Currency"`
	RecipientAccountNumber string          `json:"RecipientAccountNumber"`
	OperationStatus        OperationStatus `json:"OperationStatus"`
} // `json:"TransferFundToDynamicAccountUserResult"`

type TransferFundToDynamicAccountUserResponse struct {
	TransferFundToDynamicAccountUserResult TransferFundToDynamicAccountUserResult `json:"TransferFundToDynamicAccountUserResult"`
} // `json:"TransferFundToDynamicAccountUserResponse"`

type TransferFundToDynamicAccountUserResponseStruct struct {
	TransferFundToDynamicAccountUserResponse TransferFundToDynamicAccountUserResponse `json:"TransferFundToDynamicAccountUserResponse"`
}

func xTransferDynamic(reqData TransferFundToDynamicAccountUserStruct) (resp TransferFundToDynamicAccountUserResponseStruct, err error) {
	url := xData["xUrl"] + "/Fund/TransferFundDynamicAccountCreateUser"

	jsonData, err := json.Marshal(reqData)
	if nil != err {
		return resp, err
	}
	payload := bytes.NewReader(jsonData)
	xbody, err := XPostRequest(url, payload)
	if nil != err {
		return resp, err
	}
	return resp, json.Unmarshal(xbody, &resp)
}
