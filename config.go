package main

import (
	"fmt"
	"github.com/go-ini/ini"
	pflag "github.com/spf13/pflag"
	"os"
	"path/filepath"
)

const INIFILE string = "xtrm.ini"

var iniFileName = INIFILE
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

var nFlags *pflag.FlagSet

func InitFlags() {

	// Example command line:
	// XTRMPAY --payee nathan@xtrm.com --currency USD --amount 4.53  ^
	//      --firstname Jean-Paul --lastname Dough --description "Money Test" --debug true
	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)

	nFlags.StringP("payee", "p", "nathan@xtrm.com", "email address of payee")
	nFlags.StringP("currency", "c", "USD", "Currency to pay")
	nFlags.StringP("amount", "a", "4.53", "Amount to pay")
	nFlags.StringP("firstname", "f", "John", "Payee first name")
	nFlags.StringP("lastname", "l", "Doe", "Payee first name")
	nFlags.StringP("description", "d", "Money Test to Friend", "XTRMCMD payment method")
	nFlags.Bool("debug", false, "enable debug output")
	nFlags.BoolP("quiet", "q", false, "suppress non-error output")
	nFlags.String("profile", ini.DefaultSection, "API Access Profile")
	nFlags.BoolP("help", "h", false, "Print this help message")
	nFlags.String("payfile", "", "CSV payment file (email, firstname, lastname, amount, currency")

	err := nFlags.Parse(os.Args[1:])

	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "\nerror parsing flags: %s\n", err.Error())
		os.Exit(1)
	}

	if getFlagBool("help") {
		_, thisCmd := filepath.Split(os.Args[0])
		fmt.Print("\n", "usage for ", thisCmd, ":\n", nFlags.FlagUsages(), "\n")
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

	ex, _ := os.Executable()
	exePath := filepath.Dir(ex)
	myDir, _ := os.Getwd()

	cfg, err = ini.Load(iniFileName)
	if nil != err {
		if getFlagBool("debug") {
			xLog.Println("attempted to load ", INIFILE)
		}
		iniFileName = filepath.Join(myDir, INIFILE)
		cfg, err = ini.Load(iniFileName)

		if nil != err {
			if getFlagBool("debug") {
				xLog.Println("attempted to load ", iniFileName)
			}
			iniFileName = filepath.Join(exePath, INIFILE)
			cfg, err = ini.Load(iniFileName)
			if nil != err {
				if getFlagBool("debug") {
					xLog.Println("attempted to load ", iniFileName)
				}
				xLog.Printf("%s\n\t%s\n",
					"Failed to read config file [ xtrm.ini ]  because: ",
					err.Error())
			}
		}
	}

	if getFlagBool("debug") {
		xLog.Printf("\n\tLoading inifile from %s\n\n", iniFileName)
	}

	// Flags.Profile defaults to the default section
	// the only value in the default section is currentSection
	// so unless the profile is overridden, the default profile
	// is the last profile used.
	xData["currentSection"] = getFlagString("profile")
	xSec := loadSection(xData["currentSection"])
	if ini.DefaultSection == xData["currentSection"] {
		loadKey(xSec, "currentSection", true)
		xSec = loadSection(xData["currentSection"])
	}

	for _, v := range requiredKeys {
		loadKey(xSec, v, true)
	}
	for _, v := range optionalKeys {
		loadKey(xSec, v, false)
	}

}

func writeCurrentSectionKeys() {
	currentSection := xData["currentSection"]
	xSection, err := cfg.GetSection(currentSection)

	if nil != err {
		xLog.Fatal("internal error: no configuration section [" +
			currentSection + "]")
	}

	// required keys should not be touched here
	for _, v := range optionalKeys {
		saveIniKey(xSection, v, xData[v])
	}

	// update the currentSection
	xSection, err = cfg.GetSection(ini.DefaultSection)
	if nil != err {
		xLog.Fatal("internal error: no default section [ " +
			ini.DefaultSection + " ]")
	}
	saveIniKey(xSection, "currentSection", currentSection)

	// update ini file here
	err = cfg.SaveTo(iniFileName)
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

func saveIniKey(xSection *ini.Section, key string, val string) {
	xSection.DeleteKey(key)
	val, ok := xData[key]
	if ok {
		_, err := xSection.NewKey(key, val)
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
	msg := fmt.Sprintf("\n%s%s%s%s%s%s%s%s%s%s%s%s%s%s\n",
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

	if getFlagBool("help") {
		_, _ = fmt.Print(msg)
	} else {
		xLog.Print(msg)
	}
}
