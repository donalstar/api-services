package util

import (
	"log"
	"os"
)

var (
	TraceLog   *log.Logger
	InfoLog    *log.Logger
	WarningLog *log.Logger
	ErrorLog   *log.Logger
	WebLog     *log.Logger
)

var logDir = "/var/log/trustcloud"
var logFile = "trustcloud-api-services.log"
var webLogFile = "api-services-server.log"

func init() {
	inf := logDir + "/" + logFile

	file, err := os.OpenFile(inf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", file, ":", err)
	}

	webLogFile, err := os.OpenFile(logDir+"/"+webLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open web log file", webLogFile, ":", err)
	}

	InfoLog = log.New(file,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	WarningLog = log.New(file,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLog = log.New(file,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	WebLog = log.New(webLogFile,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
