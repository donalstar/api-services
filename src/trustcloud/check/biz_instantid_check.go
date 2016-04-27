package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"trustcloud/api"
	"trustcloud/util"
)

type LNInstantIdResponse struct {
	Response InstantIdResponse `json:"response"`
}

type InstantIdResponse struct {
	Result InstantIdResult `json:"Result"`
}

type InstantIdResult struct {
	CompanyResults CompanyResults `json:"CompanyResults"`
}

type CompanyResults struct {
	BusinessId    string             `json:"BusinessId"`
	VerifiedInput BizIdVerifiedInput `json:"VerifiedInput"`
}

type BizIdVerifiedInput struct {
	Name    string  `json:"CompanyName"`
	Address Address `json:"Address"`
	Phone   string  `json:"Phone10"`
}

type bizInstantIdChecker struct {
}

func (db *bizInstantIdChecker) GetName() string {
	return "LN_BIZ_INST_ID"
}

func (db *bizInstantIdChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {

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
		"&type=instant"

	util.InfoLog.Println("Running LN_BIZ_INST_ID checker... - url = ", url)

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	//	body, _ := ioutil.ReadFile(util.BaseDir + "/test/" + "ln_biz_id_resp.json")

	//		body, _ := ioutil.ReadFile(util.BaseDir + "/test/" + "ln_biz_instant_id2.json")

	var data LNInstantIdResponse

	err := json.Unmarshal(body, &data)

	util.InfoLog.Println("Resp ", string(body))

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	got_id := data.Response.Result.CompanyResults.BusinessId != "000000000000"

	checkedProvider.LnBizInstidMatch = "N"
	checkedProvider.LnBizInstidMatchedPhone = "N"
	checkedProvider.LnBizInstidMatchedName = "N"
	checkedProvider.LnBizInstidMatchedZip = "N"

	if got_id == true {
		util.InfoLog.Println("Got a match")

		checkedProvider.LnBizInstidMatch = "Y"

		if data.Response.Result.CompanyResults.VerifiedInput.Name != "" {
			checkedProvider.LnBizInstidMatchedName = "Y"
			util.InfoLog.Println("Co name match")
		}

		if data.Response.Result.CompanyResults.VerifiedInput.Phone != "" {
			checkedProvider.LnBizInstidMatchedPhone = "Y"
			util.InfoLog.Println("phone match")

		}

		if data.Response.Result.CompanyResults.VerifiedInput.Address.Zip != "" {
			checkedProvider.LnBizInstidMatchedZip = "Y"

			util.InfoLog.Println("zip match")
		}
	}
}
