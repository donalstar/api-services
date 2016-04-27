package legacy

/*
Interacting with the legacy TrustCloud database
*/

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"trustcloud/util"
)

// The Legacy DB interface .
type DB interface {
	GetUser(name string, password string) (*ApiAdminUser, error)
	GetSupportTicket(id string) (*SupportTicketView, error)
	GetLegacyTransaction(id string) (*LegacyTransaction, error)
	AddLegacyTransaction(transaction *LegacyTransaction) (int64, error)
	GetCannedResponses() []interface{}
	AddCannedResponse(cannedResponse *CannedResponse) (int64, error)
	DeleteCannedResponse(cannedResponse *CannedResponse) error
}

type legacyDatabase struct {
}

// The one and only database instance.
var DataBase DB

var LegacyDbMap gorp.DbMap

func init() {
	DataBase = &legacyDatabase{}
	LegacyDbMap = *initDb()
}

/**

*/
func initDb() *gorp.DbMap {

	db, err := sql.Open("mysql",
		util.Configuration.Environment.Database.User+
			":"+util.Configuration.Environment.Database.Password+
			"@tcp("+util.Configuration.Environment.Database.Host+":3306)/"+
			util.Configuration.Environment.Database.Legacy+"?parseTime=true")

	checkErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	tables := map[interface{}]string{
		ApiAdminUser{}:      "ApiAdminUser",
		CannedResponse{}:    "SupportCannedResponses",
		LegacyTransaction{}: "Transaction",
	}

	for tableType, tableName := range tables {
		dbmap.AddTableWithName(tableType, tableName).SetKeys(true, "Id")
	}

	return dbmap
}

func (db *legacyDatabase) GetUser(name string, password string) (*ApiAdminUser, error) {
	var user ApiAdminUser

	util.InfoLog.Println("Get user record for ", name)

	error := LegacyDbMap.SelectOne(&user, "select * from ApiAdminUser where strUsername = ? and strPassword = ?",
		name, password)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
	}

	return &user, error
}

func (db *legacyDatabase) GetSupportTicket(id string) (*SupportTicketView, error) {
	var supportTicket SupportTicketView

	error := LegacyDbMap.SelectOne(&supportTicket,
		"SELECT SupportTickets.*,ApiAdminUser.strUsername,ApiAdminUser.intApiAccount "+
			"FROM SupportTickets "+
			"INNER JOIN LNApiUsage on SupportTickets.strTransactionID = LNApiUsage.strTransactionID "+
			"LEFT JOIN ApiAdminUser on LNApiUsage.intUserId = ApiAdminUser.intApiAccount "+
			"WHERE SupportTickets.strTransactionID=? "+
			"AND LNApiUsage.intQueryType = 15 LIMIT 1",
		id)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
	}

	return &supportTicket, error
}

func (db *legacyDatabase) GetLegacyTransaction(id string) (*LegacyTransaction, error) {
	var transaction LegacyTransaction

	util.InfoLog.Println("Get transaction record for ", id)

	error := LegacyDbMap.SelectOne(&transaction, "select * from Transaction where idTransaction = ?",
		id)

	if error != nil {
		util.ErrorLog.Println("DB error...", error)
	}

	return &transaction, error
}

func (db *legacyDatabase) AddLegacyTransaction(transaction *LegacyTransaction) (int64, error) {
	util.InfoLog.Println("AddLegacyTransaction ", transaction)

	err := LegacyDbMap.Insert(transaction)

	util.InfoLog.Println("AddLegacyTransaction - result ", err)

	return transaction.Id, err
}

func (db *legacyDatabase) GetCannedResponses() []interface{} {
	var cannedResponses []CannedResponse
	_, err := LegacyDbMap.Select(&cannedResponses, "select * from SupportCannedResponses")

	checkErr(err, "sqlselect error")

	return cannedResponsesToIface(cannedResponses)
}

func (db *legacyDatabase) AddCannedResponse(cannedResponse *CannedResponse) (int64, error) {
	util.InfoLog.Println("AddCannedResponse ", cannedResponse)

	err := LegacyDbMap.Insert(cannedResponse)

	util.InfoLog.Println("AddCannedResponse - after ", err)

	return cannedResponse.Id, err
}

func (db *legacyDatabase) DeleteCannedResponse(cannedResponse *CannedResponse) error {
	_, err := LegacyDbMap.Delete(cannedResponse)

	return err
}

func checkErr(err error, msg string) {
	if err != nil {
		util.ErrorLog.Println("error --  ...!\n", err)
	}
}
