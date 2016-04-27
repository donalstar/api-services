package check

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"trustcloud/api"
	"trustcloud/util"
)

type LNIndividualResponse struct {
	Response IndividualResponse `json:"response"`
}

type IndividualResponse struct {
	Header Header `json:"Header"`
	Result Result `json:"Result"`
}

type Header struct {
	Status        int    `json:"Status"`
	TransactionId string `json:"TransactionId"`
}

type Result struct {
	VerifiedInput                  VerifiedInput `json:"VerifiedInput"`
	ComprehensiveVerificationIndex int           `json:"ComprehensiveVerificationIndex"`
}

type VerifiedInput struct {
	Name  FullName `json:"Name"`
	Phone string   `json:"HomePhone"`
	Zip   string   `json:"Zip5"`
}

type FullName struct {
	First string `json:"First"`
	Last  string `json:"Last"`
}

type idChecker struct {
}

func (db *idChecker) GetName() string {
	return "LN_IND_ID"
}

func (db *idChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	body := getResponse(provider)

	//	body := getTestResponse()

	var data LNIndividualResponse

	util.InfoLog.Println("LN Response ", string(body))

	checkedProvider.LnIndivRawJson = sql.NullString{String: string(body), Valid: true}

	err := json.Unmarshal(body, &data)

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	if data.Response.Result.ComprehensiveVerificationIndex == 0 {
		util.InfoLog.Println("No match for provider ", provider.OwnerName, "\n")
		checkedProvider.LnIndivMatch = "N"
	}

	if data.Response.Result.ComprehensiveVerificationIndex >= 20 {
		util.InfoLog.Println("Matched provider ", provider.OwnerName, "\n")
		checkedProvider.LnIndivMatch = "Y"

		if len(data.Response.Result.VerifiedInput.Phone) != 0 {
			checkedProvider.LnIndivMatchedPhone = "Y"
		}

		if len(data.Response.Result.VerifiedInput.Zip) != 0 {
			checkedProvider.LnIndivMatchedZip = "Y"
		}
	}
}

func getResponse(provider *api.Provider) []byte {
	user := util.Configuration.General.Connector["lnIndividual"].Id
	password := util.Configuration.General.Connector["lnIndividual"].Password
	baseUrl := util.Configuration.General.Connector["lnIndividual"].BaseUrl

	name := strings.Split(provider.OwnerName, " ")

	if len(provider.OwnerName) == 0 {
		name = strings.Split(provider.Name, " ")
	}

	util.InfoLog.Println("Name ", name)

	fname := name[0]

	lname := ""

	if len(name) == 2 {
		lname = name[1]
	}

	util.InfoLog.Println("fname ", fname, " lname", lname, " phone ", provider.Phone, " zip ", provider.Zip)

	var url = baseUrl +
		"?firstName=" + url.QueryEscape(fname) +
		"&lastName=" + url.QueryEscape(lname) +
		"&zip=" + provider.Zip +
		"&phone=" + url.QueryEscape(strconv.FormatInt(provider.Phone, 10)) +
		"&user=" + user +
		"&password=" + password

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	util.InfoLog.Println("Resp ", string(body))
	return body
}

func getTestResponse() []byte {
	fileName := util.BaseDir + "/test/" + "ln_individual2.json"
	//	fileName := util.BaseDir + "/test/" + "ln_individual_bad_phone.json"

	body, _ := ioutil.ReadFile(fileName)

	return body
}
