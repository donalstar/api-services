// select count(*) from CheckedProvider where ln_biz_raw_json like '%response%'  order by id desc

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"trustcloud/api"
	//"trustcloud/util"
	//	"bufio"
	"database/sql"
	//	"strconv"
	"strconv"
	"trustcloud/check"
	//	"trustcloud/util"
	"trustcloud/util"
)

func main() {
	fmt.Println("Updater...")
	//
	//	results := GetBizSearchResults()
	//
	//	for _, row := range results {
	//		UpdateCheckedProvider(row)
	//	}

	//	inputFile := util.BaseDir + "/test/" + "tt2.csv"

	//	inputFile := util.BaseDir + "/test/" + "missed-tt-providers.csv"
	inputFile := util.BaseDir + "/test/" + "thumbtack-data-mar17-2.csv"

	//	checkIsBusiness(inputFile)

	//	thumbtackParse(inputFile)

	//	thumbtackParse(inputFile)

	//populateThumbtackProviderDataTable(inputFile)

	//	Recon3(inputFile)

	setTTScoredFlag(inputFile)
}

func Recon3(fileName string) {

	var ids []int64
	var matched_checked []int64
	var no_matched_checked []int64
	var too_many []int64

	api.DbMap.Select(&ids, "select id from Provider where source = 'THUMBTACK'")

	for _, id := range ids {
		var checkedProviders []int64

		api.DbMap.Select(&checkedProviders, "select id from CheckedProvider where provider_id = ?", id)

		if len(checkedProviders) > 0 {
			matched_checked = append(matched_checked, id)

			if len(checkedProviders) > 1 {
				too_many = append(too_many, id)

			}
		}

		if len(checkedProviders) == 0 {
			no_matched_checked = append(no_matched_checked, id)
		}
	}

	fmt.Println("no match checked  ", no_matched_checked)
	fmt.Println("too_many  ", too_many)

	//	var allCheckedProviders []int64
	//	api.DbMap.Select(&allCheckedProviders, "select provider_id from CheckedProvider")
	//
	//	var no_match []int64
	//	for _, id := range allCheckedProviders {
	//		match := false
	//		for _, mc := range matched_checked {
	//			if mc == id {
	//				match = true
	//			}
	//		}
	//
	//		if match == false {
	//			no_match = append(no_match, id)
	//		}
	//	}
	//
	//	fmt.Println("no matches ", no_match)
}

func setTTScoredFlag(fileName string) {
	//	var ttData []api.ThumbtackProviderData
	//
	//	api.DbMap.Select(&ttData, "select * from ThumbtackProviderData")
	//
	//	for _, each := range ttData {
	//		fmt.Println("NEXT: ", each)
	//	}

	rawCSVdata := readRecordsFromCsv(fileName)

	/*
		1366,Tony,Russo,uncletonio@gmail.com,5126324022,,N,Tjscapes,216 Brushy Hill Rd,Spicewood,TX,78669,,,Sprinkler System Repair and Maintenance,35,2,3,5,0,No,0,0
		1367,Tony,Rosa,kettlebellcoachtony@gmail.com,4436947905,,N,Iron Storm Kettlebell Strength and Conditioning,650 State Rt. 3,Gambrills,MD,21054,,,Personal Training,12,0,3,5,0,No,0,0
		251,Terry,Roper,terryroper31@yahoo.com,2148305254,,N,Double Time Roofing & Remodeling,2852 North Hampton Dr,Grand Prairie,TX,75052,N,terryroper31@yahoo.com,Kitchen Remodel,,,,,,,,

	*/
	for _, each := range rawCSVdata {

		name := each[7]
		email := each[3]

		fmt.Println("name: ", name, " len ", len(each))

		var provider api.Provider

		err := api.DbMap.SelectOne(&provider, "select * from Provider where name=? and email=?",
			name, email)

		fmt.Println("got provider ", provider.Id, " ", err)

		var ttData api.ThumbtackProviderData

		err = api.DbMap.SelectOne(&ttData, "select * from ThumbtackProviderData where provider_id=?",
			provider.Id)

		fmt.Println("update TT record ", ttData, " ", err)

		ttData.Scored = "N"

		reviews := 0
		if len(each) > 17 {
			reviews, _ = strconv.Atoi(each[17])
			fmt.Println("got a review ", reviews)
			ttData.Scored = "Y"
		}

		_, err = api.DbMap.Update(&ttData)

		fmt.Println("updated? ", err)
	}
}

