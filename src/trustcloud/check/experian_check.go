package check

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"trustcloud/api"
	"trustcloud/util"
)

type experianChecker struct {
}

func (db *experianChecker) GetName() string {
	return "Experian"
}

func (db *experianChecker) Check(provider *api.Provider, checkedProvider *api.CheckedProvider) {
	body, _ := ioutil.ReadFile(util.BaseDir + "/test/" + "net-connect-prod.xml")

	input := "&NETCONNECT_TRANSACTION=" + url.QueryEscape(string(body))

	util.InfoLog.Println("Use input ", input)

	baseUrl := getLookupLink()

	client := &http.Client{}

	util.InfoLog.Println("Lookup link: ", baseUrl)

	user := util.Configuration.General.Connector["experian"].Id
	password := util.Configuration.General.Connector["experian"].Password

	data := url.Values{}
	data.Set("NETCONNECT_TRANSACTION", string(body))

	//	client := &http.Client{}
	//	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	//
	//
	//

	//	test := "&NETCONNECT_TRANSACTION=%3CNetConnectRequest%3E%0A++++%3CEAI%3EDEWJEIHM%3C%2FEAI%3E%0A++++%3CDBHost%3EBIZ_ID%3C%2FDBHost%3E%0A++++%3CReferenceId%3EBizID+Sample%3C%2FReferenceId%3E%0A++++%3CRequest%3E%0A++++++++%3CProducts%3E%0A++++++++++++%3CBizID%3E%0A++++++++++++++++%3CXMLVersion%3E03%3C%2FXMLVersion%3E%0A++++++++++++++++%3CSubscriber%3E%0A++++++++++++++++++++%3CPreamble%3ETBRC%3C%2FPreamble%3E%0A++++++++++++++++++++%3COpInitials%3E00%3C%2FOpInitials%3E%0A++++++++++++++++++++%3CSubCode%3E0308560%3C%2FSubCode%3E%0A++++++++++++++++%3C%2FSubscriber%3E%0A++++++++++++++++%3CBusinessApplicant%3E%0A++++++++++++++++++++%3CBusinessName%3EA+to+Z+Painting+Plus%3C%2FBusinessName%3E%0A++++++++++++++++++++%3CCurrentAddress%3E%0A++++++++++++++++++++++++%3CStreet%3E133+Airport+Rd%3C%2FStreet%3E%0A++++++++++++++++++++++++%3CCity%3EConcord%3C%2FCity%3E%0A++++++++++++++++++++++++%3CState%3ENH%3C%2FState%3E%0A++++++++++++++++++++++++%3CZip%3E03301%3C%2FZip%3E%0A++++++++++++++++++++%3C%2FCurrentAddress%3E%0A++++++++++++++++++++%3CPhone%3E%0A++++++++++++++++++++++++%3CNumber%3E6038564270%3C%2FNumber%3E%0A++++++++++++++++++++%3C%2FPhone%3E%0A++++++++++++++++%3C%2FBusinessApplicant%3E%0A++++++++++++++++%3COptions%3E%0A++++++++++++++++++++%3CProductOption%3E100%3C%2FProductOption%3E%0A++++++++++++++++++++%3CReferenceNumber%3ESample100%3C%2FReferenceNumber%3E%0A++++++++++++++++++++%3CVerbose%3EY%3C%2FVerbose%3E%0A++++++++++++++++%3C%2FOptions%3E%0A++++++++++++%3C%2FBizID%3E%0A++++++++%3C%2FProducts%3E%0A++++%3C%2FRequest%3E%0A%3C%2FNetConnectRequest%3E%0A0A"
	//	buffer := bytes.NewBufferString(test)

	//	req, err := http.NewRequest("POST", baseUrl, buffer)
	req, err := http.NewRequest("POST", baseUrl, bytes.NewBufferString(data.Encode()))
	req.SetBasicAuth(user, password)

	res, err := client.Do(req)
	if err != nil {
		util.ErrorLog.Println("Error : %s", err)
	}

	body, _ = ioutil.ReadAll(res.Body)

	util.InfoLog.Println("RESP -- ", string(body))
}

func getLookupLink() string {
	user := util.Configuration.General.Connector["experian"].Id
	password := util.Configuration.General.Connector["experian"].Password
	baseUrl := util.Configuration.General.Connector["experian"].BaseUrl

	baseUrl = baseUrl +
		"?lookupServiceName=AccessPoint&lookupServiceVersion=1.0" +
		"&serviceName=NetConnect&serviceVersion=2.0&responseType=text/plain"

	client := &http.Client{}

	req, err := http.NewRequest("POST", baseUrl, nil)
	req.SetBasicAuth(user, password)

	res, err := client.Do(req)
	if err != nil {
		util.ErrorLog.Println("Error : %s", err)
	}

	body, _ := ioutil.ReadAll(res.Body)

	return string(body)
}
