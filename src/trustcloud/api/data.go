package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"trustcloud/util"
)

var (
	ErrAlreadyExists = errors.New("provider already exists")
)

// The DB interface defines methods to manipulate the providers.
type DB interface {
	GetAllProjects() []interface{}
	GetProject(id int) (*Project, error)
	GetProjectStatus(id int) (*ProjectStatus, error)
	UpdateProject(id int, status string) (*ProjectStatus, error)
	GetProvider(id int) *Provider
	GetAllProviders() []interface{}
	GetAllProviderIds(source string) []int64
	GetApprovedProviders(provider *Provider) []Provider
	FindProvider(businessName string, email string, phone string) []interface{}
	AddProvider(a *Provider) (int64, error)
	AddProject(project *Project) (int, error, ProjectStatus)
	AddProjectStatus(projectStatus *ProjectStatus) (int, error)
	UpdateProjectStatus(projectStatus *ProjectStatus) (int64, error)
	AddJob(job *Job) (int64, error)
	UpdateJobStatus(job *Job) (int64, error)
	AddUser(user *User) int64
	GetJob(id int) *Job
	GetUser(id int) *User
	AddNote(note *Note) (int64, error)
	GetNotes(id int) []*Note
	DeleteNote(note *Note) error
	CreatePayment(payment *Payment) (int64, error)
	Delete(id int)
	GetProviderFromMaster(provider *Provider) int64
	AddCheckedProvider(checkedProvider *CheckedProvider) (int64, error)
	GetCheckedProvider(providerId int64) *CheckedProvider
	UpdateCheckedProvider(checkedProvider *CheckedProvider) (int64, error)
	CheckAuthorizedUser(username string, apiKey string) bool
	GetCheckers() []Checker
	GetCheckerByName(name string) Checker
	LogCheck(name string, providerId int64) (int64, error)
	GetClientPlatform(apiKey string) (*ClientPlatform, error)
}

type projectDatabase struct {
}

// The one and only database instance.
var DataBase DB

var DbMap gorp.DbMap

var Checkers []Checker

func init() {
	DataBase = &projectDatabase{}

	DbMap = *initDb()

	Checkers = DataBase.GetCheckers()

	util.ErrorLog.Println("got checkers  ", Checkers)
}

/**

*/
func initDb() *gorp.DbMap {

	db, err := sql.Open("mysql",
		util.Configuration.Environment.Database.User+
			":"+util.Configuration.Environment.Database.Password+
			"@tcp("+util.Configuration.Environment.Database.Host+":3306)/"+
			util.Configuration.Environment.Database.Name+"?parseTime=true")

	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	tables := map[interface{}]string{
		Provider{}:              "Provider",
		ProjectStatus{}:         "Project",
		Job{}:                   "Job",
		User{}:                  "User",
		Payment{}:               "Payment",
		Note{}:                  "Note",
		CheckedProvider{}:       "CheckedProvider",
		BizSearch{}:             "LnBusinessCheck",
		CheckerLog{}:            "CheckerLog",
		ClientPlatform{}:        "ClientPlatform",
		ClientCredentials{}:     "ClientCredentials",
		ThumbtackProviderData{}: "ThumbtackProviderData",
	}

	for tableType, tableName := range tables {
		dbmap.AddTableWithName(tableType, tableName).SetKeys(true, "Id")
	}

	return dbmap
}

// Add Project
func (db *projectDatabase) AddProject(project *Project) (int, error, ProjectStatus) {
	// New project

	// TODO: Should be in a transaction
	projectStatus := &ProjectStatus{
		Status:     "Pending",
		CreateDate: time.Now().Local(),
	}

	id, err := db.AddProjectStatus(projectStatus)

	if err == nil {
		project.Provider.Source = "USER"

		providerId, _ := db.AddProvider(&project.Provider)

		userId := db.AddUser(&project.User)

		jobId, _ := db.AddJob(&project.Job)

		projectStatus.JobId = jobId
		projectStatus.ProviderId = providerId
		projectStatus.UserId = userId

		projectStatus.EncodedId = EncodeId(id)

		db.UpdateProjectStatus(projectStatus)
	}

	return id, err, *projectStatus

}

