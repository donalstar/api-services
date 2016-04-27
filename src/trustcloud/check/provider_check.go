package check

import (
	"trustcloud/api"
	"trustcloud/util"
)

type providerChecker struct {
}

func (db *providerChecker) GetName() string {
	return "TC_PROVIDER"
}

func (db *providerChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	id := api.DataBase.GetProviderFromMaster(provider)

	if id != 0 {
		util.InfoLog.Println("Got a match")
		checkedProvider.InMasterList = "Y"
	}

	if id == 0 {
		checkedProvider.InMasterList = "N"
	}
}
