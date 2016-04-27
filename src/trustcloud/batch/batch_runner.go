package main

import (
	"trustcloud/api"
	"trustcloud/check"
	"trustcloud/util"
)

func main() {
	loadType := true

	var ids []int64

	if loadType == true {
		//	fileName := "sample_providers_list.csv"
		//		fileName := "ely_stokes.csv"
		//	fileName := "angies_list_1000.csv"
		//	fileName := "sample_thumbtack.csv"
		//	fileName := "thumbtack-provider-data.csv"

		// Patrick,McNeil,communitypest211@hotmail.com,5083614526,Community Pest Control,Hudson,MA

		//		fileName := "api-test-accounts.csv"
		fileName := "tt.csv"

		ids = loadFromFile(fileName)
	}

	if loadType == false {
		ids = loadFromDB()
	}

	// for each provider, run checks against FB, LN, FullContact, our in-house providers DB & Yelp
	check.AssessProviders(ids)
}

func loadFromFile(fileName string) []int64 {
	//	format := 1 // TC CSV Format
	format := 4 // Thumbtack CSV Format

	return check.Load(util.BaseDir+"/test/"+fileName, format)
}

func loadFromDB() []int64 {
	return check.GetProviderIdsFromDB(api.DataBase)

}