func (db *projectDatabase) GetProject(id int) (*Project, error) {
	util.ErrorLog.Println("get proj  ", id)

	projectStatus, error := db.GetProjectStatus(id)

	project := &Project{}

	if error == nil {
		project.Id = projectStatus.Id
		project.EncodedId = projectStatus.EncodedId
		project.Status = projectStatus.Status
		project.CreateDate = projectStatus.CreateDate

		project.Provider = *(db.GetProvider(int(projectStatus.ProviderId)))

		project.User = *(db.GetUser(int(projectStatus.UserId)))

		project.Job = *(db.GetJob(int(projectStatus.JobId)))

		project.SupportPhone = util.Configuration.General.Guarantee.SupportPhone
	}

	return project, error
}

// GetProjectStatus returns the project identified by the id, or nil.
func (db *projectDatabase) GetProjectStatus(id int) (*ProjectStatus, error) {
	var project ProjectStatus
	error := DbMap.SelectOne(&project, "select * from Project where id=?", id)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
	}

	return &project, error
}

func (db *projectDatabase) GetAllProjects() []interface{} {
	util.InfoLog.Println("GetAllProjects")

	var projectSummaries []ProjectSummary
	_, err := DbMap.Select(&projectSummaries, "select Project.id as project_id, Project.status, Project.encoded_id,"+
		"Job.type, Job.cost, Job.startdate, "+
		"Provider.name as provider_name, "+
		"User.first_name, User.last_name "+
		" from Project, Job, Provider, User "+
		" where Job.id = Project.job_id and Provider.id = Project.provider_id and User.id = Project.user_id ")

	checkErr(err, "sqlselect error")

	return projectsToIface(projectSummaries)
}

func (db *projectDatabase) UpdateProject(id int, status string) (*ProjectStatus, error) {
	util.ErrorLog.Println("update proj ", id)

	projectStatus, _ := db.GetProjectStatus(id)

	projectStatus.Status = status

	_, err := db.UpdateProjectStatus(projectStatus)

	return projectStatus, err
}

// GetAll returns all providers from the database.
func (db *projectDatabase) GetAllProviders() []interface{} {
	var providers []Provider
	_, err := DbMap.Select(&providers, "select * from Provider")

	checkErr(err, "sqlselect error")

	return toIface(providers)
}

// GetAll returns all providers from the database.
func (db *projectDatabase) GetAllProviderIds(source string) []int64 {
	var providerIds []int64
	_, err := DbMap.Select(&providerIds, "select id from Provider where source = ?", source)

	checkErr(err, "sqlselect error")

	return providerIds
}

func (db *projectDatabase) GetApprovedProviders(provider *Provider) []Provider {

	var providers []Provider
	_, err := DbMap.Select(&providers,
		"select * from Provider where type = ? and id != ? and city = ? and state = ? limit 3",
		provider.Type, provider.Id, provider.City, provider.State)

	checkErr(err, "sqlselect error")

	if len(providers) < 3 {
		var providers []Provider
		_, err = DbMap.Select(&providers,
			"select * from Provider where type = ? and id != ? limit 3",
			provider.Type, provider.Id)

		checkErr(err, "sqlselect error")
	}

	for i := range providers {
		providers[i].DisplayablePhone = util.FormatPhone(providers[i].Phone)
	}

	return providers
}

// Find returns providers that match the search criteria.
func (db *projectDatabase) FindProvider(name string, email string, phone string) []interface{} {
	query := "select * from Provider"

	params := make(map[string]string)

	if name != "" {
		params["name"] = name
	}

	if email != "" {
		params["email"] = email
	}

	if phone != "" {
		params["phone"] = phone
	}

	if len(params) > 0 {
		prefix := " where "

		i := 0

		for key, _ := range params {
			query += prefix + key + " = :" + key

			if i == 0 {
				i++
				prefix = " and "
			}

		}
	}

	var providers []Provider
	_, err := DbMap.Select(&providers, query, params)

	checkErr(err, "sqlselect error")

	return toIface(providers)
}

func GetNotesArray(notes []Note) []*Note {
	ar := make([]*Note, len(notes))

	i := 0
	for _, v := range notes {
		ar[i] = CopyNote(v)
		i++
	}

	return ar
}

/*
feels dumb
*/

