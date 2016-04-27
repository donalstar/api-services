package main

import (
	"trustcloud/api"
	"trustcloud/check"
)

func main() {
	ids := loadFromDB()

	// for each non-biz provider, run against LN biz check
	check.CheckIndividuals(ids)
}

func loadFromDB() []int64 {
	return check.GetProviderIdsFromDB(api.DataBase)
}
