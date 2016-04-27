package api

import (
	"github.com/crowdmob/paypal"
	"net/http"
	"net/url"
	"trustcloud/util"
)

var isSandbox bool

func init() {
	isSandbox = (util.Configuration.Environment.Payment.TestMode == "Y")
}

/*
Checkout with PayPal
*/
func Checkout(w http.ResponseWriter, r *http.Request, name string, transaction *Transaction) (string, error) {

	// Create the paypal Client with default http client
	client := paypal.NewDefaultClient(util.Configuration.Environment.Payment.User,
		util.Configuration.Environment.Payment.Password,
		util.Configuration.Environment.Payment.Signature,
		isSandbox)

	// Make a array of your digital-goods
	goods := []paypal.PayPalDigitalGood{paypal.PayPalDigitalGood{
		Name:     name,
		Amount:   transaction.PlanCost,
		Quantity: 1,
	}}

	// Sum amounts and get the token!
	response, err := client.SetExpressCheckoutDigitalGoods(paypal.SumPayPalDigitalGoodAmounts(&goods),
		util.Configuration.General.Guarantee.Currency,
		transaction.ReturnUrl,
		transaction.CancelUrl,
		goods,
	)

	return response.CheckoutUrl(), err
}

/*
Do Paypal sale
*/
func Sale(paymentDetails *PaymentDetails) (*paypal.PayPalResponse, error) {
	client := paypal.NewDefaultClient(util.Configuration.Environment.Payment.User,
		util.Configuration.Environment.Payment.Password,
		util.Configuration.Environment.Payment.Signature,
		isSandbox)

	response, err := client.DoExpressCheckoutSale(
		paymentDetails.Token,
		paymentDetails.PayerId,
		util.Configuration.General.Guarantee.Currency,
		paymentDetails.PlanCost)

	return response, err
}

func GetTransactionDetails(transactionId string) (*paypal.PayPalResponse, error) {
	client := paypal.NewDefaultClient(util.Configuration.Environment.Payment.User,
		util.Configuration.Environment.Payment.Password,
		util.Configuration.Environment.Payment.Signature,
		isSandbox)

	values := url.Values{}
	values.Set("METHOD", "GetTransactionDetails")
	values.Add("TRANSACTIONID", transactionId)

	return client.PerformRequest(values)
}