func CopyNote(note Note) *Note {
	copy := &Note{}

	copy.Id = note.Id
	copy.Date = note.Date
	copy.Value = note.Value
	copy.ProjectId = note.ProjectId

	return copy
}

// Get returns the provider identified by the id, or nil.
func (db *projectDatabase) GetProvider(id int) *Provider {
	var provider Provider
	error := DbMap.SelectOne(&provider, "select * from Provider where id=?", id)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
	}

	return &provider
}

func (db *projectDatabase) GetUser(id int) *User {
	var user User
	error := DbMap.SelectOne(&user, "select * from User where id=?", id)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
	}

	return &user
}

func (db *projectDatabase) GetJob(id int) *Job {

	var job Job
	error := DbMap.SelectOne(&job, "select * from Job where id=?", id)

	if error != nil {
		fmt.Printf("DB error...", error)
	}

	// get plan cost from config
	job.PlanCost = GetPlanCost(job.Cost)
	job.PlanName = util.Configuration.General.Guarantee.Name

	return &job
}

// Add creates a new provider and returns its id, or an error.
func (db *projectDatabase) AddProvider(provider *Provider) (int64, error) {
	id := CheckForProvider(provider)

	if id != 0 {
		util.WarningLog.Println("Provider already exists [id = ", id, "]")
	}

	if id == 0 {
		err := DbMap.Insert(provider)

		checkErr(err, "Database insert error")

		id = provider.Id
	}

	return id, nil
}

// Add Project Status
func (db *projectDatabase) AddProjectStatus(projectStatus *ProjectStatus) (int, error) {
	// insert rows - auto increment PKs will be set properly after the insert
	err := DbMap.Insert(projectStatus)

	return projectStatus.Id, err
}

/*
Update Project row with ids for Job, Provider, User
*/
func (db *projectDatabase) UpdateProjectStatus(projectStatus *ProjectStatus) (int64, error) {
	return DbMap.Update(projectStatus)
}

/*
Update Project row with ids for Job, Provider, User
*/
func (db *projectDatabase) UpdateJobStatus(job *Job) (int64, error) {
	return DbMap.Update(job)
}

// Add Job
func (db *projectDatabase) AddJob(job *Job) (int64, error) {
	job.StartDate, _ = util.ParseDate(job.StartDateRaw)

	err := DbMap.Insert(job)

	checkErr(err, "Database insert error")

	return job.Id, err
}

// Add Project Status
func (db *projectDatabase) AddUser(user *User) int64 {

	id := CheckForUser(user)

	if id != 0 {
		util.WarningLog.Println("User already exists [id = ", id, "]")
	}

	if id == 0 {
		// insert rows - auto increment PKs will be set properly after the insert
		err := DbMap.Insert(user)

		checkErr(err, "Database insert error")

		id = user.Id
	}

	return id
}

func (db *projectDatabase) AddNote(note *Note) (int64, error) {
	err := DbMap.Insert(note)

	return note.Id, err
}

func (db *projectDatabase) GetNotes(id int) []*Note {
	var notes []Note
	_, err := DbMap.Select(&notes, "select * from Note where project_id = ?", id)

	checkErr(err, "sqlselect error")

	return GetNotesArray(notes)
}

func (db *projectDatabase) DeleteNote(note *Note) error {
	_, err := DbMap.Delete(note)

	return err
}

func (db *projectDatabase) CreatePayment(payment *Payment) (int64, error) {
	util.InfoLog.Println("Create payment: ", payment)

	err := DbMap.Insert(payment)

	checkErr(err, "Database insert error")

	return payment.Id, err
}

// Delete removes the provider identified by the id from the database. It is a no-op
// if the id does not exist.
func (db *projectDatabase) Delete(id int) {

}

/*
See if Provider already exists
*/
func CheckForProvider(provider *Provider) int64 {
	// use phone # as unique provider identifier
	var ids []int64
	_, err := DbMap.Select(&ids, "select id from Provider where phone = ? and UPPER(name) = UPPER(?) and email = ? and source = ?",
		provider.Phone, provider.Name, provider.Email, provider.Source)

	checkErr(err, "Error retrieving providers")

	var id int64 = 0

	if len(ids) > 0 {
		id = ids[0]
	}

	return id
}

