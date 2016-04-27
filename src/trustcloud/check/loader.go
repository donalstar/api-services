package check

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"trustcloud/api"
	"trustcloud/util"
)

/*
read lines from CSV file, and load as Provider records into the DB
*/
func Load(fileName string, format int) []int64 {
	db := api.DataBase

	csvfile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, each := range rawCSVdata {
		for i, val := range each {
			util.InfoLog.Println("next row ", val)
			each[i] = strings.Trim(val, "'")
		}
	}

	var ids []int64

	for _, each := range rawCSVdata {
		provider := &api.Provider{}

		switch format {
		case 1:
			provider = CreateProviderFromFormat1(each)
			break
		case 2:
			provider = CreateProviderFromFormat3(each)
			break
		case 4:
			provider = CreateProviderFromFormat4(each)
			break
		}

		providerId := ImportRecord(provider, db)

		ids = append(ids, providerId)
	}

	return ids
}

func CreateProviderFromFormat1(each []string) *api.Provider {
	phone, _ := strconv.ParseInt(strings.Trim(each[5], " "), 10, 64)

	provider := &api.Provider{
		Name:   each[0],
		Phone:  phone,
		Email:  each[4],
		City:   each[1],
		State:  each[2],
		Zip:    each[3],
		Source: "THUMBTACK",
	}

	return provider
}

func CreateProviderFromFormat2(each []string) *api.Provider {
	phone, _ := strconv.ParseInt(each[4], 10, 64)

	provider := &api.Provider{
		Name:      each[5],
		Phone:     phone,
		Email:     each[3],
		OwnerName: each[1] + " " + each[2],
		Address1:  each[6],
		City:      each[7],
		State:     each[8],
		Zip:       each[9],
		Source:    "THUMBTACK",
	}

	fmt.Println("Created provider ", provider)

	return provider
}

func CreateProviderFromFormat3(each []string) *api.Provider {
	phone, _ := strconv.ParseInt(each[4], 10, 64)

	provider := &api.Provider{
		Name:      each[7],
		Phone:     phone,
		Email:     each[3],
		OwnerName: each[1] + " " + each[2],
		Address1:  each[8],
		City:      each[9],
		State:     each[10],
		Zip:       each[11],
		Source:    "THUMBTACK",
	}

	fmt.Println("Created provider ", provider)

	return provider
}

func CreateProviderFromFormat4(each []string) *api.Provider {

	util.InfoLog.Println("Create provider - format #4")

	phone, _ := strconv.ParseInt(each[3], 10, 64)

	util.InfoLog.Println("Create provider - name ", each[4])

	provider := &api.Provider{
		Name:      each[4],
		Phone:     phone,
		Email:     each[2],
		OwnerName: each[0] + " " + each[1],
		City:      each[5],
		State:     each[6],
		Source:    "USER",
	}

	fmt.Println("Created provider ", provider)

	return provider
}

func ImportRecord(provider *api.Provider, db api.DB) int64 {
	util.InfoLog.Println("Importing record ", provider)

	id, err := db.AddProvider(provider)

	fmt.Println("Created provider - err  ", err, " id = ", id)

	return id
}

func GetProviderIdsFromDB(db api.DB) []int64 {
	return db.GetAllProviderIds("THUMBTACK")
}
