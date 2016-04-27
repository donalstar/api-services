package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	"io/ioutil"
	"net/http"
	"strconv"
	"trustcloud/util"
)

func AddProject(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {
	project := getProjectFromPost(r)

	projectId, err, projectStatus := db.AddProject(project)

	switch err {
	case nil:
		// do a trustcheck on the provider
		trustcheckStatus := Trustcheck(projectStatus.ProviderId)

		project.Provider.TrustcheckStatus = trustcheckStatus

		projectStatus.TrustcheckStatus = trustcheckStatus.Validated

		// TODO : Location is expected to be an absolute URI, as per the RFC2616
		w.Header().Set("Location", fmt.Sprintf("/project/%d", projectId))

		SendProjectCreatedMail(project)

		return http.StatusCreated, util.Must(enc.Encode(projectStatus))
	default:
		panic(err)
	}
}

// GetProjects returns the list of projects (possibly filtered).
func GetProjects(r *http.Request, enc util.Encoder, db DB) string {
	return util.Must(enc.Encode(db.GetAllProjects()...))
}

// GetProject - the created project & status.
func GetProject(enc util.Encoder, db DB, parms martini.Params) (int, string) {
	id, err := strconv.Atoi(parms["id"])

	project, err := db.GetProject(id)

	if err != nil || project == nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, util.Must(enc.Encode(
			util.NewError(util.ErrCodeNotExist, fmt.Sprintf("the project with id %s does not exist", parms["id"]))))
	}

	return http.StatusOK, util.Must(enc.Encode(project))
}

func DecodeProjectId(enc util.Encoder, db DB, params martini.Params) (int, string) {

	util.InfoLog.Println("DecodeProjectId [", params["id"], "]")

	decodedId := DecodeId(params["id"])

	status := &ProjectId{
		Id: decodedId,
	}

	return http.StatusOK, util.Must(enc.Encode(status))
}

// UpdateProject
func UpdateProject(enc util.Encoder, db DB, params martini.Params) (int, string) {
	id, _ := strconv.Atoi(params["id"])

	status := params["status"]

	webServer := params["webserver"]

	projectStatus, _ := db.UpdateProject(id, status)

	if status != "Pending" {
		project, _ := db.GetProject(id)

		if status == "Approved" {
			SendProjectApprovedMail(project, webServer)
		}

		if status == "Declined" {
			providers := db.GetApprovedProviders(&project.Provider)

			SendProjectDeclinedMail(project, webServer, providers)
		}
	}

	return http.StatusOK, util.Must(enc.Encode(projectStatus))

}

func UpdateJob(enc util.Encoder, db DB, params martini.Params) (int, string) {
	id := int(DecodeId(params["id"]))

	status, _ := strconv.Atoi(params["status"])

	project, _ := db.GetProject(id)

	project.Job.Status = status

	db.UpdateJobStatus(&project.Job)

	SendPurchaseCompleteMail(project)

	return http.StatusOK, util.Must(enc.Encode(project.Job))
}

