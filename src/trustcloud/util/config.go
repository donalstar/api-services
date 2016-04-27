package util

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	Environment EnvConfig
	General     GeneralConfig
}

type EnvConfig struct {
	Http struct {
		Forwarder string
	}

	Database struct {
		Name     string
		Legacy   string
		User     string
		Password string
		Host     string
	}

	Mail struct {
		User     string
		Password string
		Server   string
		Port     int
		From     string
		TestMode string
	}

	Payment struct {
		User      string
		Password  string
		Signature string
		TestMode  string
	}
}

type GeneralConfig struct {
	Env struct {
		ListenPort string
	}

	Servers struct {
		Dev  string
		Prod string
	}

	Guarantee struct {
		Name            string
		PricePercentage int
		PriceRoundUp    int
		Currency        string
		SupportPhone    string
	}

	MailTemplate map[string]*struct {
		Name    string
		Subject string
	}

	Connector map[string]*struct {
		Id       string
		Password string
		BaseUrl  string
	}
}

// The one and only config instance.
var Configuration Config

var BaseDir string

func init() {

	SetBaseDirectory()

	configDir := "./"

	if BaseDir != "." {
		configDir = BaseDir + "/src/trustcloud/"
	}

	generalConfig := GeneralConfig{}
	err := gcfg.ReadFileInto(&generalConfig, configDir+"/config/trustcloud.gcfg")

	if err != nil {
		fmt.Println("Failed to parse config data: ", err)
	}

	envConfig := EnvConfig{}

	err = gcfg.ReadFileInto(&envConfig, configDir+"config/env.gcfg")

	if err != nil {
		fmt.Println("Failed to locate env.gcfg... defaulting to dev.gcfg")

		err = gcfg.ReadFileInto(&envConfig, configDir+"config/dev.gcfg")

		if err != nil {
			fmt.Println("Failed to parse config data: ", err)
		}
	}

	Configuration = Config{Environment: envConfig, General: generalConfig}
}

func IsSet(name string) bool {
	length := len(strings.TrimSpace(name))

	return (length != 0)
}

func SetBaseDirectory() {
	flag.StringVar(&BaseDir, "bd", ".", "base directory")

	flag.Parse()
}
