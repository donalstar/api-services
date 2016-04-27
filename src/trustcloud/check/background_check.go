package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"trustcloud/api"
	"trustcloud/util"
)

type BGCResponse struct {
	Response BGCResponseAttributes `json:"@attributes"`
	Product  BGCProduct            `json:"product"`
	Status   BGCStatus             `json:"pwf"`
}

type BGCResponseAttributes struct {
	OrderId string `json:"orderId"`
}

type BGCProduct struct {
}

type BGCStatus struct {
	Status     string `json:"status"`
	Reason     string `json:"reason"`
	ReasonCode string `json:"reasonCode"`
}

type backgroundChecker struct {
}

func (db *backgroundChecker) GetName() string {
	return "BACKGROUND"
}

func (db *backgroundChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	user := util.Configuration.General.Connector["backgroundchecks"].Id
	password := util.Configuration.General.Connector["backgroundchecks"].Password
	baseUrl := util.Configuration.General.Connector["backgroundchecks"].BaseUrl

	name := strings.Split(provider.Name, " ")

	firstName := name[0]
	lastName := ""

	if len(name) > 1 {
		lastName = name[1]
	}

	var url = baseUrl +
		"?firstName=" + firstName +
		"&lastName=" + lastName +
		"&user=" + user +
		"&password=" + password

	layout := "2006/01/02"

	if !time.Time(provider.DateOfBirth).IsZero() {
		dateOfBirth := provider.DateOfBirth.Format(layout)

		dates := strings.Split(dateOfBirth, "/")

		url = url +
			"&day=" + dates[0] +
			"&month=" + dates[1] +
			"&year=" + dates[2]
	}

	util.InfoLog.Println("Running BACKGROUND CHECK -- URL = ", url)

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data BGCResponse

	err := json.Unmarshal(body, &data)

	util.InfoLog.Println("DATA ", data, " status ", data.Status)

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	if len(data.Response.OrderId) > 0 {
		util.InfoLog.Println("Background check completed successfully for ", provider.Name)

		checkedProvider.BgcStatus = "N"

		if data.Status.Status == "p" {
			checkedProvider.BgcStatus = "Y"
		}
	}
}
