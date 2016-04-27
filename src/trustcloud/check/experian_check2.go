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

/*
check Experian - invoke local Java web service
*/
type ExperianCheckResponse struct {
	Result bool   `json:"result"`
	Bin    string `json:"BIN"`
}

type experianChecker2 struct {
}

func (db *experianChecker2) GetName() string {
	return "EXPERIAN_ID"
}

func (db *experianChecker2) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	baseUrl := util.Configuration.General.Connector["experian2"].BaseUrl

	/*
		http://localhost:4744/experian_check/hhgh/A%20to%20Z%20Painting%20Plus/Concord/NH/03301/6038564270
	*/
	url := baseUrl + "/" + url.QueryEscape(provider.Name) +
		"/" + url.QueryEscape(provider.Address1) +
		"/" + url.QueryEscape(provider.City) +
		"/" + provider.State +
		"/" + provider.Zip +
		"/" + strconv.FormatInt(provider.Phone, 10)

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data ExperianCheckResponse

	err := json.Unmarshal(body, &data)

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	result := string(body)

	util.InfoLog.Println("Experian2 search result : ", result)

	checkedProvider.ExperianMatch = "N"

	if data.Result == false {
		checkedProvider.ExperianMatch = "Y"
	}

	checkedProvider.ExperianBizId = sql.NullString{String: data.Bin, Valid: true}
}
