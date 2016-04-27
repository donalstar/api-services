package check

import (
	"trustcloud/api"
	"trustcloud/util"
)

type Checker interface {
	Check(provider *api.Provider, checkedProvider *api.CheckedProvider)
	GetName() string
}

var Checkers []Checker = []Checker{
	//			&idChecker{},
	//			&businessIdChecker{},
	//	&backgroundChecker{},
	&yelpChecker{},
	&facebookChecker{},
	&fullContactChecker{},
	//	&fullContactCompanyChecker{},
	//	&providerChecker{},
	//	&googleChecker{},
}

/*
Check businesses using LN Business Instantid
*/
func CheckBusinessIds(ids []int64) {
	for _, providerId := range ids {
		CheckBusinessId(providerId)
	}
}

//func BackgroundCheck() {
//	util.InfoLog.Println("Do background check")
//
//	checker := &backgroundChecker{}
//
//	providerId := 3465
//
//	db := api.DataBase
//
//	provider := db.GetProvider(int(providerId))
//
//	util.InfoLog.Println("PROV ", provider)
//
//	checkedProvider := db.GetCheckedProvider(int64(providerId))
//
//	checker.Check(provider, checkedProvider)
//
//	db.AddCheckedProvider(checkedProvider)
//
//}

func CheckBusinesses(ids []int64) {
	for _, providerId := range ids {
		CheckBusiness(providerId)
	}
}

func ExperianBusinessCheckIds(ids []int64) {
	for _, providerId := range ids {
		ExperianBusinessCheckId(providerId)
	}
}

func CheckIndividuals(ids []int64) {
	for _, providerId := range ids {
		CheckIndividual(providerId)
	}
}

func AssessProviders(ids []int64) {
	for _, providerId := range ids {

		AssessProvider(providerId)
	}
}

func CheckWebsites(ids []int64) {
	for _, providerId := range ids {
		CheckWebsite(providerId)
	}
}

/*
Override - log the check to the DB (for use with the Platform Service API)
*/
func AssessProviderAndLog(providerId int64) *api.CheckedProvider {
	return doAssessProvider(providerId, true)
}

func AssessProvider(providerId int64) *api.CheckedProvider {
	return doAssessProvider(providerId, false)
}

func doAssessProvider(providerId int64, doLog bool) *api.CheckedProvider {
	db := api.DataBase

	provider := db.GetProvider(int(providerId))

	util.InfoLog.Println("Assess provider: ", provider.Id)

	checkedProvider := db.GetCheckedProvider(providerId)

	if checkedProvider == nil {
		checkedProvider = &api.CheckedProvider{
			ProviderId: provider.Id,
		}

		util.InfoLog.Println("Adding NEW checked provider for ", providerId)
		db.AddCheckedProvider(checkedProvider)
	}

	for _, checker := range Checkers {
		util.InfoLog.Println("Running ", checker.GetName(), "checker...")

		checker.Check(provider, checkedProvider)

		if doLog {
			db.LogCheck(checker.GetName(), providerId)
		}
	}

	util.InfoLog.Println("ENV    ", util.Configuration.Environment.Database.Name)

	if util.Configuration.Environment.Database.Name == "trustcloud" {
		//			do id check
		util.InfoLog.Println("Do ID check - is business? ", provider.IsBusiness)

		if provider.IsBusiness == "Y" {
			doCheckBusinessId(providerId, checkedProvider)

			if doLog {
				db.LogCheck("LN_IND_ID", providerId)
			}
		}

		if provider.IsBusiness == "N" || provider.IsBusiness == "" {
			doCheckIndividual(providerId, checkedProvider)
			if doLog {
				db.LogCheck("LN_BIZ_INST_ID", providerId)
			}
		}
	}

	if util.Configuration.Environment.Database.Name == "trustcloud_dev" {
		util.InfoLog.Println("DEV - skipping ID check...")
	}

	db.UpdateCheckedProvider(checkedProvider)

	util.InfoLog.Println("Assess provider: ", provider.Id, " ... complete")

	return checkedProvider
}

func CheckBusiness(providerId int64) {
	db := api.DataBase
	checker := &businessIdChecker{}

	provider := db.GetProvider(int(providerId))

	if provider.IsBusiness == "Y" {
		util.InfoLog.Println("Running ", checker.GetName(), "checker...")
		util.InfoLog.Println(providerId, " is a business")

		checkedProvider := db.GetCheckedProvider(providerId)

		//if checkedProvider.LnBizRawJson.Valid == false {
		util.InfoLog.Println("NEED A BUSINESS CHECK - provider id ", providerId)

		checker.Check(provider, checkedProvider)

		db.UpdateCheckedProvider(checkedProvider)
	}
}

func CheckBusinessId(providerId int64) {
	db := api.DataBase
	checker := &bizInstantIdChecker{}

	util.InfoLog.Println("Running ", checker.GetName(), "checker...")

	provider := db.GetProvider(int(providerId))

	if provider.IsBusiness == "Y" {
		util.InfoLog.Println(providerId, " is a business")

		checkedProvider := db.GetCheckedProvider(providerId)

		checker.Check(provider, checkedProvider)

		db.UpdateCheckedProvider(checkedProvider)
	}
}

func ExperianBusinessCheckId(providerId int64) {
	db := api.DataBase
	checker := &experianChecker2{}

	util.InfoLog.Println("got checker ", checker)

	provider := db.GetProvider(int(providerId))

	if provider.IsBusiness == "Y" {
		checkedProvider := db.GetCheckedProvider(providerId)

		checker.Check(provider, checkedProvider)

		db.UpdateCheckedProvider(checkedProvider)
	}
}

func CheckIndividual(providerId int64) {

	db := api.DataBase
	checker := &idChecker{}

	util.InfoLog.Println("Running ", checker.GetName(), "checker...")

	provider := db.GetProvider(int(providerId))

	if provider.IsBusiness == "N" {
		util.InfoLog.Println("NEED AN LN INDIVIDUAL CHECK - provider id ", providerId)

		checkedProvider := db.GetCheckedProvider(providerId)

		checker.Check(provider, checkedProvider)

		db.UpdateCheckedProvider(checkedProvider)
	}
}

func doCheckBusinessId(providerId int64, checkedProvider *api.CheckedProvider) {
	db := api.DataBase
	checker := &bizInstantIdChecker{}

	util.InfoLog.Println("Running ", checker.GetName(), "checker...")

	provider := db.GetProvider(int(providerId))

	if provider.IsBusiness == "Y" {
		util.InfoLog.Println(providerId, " is a business")

		checker.Check(provider, checkedProvider)
	}
}

func doCheckIndividual(providerId int64, checkedProvider *api.CheckedProvider) {

	db := api.DataBase
	checker := &idChecker{}

	util.InfoLog.Println("Running ", checker.GetName(), "checker...")

	provider := db.GetProvider(int(providerId))

	if provider.IsBusiness == "N" {
		util.InfoLog.Println("NEED AN LN INDIVIDUAL CHECK - provider id ", providerId)

		checker.Check(provider, checkedProvider)
	}
}

// Use the FullContact company API to check the provider website
func CheckWebsite(providerId int64) {
	db := api.DataBase
	checker := &fullContactCompanyChecker{}

	provider := db.GetProvider(int(providerId))

	checkedProvider := db.GetCheckedProvider(providerId)

	checker.Check(provider, checkedProvider)

	db.UpdateCheckedProvider(checkedProvider)
}
