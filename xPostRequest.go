package main

import (
	"io"
	"io/ioutil"
	"net/http"
)

func XPostRequest(url string, payload io.Reader) (resp []byte, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if nil != err {
		return resp, err
	}
	s := xData["TokenType"] + " " + xData["AccessToken"]
	if GetFlagBool("debug") {
		xLog.Println(s)
	}
	req.Header.Add("Authorization", s)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if nil == err {
		defer DeferError(res.Body.Close)
		resp, err = ioutil.ReadAll(res.Body)
		setTimeoutExpires()
	}
	if GetFlagBool("debug") {
		xLog.Println(string(resp))
	}
	return resp, err
}
