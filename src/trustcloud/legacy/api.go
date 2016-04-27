package legacy

import (
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"strconv"
	"trustcloud/api"
	"trustcloud/util"
)

/*
Admin Login
*/
func AdminLogin(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB, params martini.Params) (int, string) {
	user := params["id"]
	password := params["password"]

	adminUser, err := db.GetUser(user, password)

	util.InfoLog.Println("Admin User ", adminUser.UserName, " ", adminUser.Password)

	code := http.StatusOK
	value := "OK"

	if err != nil {
		code = http.StatusUnauthorized
		value = "User login incorrect"
	}

	status := &AdminLoginStatus{
		Status:     code,
		Message:    value,
		Id:         adminUser.Id,
		ApiAccount: adminUser.ApiAccount,
	}

	return code, util.Must(enc.Encode(status))
}

/*
Buy Trustcheck
*/
func BuyTrustcheck(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {
	transaction := api.GetTransactionFromPost(r)

	util.InfoLog.Println("BuyTrustcheck")

	url, _ := api.Checkout(w, r, "myTrustcheck", transaction)

	paypalStatus := &api.PaypalStatus{
		Url: url,
	}

	return http.StatusCreated, util.Must(enc.Encode(paypalStatus))
}

func CompleteTrustcheckPayment(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {
	paymentDetails := api.GetPaymentDetailsFromPost(r)

	util.InfoLog.Println("CompleteTrustcheckPayment : ", paymentDetails)

	response, _ := api.Sale(paymentDetails)

	util.InfoLog.Println("response : ", response)

	transactionId := response.Values.Get("PAYMENTINFO_0_TRANSACTIONID")

	transactionDetails, _ := api.GetTransactionDetails(transactionId)

	paymentInfoJson, _ := enc.Encode(response.Values)

	amount, _ := strconv.ParseFloat(response.Values.Get("PAYMENTINFO_0_AMT"), 64)
	feeAmount, _ := strconv.ParseFloat(response.Values.Get("PAYMENTINFO_0_FEEAMT"), 64)

	transaction := &LegacyTransaction{
		TransactionObject: paymentInfoJson,
		PartnerId:         1,
		Email:             transactionDetails.Values.Get("EMAIL"),
		Payment:           amount,
		ProductId:         1,
		ProductCode:       "1",
		TransactionId:     transactionId,
		PaymentType:       "paypal",
		Fee:               feeAmount,
		ReturnUrl:         paymentDetails.ReturnUrl,
	}

	db.AddLegacyTransaction(transaction)

	util.InfoLog.Println("returning transaction  : ", transaction)

	return http.StatusCreated, util.Must(enc.Encode(transaction))
}

func GetLegacyTransaction(enc util.Encoder, db DB, params martini.Params) (int, string) {
	id := params["id"]

	transaction, err := db.GetLegacyTransaction(id)

	if err != nil || transaction == nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, util.Must(enc.Encode(
			util.NewError(util.ErrCodeNotExist, fmt.Sprintf("the transaction with id %s does not exist", params["id"]))))
	}

	return http.StatusOK, util.Must(enc.Encode(transaction))
}

func RenderCard(r *http.Request, enc util.Encoder, db DB, params martini.Params) (int, string) {
	id := params["id"]

	card_size, _ := strconv.Atoi(params["size"])

	/*
				    // handle different sizes
		    if (strstr($_GET['userid'], "-b1")) {
		        $userid = str_replace("-b1", "", $_GET["userid"]);
		        $exec = 'xvfb-run -a --server-args="-screen 0, 640x480x24" ' . FS_ROOT_INCLUDE . 'lib/cutycapt/CutyCapt --url='
		            . PLATFORM_URL . 'display/renderidcard?svg=no\&allowUpdate=no\&userid='
		            . urlencode($userid) . '\&size=b1 --out=/tmp/' . $filename . '.png --zoom-factor=2.0 --min-width=100 --min-height=100';
		    } else if (strstr($_GET['userid'], "-b2")) {
		        $userid = str_replace("-b2", "", $_GET["userid"]);
		        $exec = 'xvfb-run -a --server-args="-screen 0, 640x480x24" ' . FS_ROOT_INCLUDE . 'lib/cutycapt/CutyCapt --url='
		            . PLATFORM_URL . 'display/renderidcard?svg=no\&allowUpdate=no\&userid='
		            . urlencode($userid) . '\&size=b2 --out=/tmp/' . $filename . '.png --zoom-factor=2.0 --min-width=100 --min-height=100';
		    } else if (strstr($_GET['userid'], "-b3")) {
		        $userid = str_replace("-b3", "", $_GET["userid"]);
		        $exec = 'xvfb-run -a --server-args="-screen 0, 640x480x24" ' . FS_ROOT_INCLUDE . 'lib/cutycapt/CutyCapt --url='
		            . PLATFORM_URL . 'display/renderidcard?svg=no\&allowUpdate=no\&userid='
		            . urlencode($userid) . '\&size=b3 --out=/tmp/' . $filename . '.png --zoom-factor=2.0 --min-width=100 --min-height=100';
		    } else // regular card
		    {
		        $exec = 'xvfb-run -a --server-args="-screen 0, 640x480x24" ' . FS_ROOT_INCLUDE . 'lib/cutycapt/CutyCapt --url='
		            . PLATFORM_URL . 'display/renderidcard?svg=no\&allowUpdate=no\&userid='
		            . urlencode($_GET['userid']) . ' --out=/tmp/' . $filename . '.png --zoom-factor=2.0 --min-width=100 --min-height=100';
		    }
	*/

	util.InfoLog.Println("render card for ", id, " size: ", card_size)

	cardDetails := &CardDetails{
		Id:   id,
		Size: card_size,
	}

	renderCard(cardDetails)

	return http.StatusOK, util.Must(enc.Encode(cardDetails))
}

func GetTransaction(enc util.Encoder, db DB, params martini.Params) (int, string) {
	id := params["id"]

	supportTicket, err := db.GetSupportTicket(id)

	if err != nil || supportTicket == nil {
		// Invalid id, or does not exist
		return http.StatusNotFound, util.Must(enc.Encode(
			util.NewError(util.ErrCodeNotExist, fmt.Sprintf("the transaction with id %s does not exist", params["id"]))))
	}

	return http.StatusOK, util.Must(enc.Encode(supportTicket))
}

func GetCannedResponses(r *http.Request, enc util.Encoder, db DB) string {
	return util.Must(enc.EncodeAsArray(db.GetCannedResponses()...))
}

func AddCannedResponse(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {
	id, _ := db.AddCannedResponse(getCannedResponseFromPost(r))

	return http.StatusOK, util.Must(enc.Encode(id))
}

func DeleteCannedResponse(enc util.Encoder, db DB, params martini.Params) (int, string) {
	id, _ := strconv.Atoi(params["id"])

	response := &CannedResponse{
		Id: int64(id),
	}

	db.DeleteCannedResponse(response)

	return http.StatusOK, util.Must(enc.Encode(id))
}

func UpdateCannedResponse(w http.ResponseWriter, r *http.Request, enc util.Encoder, db DB) (int, string) {
	util.InfoLog.Println("UpdateCannedResponse")

	response := &CannedResponse{
		Id: 7,
	}

	return http.StatusOK, util.Must(enc.Encode(response))
}

/**
get note data via post
*/
func getCannedResponseFromPost(r *http.Request) *CannedResponse {
	var cannedResponse CannedResponse

	api.GetFromPost(r, &cannedResponse)

	return &cannedResponse
}
