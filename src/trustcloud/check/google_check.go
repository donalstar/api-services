package check

import (
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

type GoogleCheckResponse struct {
	Search       GoogleSearch `json:"search"`
	Match        string       `json:"match"`
	MatchedPhone string       `json:"matched_phone"`
}

type GoogleSearch struct {
	Name  string `json:"name"`
	City  string `json:"city"`
	State string `json:"state"`
	Phone string `json:"phone"`
}

type googleChecker struct {
}

func (db *googleChecker) GetName() string {
	return "Google"
}

func (db *googleChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {

	baseUrl := util.Configuration.General.Connector["google"].BaseUrl

	url := baseUrl + "?name=" + url.QueryEscape(provider.Name) +
		"&city=" + url.QueryEscape(provider.City) +
		"&state=" + provider.State +
		"&phone=" + strconv.FormatInt(provider.Phone, 10)

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data GoogleCheckResponse

	err := json.Unmarshal(body, &data)

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	result := string(body)

	util.InfoLog.Println("Google search result : ", result)

	limitExceeded := strings.Contains(result, "503")

	util.InfoLog.Println(provider.Name, " Limit exceeded? ", limitExceeded)

	checkedProvider.GoogleMatch = data.Match
	checkedProvider.GooglePhoneMatch = data.MatchedPhone
}
