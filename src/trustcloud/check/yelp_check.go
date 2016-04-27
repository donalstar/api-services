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

type YelpResponse struct {
	Businesses []Business `json:"businesses"`
}

type Business struct {
	ReviewCount   int     `json:"review_count"`
	AverageRating float64 `json:"avg_rating"`
}

type yelpChecker struct {
}

func (db *yelpChecker) GetName() string {
	return "YELP"
}

func (db *yelpChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {

	ywsid := util.Configuration.General.Connector["yelp"].Id
	baseUrl := util.Configuration.General.Connector["yelp"].BaseUrl

	matchCount := 0

	if provider.Phone != 0 {
		url := baseUrl + "?phone=" + strconv.FormatInt(provider.Phone, 10) + "&ywsid=" + ywsid

		res, _ := http.Get(url)

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var data YelpResponse

		err := json.Unmarshal(body, &data)

		if err != nil {
			fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		}

		matchCount = len(data.Businesses)

		if matchCount > 0 {
			util.InfoLog.Println("Yelp Match! - # of ratings: ", data.Businesses[0].ReviewCount,
				" Average Rating: ", data.Businesses[0].AverageRating)

			checkedProvider.YelpRating = data.Businesses[0].AverageRating
			checkedProvider.YelpReviewCount = data.Businesses[0].ReviewCount
		}
	}
}
