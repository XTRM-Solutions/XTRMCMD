package main

import (
	"fmt"
	"github.com/go-ini/ini"
	flag "github.com/spf13/pflag"
	"os"
)

var cfg *ini.File

var requiredKeys = [...]string{
	"apiAuthorizeUrl",
	"xIssuerID",
	"xClient",
	"xSecret",
	"xUrl",
	"xDefaultWallet",
}

var optionalKeys = [...]string{
	"AccessToken",
	"TokenType",
	"ExpiresIn",
	"RefreshToken",
	"ClientID",
	"Issued",
	"Expires",
}

var xData = map[string]string{
	"currentSection": ini.DefaultSection,
}

var Flags struct {
	Payee       *string
	Currency    *string
	Amount      *string
	FirstName   *string
	LastName    *string
	Description *string
	Debug       *bool
}

func InitFlags() {

	// Example command line:
	// XTRMPAY --payee nathan@xtrm.com --currency USD --amount 4.53  \
	//      --firstname Jean-Paul --lastname Dough --description "Money Test" --debug true

	Flags.Payee = flag.StringP("payee", "p", "nathan@xtrm.com", "email address of payee")
	Flags.Currency = flag.StringP("currency", "c", "USD", "Currency to pay")
	Flags.Amount = flag.StringP("amount", "a", "4.53", "Amount to pay")
	Flags.FirstName = flag.StringP("firstname", "f", "Jean-Paul", "Payee first name")
	Flags.LastName = flag.StringP("lastname", "l", "Dough", "Payee first name")
	Flags.Description = flag.StringP("description", "d", "Money Test to Friend", "Description of Money Transfer")
	Flags.Debug = flag.Bool("debug", true, "enable debug output")

	flag.Parse()
}

func InitConfig() {
	var err error
	cfg, err = ini.Load("xtrm.ini")
	if nil != err {
		fmt.Printf("Fail to read file [xtrm.ini] : %s\n", err.Error())
		os.Exit(1)
	}

	xData["currentSection"] = cfg.Section(ini.DefaultSection).Key("currentSection").String()

	xsec, err := cfg.GetSection(xData["currentSection"])

	for _, v := range requiredKeys {
		loadKey(xsec, true, v)
	}
	for _, v := range optionalKeys {
		loadKey(xsec, false, v)
	}

}

func writeCurrentSectionKeys() {
	currentSection := xData["currentSection"]
	xsec, err := cfg.GetSection(currentSection)

	if nil != err {
		fmt.Println("internal error: no configuration section [" +
			currentSection + "]")
		os.Exit(3)
	}

	for _, v := range optionalKeys {
		xsec.DeleteKey(v)
		val, ok := xData[v]
		if ok {
			_, err = xsec.NewKey(v, val)
			if nil != err {
				fmt.Println("Internal error: could not set " + "" +
					"ini file key [" + v + "] to [" + val + "]")
			}
		}
	}

	err = cfg.SaveTo("xtrm.ini")
	if nil != err {
		fmt.Println("Internal error: failed to write config file [xtrm.ini] " + err.Error())
		os.Exit(4)
	}
}

func loadKey(section *ini.Section, required bool, key string) {
	if required && !section.HasKey(key) {
		msgRequiredIniKeys()
		os.Exit(2)
	}
	xData[key] = section.Key(key).String()
}

func msgRequiredIniKeys() {
	_, _ = fmt.Fprintf(os.Stderr, "\n%s%s%s%s%s%s%s%s%s%s%s%s%s%s\n",
		"XTRM requires some initialization keys in the file XTRM.INI\n",
		"an initial file looks something like: (minimal required file)\n\n",
		"\t[DEFAULT]\n",
		"\tcurrentSection=initial\n\n",
		"\t[initial]\n",
		"\tapiAuthorizeUrl=https://zodmo.xapi.xtrm.com/oAuth/token\n",
		"\txIssuerID=SPN99999999\n",
		"\txClient=9999999_API_User\n",
		"\txSecret=gTv/g5LNOdHRkxlo/bjYxWo6YUXZWTkhjN04RnvDGls%3D\n",
		"\txUrl=zodmo.xapi.xtrm.com/API/V4\n",
		"\txDefaultWallet=123456\n",
		"\nPlease ensure this file exists with the minimum required keys in the XTRM command directory\n",
		"Please substitute in the correct values from the API integration page in the console application\n",
		"Please note all keys and values are CASE SENSITIVE")
}
