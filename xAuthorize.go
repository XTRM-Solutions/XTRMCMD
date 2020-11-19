package xtrmcmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
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

func isTokenActive(duration time.Duration) (active bool) {

	// Do we already have an access token good for at least 2 hours?
	if "" != xData["AccessToken"] {
		// is it current?
		expires := xData["Expires"]
		if "" != expires {
			// XTRM time example: "Wed, 28 Oct 2020 20:15:16 GMT"
			xtrmTimeFormat := "Mon, 02 Jan 2006 15:04:05 MST"

			timeExpires, err := time.Parse(xtrmTimeFormat, expires)
			if nil != err {
				xLog.Fatal("Internal error: could not parse time [ " +
					expires + "] as format [ " +
					xtrmTimeFormat + " ]\n\tbecause\n" +
					err.Error())
			}
			if timeExpires.After(time.Now().Add(duration)) {
				return true
			}
		}
	}
	return false
}

func xAuthorize(xmethod, xurl, xclient, xsecret string) (success bool) {

	if isTokenActive((2 * time.Hour)) {
		return true
	}

	// otherwise need to authorize
	// don't see any point in using token refresh
	payload := strings.NewReader("grant_type=password" +
		"&client_id=" + xclient +
		"&client_secret=" + xsecret)

	client := &http.Client{}
	req, err := http.NewRequest(xmethod, xurl, payload)

	if err != nil {
		xLog.Fatal(err.Error())
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if nil != err {
		xLog.Fatal(err.Error())
	}

	defer DeferError(res.Body.Close)

	xbody, err := ioutil.ReadAll(res.Body)
	if nil != err {
		xLog.Fatal(err.Error())
	}

	var tr tokenResponse
	err = json.Unmarshal(xbody, &tr)
	if nil != err {
		xLog.Fatal(err.Error())
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