func (db *projectDatabase) GetProviderFromMaster(provider *Provider) int64 {
	util.WarningLog.Println("GetProviderFromMaster for - ", provider.Name, " ", provider.Phone, " ", provider.City)

	var ids []int64
	_, err := DbMap.Select(&ids, "select id from ProviderMaster where UPPER(company_name) = UPPER(?) "+
		"and UPPER(city) = UPPER(?)",
		provider.Name, provider.City)

	checkErr(err, "Error retrieving providers")

	var id int64 = 0

	if len(ids) > 0 {
		id = ids[0]
	}

	return id
}

// Add creates a new provider and returns its id, or an error.
func (db *projectDatabase) AddCheckedProvider(checkedProvider *CheckedProvider) (int64, error) {
	//	util.InfoLog.Println("AddCheckedProvider - ", checkedProvider)

	err := DbMap.Insert(checkedProvider)

	checkErr(err, "Database insert error")

	id := checkedProvider.Id

	return id, err
}

func (db *projectDatabase) GetCheckedProvider(providerId int64) *CheckedProvider {
	var checkedProvider CheckedProvider
	error := DbMap.SelectOne(&checkedProvider, "select * from CheckedProvider where provider_id=?", providerId)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
		return nil
	}

	return &checkedProvider
}

func (db *projectDatabase) CheckAuthorizedUser(username string, apiKey string) bool {
	util.InfoLog.Println("Check authorized user")

	var ids []int64
	_, err := DbMap.Select(&ids, "select a.id from AuthorizedClientUsers a, ClientCredentials c "+
		"where a.client_credentials_id = c.client_platform_id "+
		"and username = ? and api_key = ?", username, apiKey)

	checkErr(err, "Error retrieving users")

	isAuthorized := len(ids) > 0

	util.InfoLog.Println("Authorized? ", isAuthorized)

	return isAuthorized
}

func (db *projectDatabase) UpdateCheckedProvider(checkedProvider *CheckedProvider) (int64, error) {

	util.InfoLog.Println("update checked provider ", checkedProvider, " stat ", checkedProvider.BgcStatus)
	return DbMap.Update(checkedProvider)
}

/*
See if user already exists
*/
func CheckForUser(user *User) int64 {
	// compare name + address

	var ids []int64
	_, err := DbMap.Select(&ids, "select id from User where first_name = ? "+
		"and last_name = ? and address1 = ?", user.FirstName, user.LastName, user.Address1)

	checkErr(err, "Error retrieving users")

	var id int64 = 0

	if len(ids) > 0 {
		id = ids[0]
	}

	return id
}

func (db *projectDatabase) GetCheckers() []Checker {
	var checkers []Checker
	_, err := DbMap.Select(&checkers, "select * from Checker")

	checkErr(err, "sqlselect error")

	return checkers
}

func (db *projectDatabase) GetCheckerByName(name string) Checker {
	var match Checker

	for _, checker := range Checkers {
		if name == checker.Name {
			match = checker
			break
		}
	}

	return match
}

func (db *projectDatabase) LogCheck(name string, providerId int64) (int64, error) {
	c := DataBase.GetCheckerByName(name)

	util.InfoLog.Println("Log with: check id ", c.Id, " prov id ", providerId, " time ", time.Now().Local())

	checkLog := &CheckerLog{
		ProviderId: providerId,
		CheckId:    c.Id,
		CheckTime:  time.Now().Local(),
	}

	err := DbMap.Insert(checkLog)

	checkErr(err, "Database insert error")

	return checkLog.Id, err
}

func (db *projectDatabase) GetClientPlatform(apiKey string) (*ClientPlatform, error) {
	var clientCredentials ClientCredentials
	error := DbMap.SelectOne(&clientCredentials, "select * from ClientCredentials where api_key=?", apiKey)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
	}

	var clientPlatform ClientPlatform
	error = DbMap.SelectOne(&clientPlatform, "select * from ClientPlatform where id=?", clientCredentials.ClientPlatformId)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
	}

	return &clientPlatform, error
}

func (provider *Provider) String() string {
	return fmt.Sprintf("%s %s %s", provider.Name, provider.Email, provider.Phone)
}

func checkErr(err error, msg string) {
	if err != nil {
		util.ErrorLog.Println("error --  ...!\n", err)
	}
}
