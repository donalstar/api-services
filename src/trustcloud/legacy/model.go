package legacy

import (
	"database/sql"
	"time"
)

type CardDetails struct {
	Id   string
	Size int
}

type User struct {
	Id               int64  `db:"id"`
	FirstName        string `db:"first_name" json:"firstname"`
	LastName         string `db:"last_name"  json:"lastname"`
	Email            string `db:"email"`
	Phone            int64  `db:"phone"`
	Address1         string `db:"address1" json:"address1"`
	Address2         string `db:"address2" json:"address2"`
	City             string `db:"city" json:"city"`
	State            string `db:"state" json:"state"`
	Zip              string `db:"zip" json:"zip"`
	DisplayablePhone string `db:"-"`
}

type SupportTicketView struct {
	TicketId          int64          `db:"intTicketID"`
	TransactionId     string         `db:"strTransactionID"`
	FirstName         string         `db:"strFirstName"`
	MiddleName        sql.NullString `db:"strMiddleName"`
	LastName          string         `db:"strLastName"`
	Address1          string         `db:"strAddress1"`
	Address2          sql.NullString `db:"strAddress2"`
	City              string         `db:"strCity"`
	State             string         `db:"strState"`
	Zip               string         `db:"strZip"`
	Var1              sql.NullString `db:"strVar1"`
	Var2              sql.NullString `db:"strVar2"`
	Var3              sql.NullString `db:"strVar3"`
	Phone             string         `db:"strPhone"`
	Email             string         `db:"strEmail"`
	RequestDatetime   time.Time      `db:"dtmRequestDatetime"`
	AssignedTo        int64          `db:"intAssignedTo"`
	Description       sql.NullString `db:"strDescription"`
	Comments          sql.NullString `db:"strComments"`
	CaseComments      sql.NullString `db:"strCaseComments"`
	StatusId          int64          `db:"intStatusID"`
	Type              string         `db:"strType"`
	Priority          int64          `db:"intPriority"`
	EmailStatus       int64          `db:"intEmailStatus"`
	Timestamp         time.Time      `db:"dtmTimeStamp"`
	CustomerResponses sql.NullString `db:"strCustomerResponses"`
	Username          sql.NullString `db:"strUsername"`
	ApiAccount        sql.NullString `db:"intApiAccount"`
}

type AdminLoginStatus struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	Id         int64  `json:"id"`
	ApiAccount int64  `json:"api_account"`
}

type AdminUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LegacyTransaction struct {
	Id                int64          `db:"idTransaction"`
	Email             string         `db:"strEmail"`
	TransactionObject string         `db:"jsontransactionobject"`
	PartnerId         int            `db:"intpartnerid"`
	Payment           float64        `db:"decpayment"`
	ProductId         int            `db:"intproductid"`
	ProductCode       string         `db:"strproductcode"`
	TransactionId     string         `db:"strtransactionid"`
	Fee               float64        `db:"decfee"`
	IpAddress         string         `db:"stripaddress"`
	PaymentType       string         `db:"enumpaymenttype"`
	ReturnUrl         string         `db:"strreturnurl"`
	Delivered         bool           `db:"booldelivered"`
	DeliveredTime     sql.NullString `db:"dtmdeliveredtime'"`
	Query             sql.NullString `db:"blobQuery"`
	Deliverable       sql.NullString `db:"blobDeliverable"`
	ReportFilename    sql.NullString `db:"strreportfilename"`
	CustomVar1        sql.NullString `db:"customvar1"`
	CustomVar2        sql.NullString `db:"customvar2"`
	PostbackUrl       sql.NullString `db:"postBackURL"`
	Transactioncol    sql.NullString `db:"Transactioncol"`
	ApiUsername       sql.NullString `db:"strApiUsername"`
}

type ApiAdminUser struct {
	Id                int64  `db:"idApiUser"`
	ApiAccount        int64  `db:"intApiAccount"`
	UserName          string `db:"strUsername"`
	Password          string `db:"strPassword"`
	AllowDetailView   int    `db:"intAllowDetailView"`
	AllowResponseView int    `db:"intAllowResponseView"`
}

type CannedResponse struct {
	Id          int64  `db:"ID" json:"id"`
	Response    string `db:"Response" json:"response"`
	Description string `db:"Description" json:"description"`
}