func Recon2(fileName string) {
	var ids []int64

	rawCSVdata := readRecordsFromCsv(fileName)

	for _, each := range rawCSVdata {
		email := each[3]
		name := each[7]

		var providers []api.Provider

		api.DbMap.Select(&providers, "select * from Provider where name = ? and email = ? and source = 'THUMBTACK'",
			name, email)

		fmt.Println("matched ", email, " ", name, " : id ", providers[0].Id)

		ids = append(ids, providers[0].Id)
	}

	providerIds := check.GetProviderIdsFromDB(api.DataBase)

	var matched []int64

	for _, id := range providerIds {
		for _, tid := range ids {
			if tid == id {
				fmt.Println("matched P ID", id, " to TID ", tid)
				matched = append(matched, id)
				break
			}
		}
	}

	var no_match []int64
	for _, id := range providerIds {
		match := false
		for _, tid := range matched {
			if tid == id {
				match = true
				break
			}
		}

		if match == false {
			no_match = append(no_match, id)
		}
	}

	fmt.Println("NO MATCH ", no_match)
}

func Recon(fileName string) {
	fmt.Println("Recon")

	badCount := 0

	rawCSVdata := readRecordsFromCsv(fileName)

	for _, each := range rawCSVdata {
		email := each[3]
		name := each[7]

		var providers []api.Provider

		api.DbMap.Select(&providers, "select * from Provider where name = ? and email = ? and source = 'THUMBTACK'",
			name, email)

		if len(providers) == 0 {
			badCount = badCount + 1
			fmt.Println("NO MATCH: ", name, email, " bad count ", badCount)

			var p2 []api.Provider
			api.DbMap.Select(&p2, "select * from Provider where email = ? and source = 'THUMBTACK'",
				email)

			if len(p2) == 0 {
				fmt.Println("\tNO MATCH with just email")

			}

			if len(p2) > 0 {
				fmt.Println("\tNMATCHED with just email")

				p2[0].Name = name
				api.DbMap.Update(&p2[0])
			}
		}

	}
}

func GetBizSearchResults() []api.BizSearch {
	fmt.Println("GetBizSearchResults")

	var results []api.BizSearch

	api.DbMap.Select(&results, "select * from LnBusinessCheck")

	return results
}

func UpdateCheckedProvider(bizSearchResult api.BizSearch) {
	id := bizSearchResult.ProviderId

	var providers []api.CheckedProvider
	api.DbMap.Select(&providers, "select * from CheckedProvider where provider_id = ?", id)

	for _, row := range providers {
		fmt.Println("Updating checked provider row ", row.Id, " with JSON ", bizSearchResult.LnBizRawJson)

		bizMatch := bizSearchResult.LnBizMatch

		if bizMatch == "" {
			bizMatch = "N"
		}

		bizMatchedPhone := bizSearchResult.LnBizMatchedPhone

		if bizMatchedPhone == "" {
			bizMatchedPhone = "N"
		}

		rawJson := bizSearchResult.LnBizRawJson

		row.LnBizMatch = bizMatch
		row.LnBizMatchedPhone = bizMatchedPhone
		row.LnBizRawJson.String = rawJson
		row.LnBizRawJson.Valid = true

		api.DbMap.Update(&row)
	}

}

