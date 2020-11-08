package main

import (
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
	PayeeEmail  *string
	Currency    *string
	Amount      *string
	FirstName   *string
	LastName    *string
	Description *string
	Debug       *bool
	Profile     *string
	Help        *bool
}

func InitFlags() {

	// Example command line:
	// XTRMPAY --payee nathan@xtrm.com --currency USD --amount 4.53  \
	//      --firstname Jean-Paul --lastname Dough --description "Money Test" --debug true

	Flags.PayeeEmail = flag.StringP("payee", "p", "nathan@xtrm.com", "email address of payee")
	Flags.Currency = flag.StringP("currency", "c", "USD", "Currency to pay")
	Flags.Amount = flag.StringP("amount", "a", "4.53", "Amount to pay")
	Flags.FirstName = flag.StringP("firstname", "f", "Jean-Paul", "Payee first name")
	Flags.LastName = flag.StringP("lastname", "l", "Dough", "Payee first name")
	Flags.Description = flag.StringP("description", "d", "Money Test to Friend", "Description of Money Transfer")
	Flags.Debug = flag.Bool("debug", true, "enable debug output")
	Flags.Profile = flag.String("profile", ini.DefaultSection, "API Access Profile")
	Flags.Help = flag.BoolP("help", "h", false, "Print this help message")

	flag.Parse()

	if *Flags.Help {
		flag.Usage()
		msgRequiredIniKeys()
		os.Exit(0)
	}
}

func loadSection(profile string) (section *ini.Section) {
	var err error
	section, err = cfg.GetSection(profile)
	if nil != err {
		xLog.Fatal("could not fetch .INI file profile / section [ " +
			xData["currentSection"] + " ] because: + " +
			err.Error())
	}
	return section
}

func InitConfig() {
	var err error
	cfg, err = ini.Load("xtrm.ini")
	if nil != err {
		xLog.Fatalf("%s\n\t%s\n",
			"Failed to read config file [ xtrm.ini ]  because: ",
			err.Error())
	}

	xData["currentSection"] = *Flags.Profile
	xsec := loadSection(*Flags.Profile)
	if ini.DefaultSection == *Flags.Profile {
		loadKey(xsec, "currentSection", true)
		xsec = loadSection(xData["currentSection"])
	}

	for _, v := range requiredKeys {
		loadKey(xsec, v, true)
	}
	for _, v := range optionalKeys {
		loadKey(xsec, v, false)
	}

}

func writeCurrentSectionKeys() {
	currentSection := xData["currentSection"]
	xsec, err := cfg.GetSection(currentSection)

	if nil != err {
		xLog.Fatal("internal error: no configuration section [" +
			currentSection + "]")
	}

	// required keys should not be touched here
	for _, v := range optionalKeys {
		saveIniKey(xsec, v, xData[v])
	}

	// update the currentSection
	xsec, err = cfg.GetSection(ini.DefaultSection)
	if nil != err {
		xLog.Fatal("internal error: no default section [ " +
			ini.DefaultSection + " ]")
	}
	saveIniKey(xsec, "currentSection", currentSection)

	// update ini file here
	err = cfg.SaveTo("xtrm.ini")
	if nil != err {
		xLog.Fatal("Internal error: failed to write config file [ xtrm.ini ] because: \n\t" + err.Error())
	}
}

func loadKey(section *ini.Section, key string, required bool) {
	if required && !section.HasKey(key) {
		msgRequiredIniKeys()
		xLog.Fatal("missing required key [" + key + " ] in section [ " +
			section.Name() + " ]")
	}
	xData[key] = section.Key(key).String()
}

func saveIniKey(xsec *ini.Section, key string, val string) {
	xsec.DeleteKey(key)
	val, ok := xData[key]
	if ok {
		_, err := xsec.NewKey(key, val)
		if nil != err {
			xLog.Fatalf("%s%s%s%s%s%s",
				"Could not set key [ ",
				val,
				"] to value [ ",
				val,
				" ] because:\n\t",
				err.Error())
		}
	}
}

func msgRequiredIniKeys() {
	xLog.Printf("\n%s%s%s%s%s%s%s%s%s%s%s%s%s%s\n",
		"This program REQUIRES some initialization keys in the file XTRM.INI\n",
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
		"Please note all keys and values are CASE SENSITIVE as token, secret, and URLs may be case sensitive")
}
