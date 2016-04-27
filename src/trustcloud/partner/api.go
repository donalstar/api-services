package partner

import (
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"trustcloud/api"
	"trustcloud/check"
	"trustcloud/util"
)

func AddProvider(w http.ResponseWriter, r *http.Request, enc util.Encoder, db api.DB) (int, string) {
	data := getPartnerProviderFromPost(r)

	authenticated := db.CheckAuthorizedUser(data.Username, data.Key)

	if authenticated == false {
		util.ErrorLog.Println("Authorization failure ", data.Username, " ", data.Key)

		status := &ProviderStatus{
			StatusCode: http.StatusUnauthorized,
			Message:    "Request not authorized - check your key and username",
		}

		return http.StatusUnauthorized, util.Must(enc.Encode(status))
	}

	util.InfoLog.Println("Authorization ok")

	provider := data.Provider

	clientPlatform, _ := db.GetClientPlatform(data.Key)

	util.InfoLog.Println("Got client platform ", clientPlatform)

	provider.ClientPlatformId = clientPlatform.Id

	// check & see if provider exists - create, if not, update otherwise
	matchedProviders := db.FindProvider(provider.Name, provider.Email, strconv.FormatInt(provider.Phone, 10))

	var providerId int64

	if matchedProviders == nil {
		util.InfoLog.Println("Partner API - no provider match -- create a new one")

		id, err := db.AddProvider(&provider)

		if err == nil {
			providerId = id
			util.InfoLog.Println("Created new provider - id: ", providerId)
		}
	}

	if matchedProviders != nil {
		util.InfoLog.Println("Provider already exists")

		matchedProvider := matchedProviders[0].(api.Provider)

		providerId = matchedProvider.Id
	}

	status := &ProviderStatus{
		StatusCode: 100,
		Message:    "Error",
	}

	status = CheckProvider(providerId)

	return http.StatusCreated, util.Must(enc.Encode(status))
}

// UpdateProvider
func UpdateProvider(enc util.Encoder, db api.DB, params martini.Params) int {
	util.InfoLog.Println("Partner API - Update Provider")

	return http.StatusCreated
}

//GetProviderCard
func GetProviderCard(enc util.Encoder, db api.DB, params martini.Params) string {
	id, _ := strconv.Atoi(params["id"])

	provider := db.GetProvider(id)

	checkedProvider := db.GetCheckedProvider(int64(id))

	score := 0

	if checkedProvider != nil {
		score = check.CalculateTrustscore(checkedProvider)
	}

	providerData := &util.Provider2{
		Name:  provider.Name,
		City:  provider.City,
		State: provider.State,
		Score: score,
	}

	return util.GetCard(*providerData)
}

/**
get project data via post
*/
func getPartnerProviderFromPost(r *http.Request) *PartnerInputData {
	var data PartnerInputData

	api.GetFromPost(r, &data)

	util.InfoLog.Println("Use partner input data ", data)

	util.InfoLog.Println("date of birth ", data.Provider.DateOfBirth)

	return &data
}

func CheckProvider(providerId int64) *ProviderStatus {
	checkedProvider := check.AssessProviderAndLog(providerId)

	util.InfoLog.Println("Checked provider: ", checkedProvider)

	score := check.CalculateTrustscore(checkedProvider)

	// todo: save score data to "score" DB table

	status := &ProviderStatus{
		StatusCode:   200,
		Message:      "OK",
		TrustcloudId: providerId,
		Score:        score,
	}

	return status
}
