package check

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"trustcloud/api"
	"trustcloud/util"
)

/*
INFO: 2015/04/13 13:43:23 full_contact_company_check.go:129: body :  {
  "status" : 404,
  "message" : "Searched within last 14 days. No results found for this Id.",
  "requestId" : "876a289c-d4f4-44d5-bb22-e4c01de83dbc"
}

*/
type FullContactCompany struct {
	Status         int                    `json: "status"`
	Organization   Organization           `json: "organization"`
	SocialProfiles []CompanySocialProfile `json: "socialProfiles"`
}

type Organization struct {
	Name string `json: "name"`
}

type CompanySocialProfile struct {
	Type             string          `json:"typeId"`
	Url              string          `json: "url"`
	Username         string          `json: "username"`
	Id               string          `json: "id"`
	TwitterFollowers json.RawMessage `json:"followers"`
}

type fullContactCompanyChecker struct {
}

func (db *fullContactCompanyChecker) GetName() string {
	return "FULL_CONTACT_BIZ"
}

func (db *fullContactCompanyChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	retryCount := 0

	website := provider.Website

	if checkedProvider.FcBizMatchedOrg == "Y" {
		util.InfoLog.Println("Already completed for id ", provider.Id, " - skip check")
	}

	checkedProvider.FcBizMatch = "N"
	checkedProvider.FcBizMatchedOrg = "N"
	checkedProvider.FcBizFbPage = sql.NullString{Valid: false}
	checkedProvider.FcBizTwitterId = sql.NullString{Valid: false}
	checkedProvider.FcBizTwitterFollowers = 0
	checkedProvider.FcBizLinkedinMatch = "N"

	if website.Valid == true && (strings.Contains(website.String, "facebook") == false) {
		data, err := doCheck(provider, checkedProvider, &retryCount)

		if err != nil {
			fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		}

		if data == nil {
			util.ErrorLog.Println("Failed to get result")
			return
		}

		for data.Status == 202 && retryCount < 4 {
			util.InfoLog.Println("FullContact search queued -- need to re-try...")

			time.Sleep(60000 * time.Millisecond)

			data, err = doCheck(provider, checkedProvider, &retryCount)
		}

		if data.Status == 200 {
			checkedProvider.FcBizMatch = "Y"

			if data.Organization.Name != "" {
				util.InfoLog.Println("Matched an organization: ", data.Organization.Name)

				checkedProvider.FcBizMatchedOrg = "Y"

				if strings.ToLower(provider.Name) == strings.ToLower(data.Organization.Name) {
					util.InfoLog.Println("Provider name matched")
				}

				if len(data.SocialProfiles) > 0 {
					for _, profile := range data.SocialProfiles {
						if profile.Type == "facebook" {
							util.InfoLog.Println("Facebook Page -  ", profile.Id)

							checkedProvider.FcBizFbPage = sql.NullString{String: profile.Id, Valid: true}
						}

						if profile.Type == "twitter" {
							checkedProvider.TwitterId = profile.Username

							checkedProvider.FcBizTwitterId = sql.NullString{String: profile.Username, Valid: true}

							if len(profile.TwitterFollowers) > 0 {
								checkedProvider.FcBizTwitterFollowers, _ = strconv.Atoi(string(profile.TwitterFollowers))

								util.InfoLog.Println("Twitter followers::: -  ", checkedProvider.FcBizTwitterFollowers)
							}
						}

						if profile.Type == "linkedincompany" {
							fmt.Println("Found company on LinkedIn ")
							checkedProvider.FcBizLinkedinMatch = "Y"
						}
					}

				}
			}
		}
	}
}

func doCheck(provider *api.Provider, checkedProvider *api.CheckedProvider, count *int) (*FullContactCompany, error) {
	util.InfoLog.Println("CheckFullContact for company ", provider.Id)
	//
	*count = *count + 1

	key := util.Configuration.General.Connector["fullcontactCompany"].Id
	baseUrl := util.Configuration.General.Connector["fullcontactCompany"].BaseUrl

	url := baseUrl + "?domain=" + provider.Website.String + "&apiKey=" + key + "&prettyPrint=true"

	util.InfoLog.Println("Use URL: ", url)

	res, err := http.Get(url)

	if err != nil {
		util.ErrorLog.Println("FullContact error for provider ", provider.Id, " : ", err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	util.InfoLog.Println("body : ", string(body))

	//body, _ := ioutil.ReadFile(util.BaseDir + "/test/" + "fc_company_resp.json")

	var data FullContactCompany

	err = json.Unmarshal(body, &data)

	util.WarningLog.Println("After Check -- count = ", *count)

	return &data, err
}
