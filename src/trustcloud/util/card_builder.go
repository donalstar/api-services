package util

import (
	"bytes"
	"text/template"
)

type Provider2 struct {
	Name  string
	City  string
	State string
	Zip   string
	Score int
}

// GetTrustcard
func GetCard(provider Provider2) string {
	var err error
	var doc bytes.Buffer

	t, err := template.ParseFiles("templates/card.tpl")

	if err != nil {
		ErrorLog.Println("error trying to parse mail template card.tpl", err)
	}

	err = t.Execute(&doc, provider)
	if err != nil {
		ErrorLog.Println("error trying to execute mail template card.tpl", err)
	}

	html := string(doc.Bytes())

	InfoLog.Println("out --  ", html)

	return "document.write('" + html + "')"
}
