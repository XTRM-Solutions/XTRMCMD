package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	ClientID     string `json:"client_id"`
	Issued       string `json:".issued"`
	Expires      string `json:".expires"`
}

func xAuthorize(xmethod, xurl, xclient, xsecret string) (success bool) {

	payload := strings.NewReader("grant_type=password" +
		"&client_id=" + xclient +
		"&client_secret=" + xsecret)

	client := &http.Client{}
	req, err := http.NewRequest(xmethod, xurl, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if nil != err {
		fmt.Println(err.Error())
		return
	}

	defer DeferError(res.Body.Close)

	xbody, err := ioutil.ReadAll(res.Body)
	if nil != err {
		fmt.Println(err.Error())
		return
	}

	var tr tokenResponse
	err = json.Unmarshal(xbody, &tr)
	if nil != err {
		if *Flags.Debug {
			fmt.Println("Internal error: unmarshal error in TokenResponse: " + err.Error())
		}
		return false
	}
	/*
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		ClientID     string `json:"client_id"`
		Issued       string `json:".issued"`
		Expires      string `json:".expires"`
	*/

	xData["AccessToken"] = tr.AccessToken
	xData["TokenType"] = tr.TokenType
	xData["ExpiresIn"] = strconv.Itoa(tr.ExpiresIn)
	xData["RefreshToken"] = tr.RefreshToken
	xData["ClientID"] = tr.ClientID
	xData["Issued"] = tr.Issued
	xData["Expires"] = tr.Expires

	writeCurrentSectionKeys()

	return len(tr.AccessToken) > 0

}
