package check

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"trustcloud/api"
	"trustcloud/util"
)

type LNResponse struct {
	Response Response `json:"response"`
}

type Response struct {
	RecordCount int     `json:"RecordCount"`
	Records     Records `json:"Records"`
}

type Records struct {
	Record json.RawMessage `json:"Record"`
}

type Record struct {
	Address     Address `json:"Address"`
	Phone       string  `json:"Phone10"`
	CompanyName string  `json:"CompanyName"`
	BusinessId  string  `json:"BusinessId"`
}

type Address struct {
	StreetNumber string `json:"StreetNumber"`
	StreetName   string `json:"StreetName"`
	StreetSuffix string `json:"StreetSuffix"`
	City         string `json:"City"`
	State        string `json:"City"`
	Zip          string `json:"Zip5"`
}

type businessIdChecker struct {
}

func (db *businessIdChecker) GetName() string {
	return "LN_BIZ_ID"
}

func (db *businessIdChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	user := util.Configuration.General.Connector["ln"].Id
	password := util.Configuration.General.Connector["ln"].Password
	baseUrl := util.Configuration.General.Connector["ln"].BaseUrl

	var url = baseUrl +
		"?name=" + url.QueryEscape(provider.Name) +
		"&address1=" + url.QueryEscape(provider.Address1) +
		"&city=" + url.QueryEscape(provider.City) +
		"&state=" + provider.State +
		"&zip=" + provider.Zip +
		"&phone=" + url.QueryEscape(strconv.FormatInt(provider.Phone, 10)) +
		"&user=" + user +
		"&password=" + password +
		"&type=search"

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	//	body, _ := ioutil.ReadFile(util.BaseDir + "/test/" + "record.json")

	//	body, _ := ioutil.ReadFile(util.BaseDir + "/test/" + "res.json")

	var data LNResponse

	util.InfoLog.Println("LN Response ", string(body))

	checkedProvider.LnBizRawJson = sql.NullString{String: string(body), Valid: true}

	err := json.Unmarshal(body, &data)

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	var records []Record

	if data.Response.RecordCount == 1 {
		var record Record

		records = append(records, record)

		err = json.Unmarshal(data.Response.Records.Record, &records[0])
	}

	if data.Response.RecordCount > 1 {
		err = json.Unmarshal(data.Response.Records.Record, &records)
	}

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	matchCount := data.Response.RecordCount

	util.InfoLog.Println(" # of id matches ", matchCount)

	if matchCount > 0 {
		checkedProvider.LnBizMatch = "Y"

		for _, record := range records {
			phone, _ := strconv.ParseInt(string(record.Phone), 10, 64)

			if phone == provider.Phone {
				checkedProvider.LnBizMatchedPhone = "Y"
				break
			}
		}
	}
}
