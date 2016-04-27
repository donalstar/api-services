package main

import (
	"fmt"
	"trustcloud/api"
	"trustcloud/check"
)

func main() {
	ids := loadFromDB()

	// for each non-biz provider, run against LN biz instantid check
	check.CheckWebsites(ids)
}

func loadFromDB() []int64 {
	var ids []int64

	_, err := api.DbMap.Select(&ids, "select p.id from Provider p, CheckedProvider c "+
		" where p.source = 'THUMBTACK' and website <> ''"+
		" and p.id = c.provider_id")

	fmt.Println("Err ", err)
	return ids
}
