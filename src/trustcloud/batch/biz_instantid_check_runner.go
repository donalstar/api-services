package main

import (
	"fmt"
	"trustcloud/api"
	"trustcloud/check"
)

func main() {
	ids := loadFromDB()

	// for each non-biz provider, run against LN biz instantid check
	check.CheckBusinessIds(ids)
}

func loadFromDB() []int64 {
	var ids []int64

	_, err := api.DbMap.Select(&ids, "select p.id from Provider p, CheckedProvider c "+
		"where p.is_business = 'Y' and p.source = 'THUMBTACK'"+
		"and p.id = c.provider_id limit 100")

	fmt.Println("Err ", err)
	return ids
}
