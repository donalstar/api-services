package main

import (
	"fmt"
	"trustcloud/api"
	"trustcloud/check"
)

func main() {
	ids := loadFromDB()

	// for each biz provider, run against Experian
	check.ExperianBusinessCheckIds(ids)
}

func loadFromDB() []int64 {
	var ids []int64

	_, err := api.DbMap.Select(&ids, "select p.id from Provider p, CheckedProvider c "+
		"where p.is_business = 'Y' and p.source = 'THUMBTACK'"+
		"and p.id = c.provider_id limit 100")

	fmt.Println("Err ", err)
	return ids
}
