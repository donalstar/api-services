package check

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"trustcloud/api"
	"trustcloud/util"
)

type FullContact struct {
	Code           int             `json: "code"`
	Message        string          `json: "message"`
	ContactInfo    ContactInfo     `json: "contactInfo"`
	SocialProfiles []SocialProfile `json: "socialProfiles"`
}

type ContactInfo struct {
	Websites []Website `json: "websites"`
}

type Website struct {
	Url string `json: "url"`
}

type SocialProfile struct {
	Type             string          `json: "type"`
	Url              string          `json: "url"`
	Username         string          `json: "username"`
	Id               string          `json: "id"`
	TwitterFollowers json.RawMessage `json:"followers"`
}

type fullContactChecker struct {
}

func (db *fullContactChecker) GetName() string {
	return "FULL_CONTACT_IND"
}

func (db *fullContactChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	key := util.Configuration.General.Connector["fullcontact"].Id
	baseUrl := util.Configuration.General.Connector["fullcontact"].BaseUrl

	// NB: Only works for individuals, not businesses!
	url := baseUrl + "?email=" + provider.Email + "&apiKey=" + key + "&prettyPrint=true"

	util.InfoLog.Println("Use URL: ", url)

	res, err := http.Get(url)

	if err != nil {
		util.ErrorLog.Println("FullContact error for provider ", provider.Id, " : ", err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var data FullContact

	err = json.Unmarshal(body, &data)

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	if len(data.ContactInfo.Websites) > 0 {
		util.InfoLog.Println("Has a website: ", data.ContactInfo.Websites[0])

		checkedProvider.UserWebsite = data.ContactInfo.Websites[0].Url
	}

	util.InfoLog.Println(data.SocialProfiles)

	if len(data.SocialProfiles) > 0 {
		util.InfoLog.Println("Has social media profiles #: ", len(data.SocialProfiles))

		for _, profile := range data.SocialProfiles {
			if profile.Type == "facebook" {
				util.InfoLog.Println("Facebook Page -  ", profile.Id)
			}

			if profile.Type == "twitter" {
				checkedProvider.TwitterId = profile.Username

				if len(profile.TwitterFollowers) > 0 {
					checkedProvider.TwitterFollowers, _ = strconv.Atoi(string(profile.TwitterFollowers))

					util.InfoLog.Println("Twitter followers::: -  ", checkedProvider.TwitterFollowers)
				}
			}

			if profile.Type == "linkedin" {
				checkedProvider.LinkedinId = "Y"
			}
		}
	}
}
