package api

import (
	"database/sql"
	"time"
)

type Project struct {
	Id           int
	Job          Job       `json:"job"`
	Provider     Provider  `json:"provider"`
	User         User      `json:"user"`
	Status       string    `json:"status"`
	CreateDate   time.Time `json:"create_date"`
	Server       string
	SupportPhone string `json:"support_phone"`
	EncodedId    string `json:"encoded_id"`
}

type ProjectStatus struct {
	Id               int       `db:"id"`
	EncodedId        string    `db:"encoded_id" json:"encoded_id"`
	ProviderId       int64     `db:"provider_id"`
	UserId           int64     `db:"user_id"`
	JobId            int64     `db:"job_id"`
	Status           string    `db:"status"`
	TrustcheckStatus bool      `db:"-"`
	CreateDate       time.Time `db:"create_date"`
}

type ProjectId struct {
	Id int64 `json:"id"`
}

type ProjectSummary struct {
	Id            int       `db:"project_id" json:"id"`
	EncodedId     string    `db:"encoded_id" json:"encoded_id"`
	JobType       string    `db:"type" json:"type"`
	JobCost       string    `db:"cost" json:"cost"`
	StartDate     time.Time `db:"startdate"`
	StartDateRaw  string    `db:"-" json:"start_date"`
	UserFirstName string    `db:"first_name" json:"first_name"`
	UserLastName  string    `db:"last_name" json:"last_name"`
	ProviderName  string    `db:"provider_name" json:"provider_name"`
	Status        string    `db:"status" json:"status"`
	CreateDate    time.Time `db:"create_date" json: create_date`
}

type Job struct {
	Id                   int64     `db:"id"`
	Type                 string    `db:"type" json:"type"`
	Cost                 int       `db:"cost" json:"cost"`
	Description          string    `db:"description" json:"description"`
	StartDateRaw         string    `db:"-"`
	StartDate            time.Time `db:"startdate" json:"startdate"`
	PlanCost             float64   `db:"-" json:"plan_cost"`
	PlanName             string    `db:"-" json:"plan_name"`
	DisplayableCost      string    `db:"-"`
	DisplayableStartDate string    `db:"-"`
	Status               int       `db:"status" json:"status"`
}

type Provider struct {
	Id               int64            `db:"id" json: "id"`
	Name             string           `db:"name" json:"name"`
	Phone            int64            `db:"phone" json:"phone,string"`
	Type             string           `db:"type" json:"type"`
	Email            string           `db:"email" json:"email"`
	OwnerName        string           `db:"ownername" json:"ownername"`
	Address1         string           `db:"address1" json:"address1"`
	Address2         string           `db:"address2" json:"address2"`
	City             string           `db:"city" json:"city"`
	State            string           `db:"state" json:"state"`
	Zip              string           `db:"zip" json:"zip"`
	Category         sql.NullString   `db:"category" json:"category"`
	Website          sql.NullString   `db:"website" json:"website"`
	IsBusiness       string           `db:"is_business" json:"is_business"`
	DateOfBirth      time.Time        `db:"date_of_birth" json:"date_of_birth"`
	Source           string           `db:"source" json: "source"`
	ClientPlatformId int64            `db:"client_platform_id" json: "client_platform_id"`
	TrustcheckStatus TrustcheckStatus `db:"-"`
	DisplayablePhone string           `db:"-"`
}

type User struct {
	Id               int64  `db:"id"`
	FirstName        string `db:"first_name" json:"firstname"`
	LastName         string `db:"last_name"  json:"lastname"`
	Email            string `db:"email"`
	Phone            int64  `db:"phone" json:"phone,string"`
	Address1         string `db:"address1" json:"address1"`
	Address2         string `db:"address2" json:"address2"`
	City             string `db:"city" json:"city"`
	State            string `db:"state" json:"state"`
	Zip              string `db:"zip" json:"zip"`
	DisplayablePhone string `db:"-"`
}

type Checker struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type CheckerLog struct {
	Id         int64     `db:"id"`
	ProviderId int64     `db:"provider_id"`
	CheckId    int64     `db:"check_id"`
	CheckTime  time.Time `db:"check_time"`
}

type ClientPlatform struct {
	Id               int64  `db:"id"`
	Name             string `db:"name"`
	Type             string `db:"type"`
	DoIdVerification int    `db:"do_id_verification"`
	DoBkgCheck       int    `db:"do_bkg_check"`
}

type ClientCredentials struct {
	Id               int64  `db:"id"`
	ClientPlatformId int64  `db:"client_platform_id"`
	ApiKey           string `db:"api_key"`
}

type Payment struct {
	Id                int64   `db:"id"`
	ProjectId         int64   `db:"project_id"`
	TransactionId     string  `db:"transaction_id"`
	Email             string  `db:"email"`
	Fee               float64 `db:"fee"`
	CurrencyCode      string  `db:"currency_code"`
	MerchantAccountId string  `db:"merchant_account_id"`
	TransactionType   string  `db:"transaction_type"`
	OrderTime         string  `db:"order_time"`
	Amount            float64 `db:"amount"`
}

