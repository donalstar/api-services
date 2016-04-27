package partner

import "trustcloud/api"

type ProviderStatus struct {
	StatusCode   int
	Message      string
	TrustcloudId int64
	Score        int
}

type PartnerInputData struct {
	Provider    api.Provider `json:"provider"`
	PartnerData PartnerData  `json:"partner_data"`
	Key         string       `json:"key"`
	Username    string       `json:"username"`
}

type PartnerData struct {
	Rating float64 `json:"rating"`
	Jobs   int     `json:"jobs"`
}
