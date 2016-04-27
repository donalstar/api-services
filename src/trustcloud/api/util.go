package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"net/url"
	"strconv"
	"trustcloud/util"
)

var (
	cipherKey = []byte("xxxxxxxxxxxxxxxx")
)

// Starting pricing can be 5% of est. job cost with a min of $5 and no more than
// $1K of coverage, so $5 - $50.
// We should round up, perhaps to the nearest $5?
func GetPlanCost(jobCost int) float64 {
	planCost := (jobCost * util.Configuration.General.Guarantee.PricePercentage) / 100

	planCostRem := planCost % util.Configuration.General.Guarantee.PriceRoundUp

	if planCostRem > 0 {
		planCost = planCost + (util.Configuration.General.Guarantee.PriceRoundUp - planCostRem)
	}

	return float64(planCost)
}

func toIface(v []Provider) []interface{} {
	if len(v) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(v))
	for i, v := range v {
		ifs[i] = v
	}
	return ifs
}

func projectsToIface(projectSummaries []ProjectSummary) []interface{} {

	if len(projectSummaries) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(projectSummaries))

	for i, v := range projectSummaries {
		ifs[i] = v
	}
	return ifs
}

func notesToIface(v []*Note) []interface{} {

	if len(v) == 0 {
		return nil
	}
	ifs := make([]interface{}, len(v))
	for i, v := range v {
		ifs[i] = v
	}
	return ifs
}

func LogAdminUserPostParameters(adminUser AdminUser) {
	util.InfoLog.Println("user", adminUser.Username)
}

func LogTransactionPostParameters(transaction Transaction) {
	util.InfoLog.Println("transaction.project_id", transaction.ProjectId)
	util.InfoLog.Println("transaction.return_url", transaction.ReturnUrl)
	util.InfoLog.Println("transaction.cancel_url", transaction.CancelUrl)
}

func LogNotePostParameters(note Note) {
	util.InfoLog.Println("note.Value", note.Value)
	util.InfoLog.Println("note.Date", note.Date)
	util.InfoLog.Println("note.ProjectId", note.ProjectId)
}

func LogProjectPostParameters(project Project) {
	util.InfoLog.Println("project.provider.company", project.Provider.Name)
	util.InfoLog.Println("project.provider.phone ", project.Provider.Phone)
	util.InfoLog.Println("project.provider.type ", project.Provider.Type)
	util.InfoLog.Println("project.provider.email ", project.Provider.Email)
	util.InfoLog.Println("project.provider.ownername ", project.Provider.OwnerName)
	util.InfoLog.Println("project.provider.address1 ", project.Provider.Address1)
	util.InfoLog.Println("project.provider.address2 ", project.Provider.Address2)
	util.InfoLog.Println("project.provider.city ", project.Provider.City)
	util.InfoLog.Println("project.provider.state ", project.Provider.State)
	util.InfoLog.Println("project.provider.zip ", project.Provider.Zip)

	util.InfoLog.Println("project.job.type ", project.Job.Type)
	util.InfoLog.Println("project.job.cost ", project.Job.Cost)
	util.InfoLog.Println("project.job.startdate (raw) ", project.Job.StartDateRaw)

	util.InfoLog.Println("project.user.first_name ", project.User.FirstName)
	util.InfoLog.Println("project.user.last_name ", project.User.LastName)
	util.InfoLog.Println("project.user.email ", project.User.Email)
	util.InfoLog.Println("project.user.phone ", project.User.Phone)
	util.InfoLog.Println("project.user.address1 ", project.User.Address1)
	util.InfoLog.Println("project.user.address2 ", project.User.Address2)
	util.InfoLog.Println("project.user.city ", project.User.City)
	util.InfoLog.Println("project.user.state ", project.User.State)
	util.InfoLog.Println("project.user.zip  ", project.User.Zip)
}

func SetDisplayableFields(project *Project) {

	project.Job.DisplayableCost = "$" + strconv.Itoa(project.Job.Cost) + ".00"

	project.Job.DisplayableStartDate = util.FormatDateForDisplay(project.Job.StartDate)

	project.Provider.DisplayablePhone = util.FormatPhone(project.Provider.Phone)

	util.InfoLog.Println("format phone / after ", project.Provider.DisplayablePhone)
	project.User.DisplayablePhone = util.FormatPhone(project.User.Phone)
}

func EncodeId(id int) string {
	data := []byte(strconv.FormatInt(int64(id), 10))

	ciphertext, _ := encrypt(cipherKey, data)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	// 2 times
	return url.QueryEscape(url.QueryEscape(str))
}

func DecodeId(id string) int64 {

	unescapedId, _ := url.QueryUnescape(id)

	data, _ := base64.StdEncoding.DecodeString(unescapedId)

	result, _ := decrypt(cipherKey, data)

	i, _ := strconv.ParseInt(string(result), 10, 64)

	return i
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func decrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}
	return data, nil
}
