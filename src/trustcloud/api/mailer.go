package api

import (
	"bytes"
	"net/smtp"
	"strconv"
	"strings"
	"text/template"
	"trustcloud/util"
)

var emailAuth smtp.Auth

func init() {
	emailAuth = smtp.PlainAuth("",
		util.Configuration.Environment.Mail.User,
		util.Configuration.Environment.Mail.Password,
		util.Configuration.Environment.Mail.Server,
	)
}

/*
Send a mail message to the user after project is created/saved
*/

func SendProjectCreatedMail(project *Project) error {
	return SendMail(GetMailTemplateData(project, "", nil), "created")
}

/*
Send a mail message to the user after project is approved
*/

func SendProjectApprovedMail(project *Project, webServer string) error {
	project.Server = webServer

	return SendMail(GetMailTemplateData(project, webServer, nil), "approved")
}

/*
Send a mail message to the user after project is declined
*/

func SendProjectDeclinedMail(project *Project, webServer string, providers []Provider) error {
	project.Server = webServer

	return SendMail(GetMailTemplateData(project, webServer, providers), "declined")
}

func SendPurchaseCompleteMail(project *Project) error {

	return SendMail(GetMailTemplateData(project, "", nil), "complete")
}

func GetMailTemplateData(project *Project, webServer string, providers []Provider) *MailTemplateData {
	templateData := &MailTemplateData{}

	SetDisplayableFields(project)

	templateData.Project = *project

	templateData.Providers = providers

	return templateData
}

/*
Send a mail message to the user
*/

func SendMail(templateData *MailTemplateData, mailType string) error {
	recipient := templateData.Project.User.Email

	util.InfoLog.Println("Sending email type [", mailType, "] to ", recipient)

	return SendMailFromTemplate(
		"templates/"+util.Configuration.General.MailTemplate[mailType].Name,
		util.Configuration.General.MailTemplate[mailType].Subject,
		recipient,
		templateData)
}

func SendMailFromTemplate(templateName string, subject string, recipient string, data interface{}) error {
	var err error
	var doc bytes.Buffer

	t, err := template.ParseFiles(templateName)

	if err != nil {
		util.ErrorLog.Println("error trying to parse mail template ", templateName)
	}

	err = t.Execute(&doc, data)
	if err != nil {
		util.ErrorLog.Println("error trying to execute mail template ", templateName)
	}

	util.InfoLog.Println("Send to recipient: ", recipient)

	return SendMailMessage(recipient, subject, doc.Bytes())
}

/*
Send an email message
*/
func SendMailMessage(recipient string, subject string, body []byte) error {

	if DoSend(recipient) {
		message := "To: " + recipient + "\r\nSubject: " +
			subject + "\r\n\r\n" + string(body)

		util.InfoLog.Println("Email sent")

		return smtp.SendMail(util.Configuration.Environment.Mail.Server+
			":"+strconv.Itoa(util.Configuration.Environment.Mail.Port),
			emailAuth, util.Configuration.Environment.Mail.From,
			[]string{recipient},
			[]byte(message))
	}

	if !DoSend(recipient) {
		util.InfoLog.Println("Test mode - no email sent")
	}

	return nil
}

/*
If the TestMode setting is on - emails will only be sent to a @trustcloud.com account, ignored otherwise
*/
func DoSend(recipient string) bool {
	mode := IsTestMode()

	if mode == true {
		if strings.HasSuffix(recipient, "trustcloud.com") {
			mode = false
		}
	}

	return !mode
}

/*
If the TestMode setting is on - emails will only be sent to a @trustcloud.com account, ignored otherwise
*/
func IsTestMode() bool {
	return util.Configuration.Environment.Mail.TestMode == "Y"
}