func AddNote(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {
	noteId, _ := db.AddNote(getNoteFromPost(r))

	return http.StatusOK, util.Must(enc.Encode(noteId))
}

// GetProjects returns the list of projects (possibly filtered).
func GetNotes(r *http.Request, enc util.Encoder, db DB, params martini.Params) string {
	id, _ := strconv.Atoi(params["id"])

	iff := notesToIface(db.GetNotes(id))

	return util.Must(enc.Encode(iff))
}

func DeleteNote(enc util.Encoder, db DB, params martini.Params) (int, string) {
	projectId, _ := strconv.Atoi(params["project_id"])
	noteId, _ := strconv.Atoi(params["id"])

	note := &Note{
		Id:        int64(noteId),
		ProjectId: int64(projectId),
	}

	util.InfoLog.Println("Deleting note ", projectId, " ", noteId)

	db.DeleteNote(note)

	return http.StatusOK, util.Must(enc.Encode(noteId))
}

/*
Do a trustcheck against a provider
*/
func DoTrustcheck(enc util.Encoder, parms martini.Params) string {
	id, _ := strconv.Atoi(parms["id"])

	return util.Must(enc.Encode(Trustcheck(int64(id))))
}

/*
Buy the service guarantee
*/
func BuyGuarantee(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {
	transaction := GetTransactionFromPost(r)

	util.InfoLog.Println("BuyGuarantee - for project id ", transaction.ProjectId)

	url, _ := Checkout(w, r, util.Configuration.General.Guarantee.Name, transaction)

	paypalStatus := &PaypalStatus{
		Url: url,
	}

	return http.StatusCreated, util.Must(enc.Encode(paypalStatus))
}

/*
Complete the PayPal sale
*/
func CompletePayment(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {
	paymentDetails := GetPaymentDetailsFromPost(r)

	util.InfoLog.Println("CompletePayment : ", paymentDetails)

	response, _ := Sale(paymentDetails)

	transactionDetails, _ := GetTransactionDetails(response.Values.Get("PAYMENTINFO_0_TRANSACTIONID"))

	amount, _ := strconv.ParseFloat(response.Values.Get("PAYMENTINFO_0_AMT"), 64)
	feeAmount, _ := strconv.ParseFloat(response.Values.Get("PAYMENTINFO_0_FEEAMT"), 64)

	payment := &Payment{
		ProjectId:         paymentDetails.ProjectId,
		TransactionId:     response.Values.Get("PAYMENTINFO_0_TRANSACTIONID"),
		Email:             transactionDetails.Values.Get("EMAIL"),
		Fee:               feeAmount,
		CurrencyCode:      response.Values.Get("PAYMENTINFO_0_CURRENCYCODE"),
		MerchantAccountId: response.Values.Get("PAYMENTINFO_0_SECUREMERCHANTACCOUNTID"),
		TransactionType:   response.Values.Get("PAYMENTINFO_0_TRANSACTIONTYPE"),
		OrderTime:         response.Values.Get("PAYMENTINFO_0_ORDERTIME"),
		Amount:            amount,
	}

	// log to DB
	db.CreatePayment(payment)

	return http.StatusCreated, util.Must(enc.Encode(response))
}

/*
Do a trustcheck against a provider
*/
func Trustcheck(id int64) TrustcheckStatus {
	util.InfoLog.Println("Do Trustcheck for ", id)

	isValid := (id == 98) // 98 = ACME  (quick test)
	//
	//	fmt.Println("VALID?", isValid)

	trustcheckStatus := &TrustcheckStatus{
		ProviderId: id,
		Score:      725,
		Validated:  isValid,
	}

	return *trustcheckStatus
}

/*
Admin Login
*/
func AdminLogin(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {

	adminUser := getAdminUserFromPost(r)

	result, _ := CheckLogin(*adminUser)

	status := &AdminLoginStatus{
		Message: "OK",
	}

	code := http.StatusOK

	if result {
		code = http.StatusUnauthorized
		status.Message = "Login failed"
	}

	status.Status = code

	return http.StatusOK, util.Must(enc.Encode(status))
}

// GetProvider returns the requested provider.
func GetProvider(enc util.Encoder, db DB, parms martini.Params) (int, string) {

	id, err := strconv.Atoi(parms["id"])
	provider := db.GetProvider(id)

	if err != nil || provider == nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, util.Must(enc.Encode(
			util.NewError(util.ErrCodeNotExist, fmt.Sprintf("the provider with id %s does not exist", parms["id"]))))
	}
	return http.StatusOK, util.Must(enc.Encode(provider))
}

// GetProviders returns the list of providers (possibly filtered).
func GetProviders(r *http.Request, enc util.Encoder, db DB) string {
	// Get the query string arguments, if any
	qs := r.URL.Query()
	name, email, phone := qs.Get("name"), qs.Get("email"), qs.Get("phone")

	if name != "" || email != "" || phone != "" {
		// At least one filter, use Find()
		return util.Must(enc.Encode(db.FindProvider(name, email, phone)...))
	}
	// Otherwise, return all providers
	return util.Must(enc.Encode(db.GetAllProviders()...))
}

/**
get project data via post
*/
func getProjectFromPost(r *http.Request) *Project {
	var project Project

	err := GetFromPost(r, &project)

	if err == nil {
		LogProjectPostParameters(project)
	}

	return &project
}

/**
get transaction data via post
*/
func GetTransactionFromPost(r *http.Request) *Transaction {
	var transaction Transaction

	err := GetFromPost(r, &transaction)

	if err == nil {
		LogTransactionPostParameters(transaction)
	}

	return &transaction
}

/**
get transaction data via post
*/
func GetPaymentDetailsFromPost(r *http.Request) *PaymentDetails {
	var paymentDetails PaymentDetails

	GetFromPost(r, &paymentDetails)

	return &paymentDetails
}

/**
get note data via post
*/
func getNoteFromPost(r *http.Request) *Note {
	var note Note

	err := GetFromPost(r, &note)

	if err == nil {
		LogNotePostParameters(note)
	}

	return &note
}

/**
get transaction data via post
*/
func getAdminUserFromPost(r *http.Request) *AdminUser {
	var adminUser AdminUser

	err := GetFromPost(r, &adminUser)

	if err == nil {
		LogAdminUserPostParameters(adminUser)
	}

	return &adminUser
}

func GetFromPost(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)

	util.InfoLog.Println("POST " + string(body))

	if err == nil {
		err = json.Unmarshal(body, v)

		if err != nil {
			util.ErrorLog.Println("Unable to unmarshall the JSON request", err)
		}
	}

	return err
}
