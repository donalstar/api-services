package main

import (
	"trustcloud/partner"
	"trustcloud/util"
)

func main() {
	util.InfoLog.Println("partner_api_check_test")

	ids := []int64{
		3304,
		3287,
		3252,
		3244,
		3316,
		3314,
		3296,
		3288,
		3316,
		3314,
		3427,
		3320,
		3317,
		3315,
		2429,
		3432,
		1337,
		1338,
		3470,
		3471,
		3472,
		3473,
		3474,
		3475,
		3476,
		3477,
		3478,
		3479,
		3480,
		3481,
		3482,
		3483,
		3484,
		3485,
		3486,
		3487,
		3488,
		3489,
		3490,
		3491,
		3492,
		3493,
		3494,
		3495,
	}

	for _, id := range ids {
		util.InfoLog.Println("assess provider ", id)

		status := partner.CheckProvider(id)

		util.InfoLog.Println("result: ", status)
	}
	//	func checkProvider(providerId int64) *ProviderStatus {
	//
	//	}
}