func populateThumbtackProviderDataTable(fileName string) {
	rawCSVdata := readRecordsFromCsv(fileName)

	for _, each := range rawCSVdata {
		if len(each) >= 14 {
			email := each[3]

			// Get Provider record & update
			var provider api.Provider

			error := api.DbMap.SelectOne(&provider, "select * from Provider where email=?", email)

			if error == nil {

				var ttData []api.ThumbtackProviderData

				api.DbMap.Select(&ttData, "select * from ThumbtackProviderData where provider_id=?", provider.Id)

				if len(ttData) == 0 {
					fmt.Println("no match in ThumbtackProviderData for ", email)

					if provider.Email == email {
						fmt.Println("# pieces ", len(each))

						bids := 0
						if len(each) > 15 {
							bids, _ = strconv.Atoi(each[15])
						}

						hires := 0
						if len(each) > 16 {
							hires, _ = strconv.Atoi(each[16])
						}

						reviews := 0
						if len(each) > 17 {
							reviews, _ = strconv.Atoi(each[17])
						}

						averageReviewScore := 0.0
						if len(each) > 18 {
							averageReviewScore, _ = strconv.ParseFloat(each[18], 64)
						}

						averageReimbursement := 0.0
						if len(each) > 19 {
							averageReimbursement, _ = strconv.ParseFloat(each[21], 64)
						}

						complaints := 0
						if len(each) > 20 {
							complaints, _ = strconv.Atoi(each[19])
						}

						accDeleted := "N"
						if len(each) > 21 {
							accDeleted = each[20]
						}

						bgVerified := "N"
						if len(each) > 22 {
							bgVerified = each[22]
						}

						thumbtackProviderData := &api.ThumbtackProviderData{
							ProviderId:           provider.Id,
							Bids:                 bids,
							Hires:                hires,
							Reviews:              reviews,
							HasValidWebsite:      each[12],
							AverageReviewScore:   averageReviewScore,
							Complaints:           complaints,
							AverageReimbursement: averageReimbursement,
							AccountDeleted:       accDeleted,
							BackgroundVerified:   bgVerified,
						}

						api.DbMap.Insert(thumbtackProviderData)

						fmt.Println("created thumbtackProviderData ", thumbtackProviderData)
					}
				}
			}
		}
	}

}

func readRecordsFromCsv(fileName string) [][]string {
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

	return rawCSVdata
}

func thumbtackParse(fileName string) {
	csvfile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return
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
		if len(each) >= 14 {
			email := each[3]
			fmt.Println("email ", each[3])

			/*
				12 - valid web
				13 - website
				14 - category
				15 bids
				16 hires
				17 reviews
				18 av rev score
				19 # complaints
				20 acc deleted
				21 av reimburs
				22 bk verified
			*/
			//		fmt.Println("valid web ", each[12])
			website := each[13]

			category := each[14]

			if category != "" || website != "" {

				// Get Provider record & update
				var provider api.Provider

				fmt.Println("got  provider with email ", email)
				error := api.DbMap.SelectOne(&provider, "select * from Provider where email=?", email)

				fmt.Println("got  provider ", provider, " err ", error)

				if error == nil {
					// update provider
					if provider.Email == email {
						fmt.Println("Matched provider email")

						/*
							update: website, category
						*/

						fmt.Println("update with ", category, " ", website)

						provider.Category = sql.NullString{String: category, Valid: true}
						provider.Website = sql.NullString{String: website, Valid: true}

						_, err := api.DbMap.Update(&provider)
						fmt.Println("updated provider ", provider.Id, " err ", err)
					}
				}
			}
			//		fmt.Println("bids ", each[15])
			//		fmt.Println("hires ", each[16])
			//		fmt.Println("reviews ", each[17])
			//		fmt.Println("av review score ", each[18])
			//		fmt.Println("complaints ", each[19])
			//		fmt.Println("acc deleted ", each[20])
			//		fmt.Println("av reimburs ", each[21])
			//		fmt.Println("bk verified ", each[22])
		}

	}

}

func checkIsBusiness(fileName string) {
	csvfile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = 6 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, each := range rawCSVdata {
		fmt.Println("Next line [", each, "]")

		email := each[2]

		fmt.Println("email ", email, " is indiv. ", each[4])

		// Get Provider record & update
		var provider api.Provider

		error := api.DbMap.SelectOne(&provider, "select * from Provider where email=?", email)

		if error == nil {
			// update provider
			if provider.Email == email {
				fmt.Println("Matched provider email")

				isBusiness := "N"

				if each[4] == "N" {
					isBusiness = "Y"
				}

				provider.IsBusiness = isBusiness

				fmt.Println("is biz? ", isBusiness)
				api.DbMap.Update(&provider)
			}
		}
	}
}
