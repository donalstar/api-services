package check

import (
	"trustcloud/api"
	"trustcloud/util"
)

// Get Trustscore
func CalculateTrustscore(checkedProvider *api.CheckedProvider) int {
	//	return barebonesAlgorithm(checkedProvider)

	return thumbtackAlgorithm(checkedProvider)
}

func barebonesAlgorithm(checkedProvider *api.CheckedProvider) int {
	util.InfoLog.Println("Calculate trustscore (bare-bones):")

	/*
		bare-bones algorithm:

		- threshold = 500
		- in provider DB = 150
		- yelp  = 100
		- 	rating = 100 x (rating/5)
		- 	reviews = >10, 50, 1 - 10, 25
		- full contact : 100
		- 	website : 40
		- 	twitter id : 40
		- 	twitter followers : >20, 20, 1-20, 10
		- facebook : 150
		- 	page url :60
		- 	matched phone :60
		- 	likes: >100, 30, 1-100, 15
	*/

	score := 500

	util.InfoLog.Println("base: 500")

	if checkedProvider.InMasterList == "Y" {
		score = score + 150
		util.InfoLog.Println("checked provider: add 150")
	}

	// yelp
	yelpRatingPoints := (checkedProvider.YelpRating / float64(5)) * float64(100)

	util.InfoLog.Println("yelp: add ", yelpRatingPoints)

	score = score + int(yelpRatingPoints)

	if checkedProvider.YelpReviewCount >= 10 {
		score = score + 50

		util.InfoLog.Println("yelp > 10 reviews: add 25")
	}

	if checkedProvider.YelpReviewCount > 0 {
		score = score + 25
	}

	// full contact
	if checkedProvider.UserWebsite != "" {
		score = score + 40

		util.InfoLog.Println("has user website: add 40")
	}

	if checkedProvider.TwitterId != "" {
		score = score + 40

		util.InfoLog.Println("has twitter id: add 40")
	}

	if checkedProvider.TwitterFollowers > 20 {
		score = score + 20

		util.InfoLog.Println("has >20 twitter followers: add 20")
	}

	if checkedProvider.TwitterFollowers > 0 {
		score = score + 10

		util.InfoLog.Println("has >0 twitter followers: add 10")
	}

	// facebook
	if checkedProvider.FbPageUrl != "" {
		score = score + 60

		util.InfoLog.Println("has FB page: add 60")
	}

	if checkedProvider.FbMatchedPhone == "Y" {
		score = score + 60

		util.InfoLog.Println("has FB matched phone: add 60")
	}

	if checkedProvider.FbLikes >= 100 {
		score = score + 30

		util.InfoLog.Println("has >100 FB likes: add 30")
	}

	if checkedProvider.FbLikes > 0 {
		score = score + 15

		util.InfoLog.Println("has >0 FB likes: add 15")
	}

	return score
}

func thumbtackAlgorithm(checkedProvider *api.CheckedProvider) int {
	util.InfoLog.Println("Calculate trustscore (thumbtack algorithm):")

	/*
			500
			+FB Likes Flag*100
			+Yelp Flag*100
			+LN Flag*100
			+Yelp Reviews Flag*25
			+Website*10+Twitter Flag*50+Twitter Follows Flag*25

		All flags are 0 or 1 depending if yes there is an account / verification or not, except Yelp Reviews
		is more than 3 reviews and Twitter Follows Flag is more than 30.

	*/
	score := 500

	util.InfoLog.Println("base: 500")

	util.InfoLog.Println("checkedProvider.LnBizInstidMatch? ", checkedProvider.LnBizInstidMatch)

	if checkedProvider.LnBizMatch == "Y" || checkedProvider.LnIndivMatch == "Y" ||
		checkedProvider.LnBizInstidMatch == "Y" {
		score = score + 100
		util.InfoLog.Println("has Lexis Nexis match: add 100")
	}

	if checkedProvider.FbLikes > 0 {
		score = score + 50

		util.InfoLog.Println("has >0 FB likes: add 50")
	}

	if checkedProvider.FbPageUrl != "" {
		score = score + 50

		util.InfoLog.Println("has FB Page: add 50")
	}

	if checkedProvider.FbMatchedPhone != "" {
		score = score + 50

		util.InfoLog.Println("mached FB Phone: add 50")
	}

	// yelp
	if checkedProvider.YelpRating > 3 {
		score = score + 100
		util.InfoLog.Println(">has Yelp rating > 3: add 100")
	}

	if checkedProvider.YelpRating > 2 && checkedProvider.YelpRating <= 3 {
		score = score + 100
		util.InfoLog.Println(">has Yelp rating 2-3: add 50")
	}

	if checkedProvider.YelpReviewCount > 3 {
		score = score + 25
		util.InfoLog.Println(">3 Yelp reviews: add 25")
	}

	if checkedProvider.YelpReviewCount > 0 && checkedProvider.YelpReviewCount <= 3 {
		score = score + 25
		util.InfoLog.Println(">0 Yelp reviews: add 10")
	}

	// full contact
	if checkedProvider.UserWebsite != "" {
		score = score + 10

		util.InfoLog.Println("has user website: add 10")
	}

	if checkedProvider.TwitterId != "" {
		score = score + 50

		util.InfoLog.Println("has twitter id: add 50")
	}

	if checkedProvider.TwitterFollowers > 30 {
		score = score + 25

		util.InfoLog.Println("has >30 twitter followers: add 25")
	}

	if checkedProvider.TwitterFollowers > 5 {
		score = score + 10

		util.InfoLog.Println("has >5 twitter followers: add 10")
	}

	if checkedProvider.LinkedinId == "Y" {
		score = score + 50

		util.InfoLog.Println("has LinkedIn Id: add 50")
	}

	return score
}
