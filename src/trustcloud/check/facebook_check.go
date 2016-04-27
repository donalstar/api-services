package check

import (
	"encoding/json"
	"fmt"
	fb "github.com/huandu/facebook"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"trustcloud/api"
	"trustcloud/util"
)

type FacebookAccessToken struct {
	Token string `json: "access_token"`
}

type FacebookPage struct {
	Id          string   `json: "id"`
	Likes       int      `json: "likes"`
	Location    Location `json: "location"`
	Phone       string   `json: "phone"`
	ParsedPhone string
}

type Location struct {
	Street string `json: "street"`
	City   string `json: "city"`
	State  string `json: "state"`
}

type facebookChecker struct {
}

var fbAccessToken string

func init() {
	fbAccessToken = GetAccessToken()
}

func (db *facebookChecker) GetName() string {
	return "FACEBOOK"
}

func (db *facebookChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	pageIds := GetPageMatches(provider)

	util.InfoLog.Println("FB check - provider name ", provider.Name)

	util.InfoLog.Println("Page IDs ", pageIds)

	for _, pageId := range pageIds {
		data := GetPage(pageId)

		// Get page, then search for biz phone to see if there's a match
		util.InfoLog.Println("Page Data ", data)

		if (provider.City == data.Location.City) && (provider.State == data.Location.State) {
			matchedPhone := "N"

			providerPhone := strconv.FormatInt(int64(provider.Phone), 10)

			if providerPhone == data.ParsedPhone {
				matchedPhone = "Y"
			}

			checkedProvider.FbMatchedPhone = matchedPhone

			checkedProvider.FbPageUrl = "http://graph.facebook.com/" + pageId

			checkedProvider.FbLikes = data.Likes
		}
	}
}

func GetPage(id string) FacebookPage {
	baseUrl := util.Configuration.General.Connector["facebook"].BaseUrl

	url := baseUrl + "/" + id

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data FacebookPage

	err := json.Unmarshal(body, &data)

	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	ph := util.ParsePhoneString(data.Phone)

	if ph != nil {
		data.ParsedPhone = *ph
	}

	return data
}

func GetPageMatches(provider *api.Provider) []string {
	queryString := provider.Name

	var pageIds []string

	GetPageIdsForQuery(url.QueryEscape(queryString), &pageIds)

	GetPageIdsForQuery(queryString, &pageIds)

	return pageIds
}

func GetPageIdsForQuery(queryString string, pageIds *[]string) {
	res2, _ := fb.Get("/search", fb.Params{
		"access_token": fbAccessToken,
		"type":         "page",
		"q":            queryString,
	})

	var items []fb.Result

	util.InfoLog.Println("FB q ", queryString)

	util.InfoLog.Println("FB result ", res2)

	err := res2.DecodeField("data", &items)

	if err != nil {
		util.InfoLog.Println("An error has happened ", err)
		return
	}

	p := *pageIds
	for _, v := range items {
		*pageIds = append(p, v["id"].(string))
	}
}

func GetAccessToken() string {
	appId := util.Configuration.General.Connector["facebook"].Id
	clientSecret := util.Configuration.General.Connector["facebook"].Password
	baseUrl := util.Configuration.General.Connector["facebook"].BaseUrl

	// get access token
	url := baseUrl + "/oauth/access_token?client_id=" +
		appId + "&client_secret=" + clientSecret + "&grant_type=client_credentials"

	res, _ := http.Get(url)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	pieces := strings.Split(string(body), "=")

	return pieces[1]
}