type TrustcheckStatus struct {
	ProviderId int64
	Score      int
	Validated  bool
}

type AdminLoginStatus struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Note struct {
	Id        int64     `db:"id" json:"id"`
	Date      time.Time `db:"date" json:"date"`
	Value     string    `db:"value" json:"value"`
	ProjectId int64     `db:"project_id" json:"project_id,string"`
}

type PaypalStatus struct {
	Url string `json:"url"`
}

type Transaction struct {
	ProjectId string  `json:"project_id"`
	PlanCost  float64 `json:"plan_cost"`
	ReturnUrl string  `json:"return_url"`
	CancelUrl string  `json:"cancel_url"`
}

type PaymentDetails struct {
	ProjectId int64   `json:"project_id"`
	PayerId   string  `json:"payer_id"`
	Token     string  `json:"token"`
	PlanCost  float64 `json:"plan_cost"`
	ReturnUrl string  `json:"return_url"`
}

type AdminUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Product struct {
	Name     string
	Price    float64
	Currency string
}

type MailTemplateData struct {
	Project   Project
	Providers []Provider
}

type ProviderMaster struct {
	Id          int64   `db:"id"`
	CompanyName string  `db:"company_name" json:"company_name"`
	City        string  `db:"city" json:"city"`
	State       string  `db:"state" json:"state"`
	Zip         string  `db:"zip" json:"zip"`
	Email       string  `db:"email" json:"email"`
	Phone       int64   `db:"phone" json:"phone"`
	Source      float64 `db: "source" json:"source"`
}

type CheckedProvider struct {
	Id                      int64          `db:"id"`
	ProviderId              int64          `db:"provider_id"`
	FbPageUrl               string         `db:"fb_page_url"`
	FbMatchedPhone          string         `db:"fb_matched_phone"`
	FbLikes                 int            `db:"fb_likes"`
	YelpRating              float64        `db:"yelp_rating"`
	YelpReviewCount         int            `db:"yelp_review_count"`
	UserWebsite             string         `db:"fc_user_website"`
	TwitterId               string         `db:"fc_twitter_id"`
	TwitterFollowers        int            `db:"fc_twitter_followers"`
	LinkedinId              string         `db:"fc_linkedin_id"`
	InMasterList            string         `db:"in_master_list"`
	LnBizMatch              string         `db:"ln_biz_match"`
	LnBizMatchedPhone       string         `db:"ln_biz_matched_phone"`
	LnBizRawJson            sql.NullString `db:"ln_biz_raw_json"`
	LnIndivMatch            string         `db:"ln_indiv_match"`
	LnIndivMatchedPhone     string         `db:"ln_indiv_matched_phone"`
	LnIndivMatchedZip       string         `db:"ln_indiv_matched_zip"`
	LnIndivRawJson          sql.NullString `db:"ln_indiv_raw_json"`
	GoogleMatch             string         `db:"google_match"`
	GooglePhoneMatch        string         `db:"google_phone_match"`
	LnBizInstidMatch        string         `db:"ln_bizinstid_match"`
	LnBizInstidMatchedPhone string         `db:"ln_bizinstid_matched_phone"`
	LnBizInstidMatchedName  string         `db:"ln_bizinstid_matched_name"`
	LnBizInstidMatchedZip   string         `db:"ln_bizinstid_matched_zip"`
	FcBizMatch              string         `db:"fc_biz_match"`
	FcBizMatchedOrg         string         `db:"fc_biz_matched_org"`
	FcBizFbPage             sql.NullString `db:"fc_biz_fb_page"`
	FcBizTwitterId          sql.NullString `db:"fc_biz_twitter_id"`
	FcBizTwitterFollowers   int            `db:"fc_biz_twitter_followers"`
	FcBizLinkedinMatch      string         `db:"fc_biz_linkedin_match"`
	ExperianMatch           string         `db:"experian_match"`
	ExperianBizId           sql.NullString `db:"experian_biz_id"`
	BgcStatus               string         `db:"bgc_status"`
}

type BizSearch struct {
	Id                int64  `db:"id"`
	ProviderId        int64  `db:"provider_id"`
	LnBizMatch        string `db:"ln_biz_match"`
	LnBizMatchedPhone string `db:"ln_biz_matched_phone"`
	LnBizRawJson      string `db:"ln_biz_raw_json"`
}

type ThumbtackProviderData struct {
	Id                   int64   `db:"id"`
	ProviderId           int64   `db:"provider_id"`
	Bids                 int     `db:"number_bids"`
	Hires                int     `db:"number_hires"`
	Reviews              int     `db:"number_reviews"`
	HasValidWebsite      string  `db:"is_valid_website"`
	AverageReviewScore   float64 `db:"av_review_score"`
	Complaints           int     `db:"number_complaints"`
	AverageReimbursement float64 `db:"av_reimbursement"`
	AccountDeleted       string  `db:"account_deleted"`
	BackgroundVerified   string  `db:"backgd_verified"`
	Scored               string  `db:"scored"`
}
