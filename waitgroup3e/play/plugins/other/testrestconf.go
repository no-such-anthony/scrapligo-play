package other

import (
	"fmt"
	"main/play/app"
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

type TestRestConf struct {
	app.TaskBase
	Filter string	// restconf filter
}

func (s *TestRestConf) Info() app.TaskBase {
	return s.TaskBase
}

func (s *TestRestConf) Run(h *app.Host, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name
	
	// ==== Custom
	certInfo := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: certInfo}
	url := "https://" + h.Hostname + ":443/restconf/data/Cisco-IOS-XE-native:native/" + 
				s.Filter + "?content=config&depth=65535"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/yang-data+json")
	req.Header.Add("Accept", "application/yang-data+json")
	req.SetBasicAuth(h.Username, h.Password)

	rcRes, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		res["result"] = err
		res["failed"] = true
		return res, err
	}
	defer rcRes.Body.Close()
	body, _ := ioutil.ReadAll(rcRes.Body)
	res["result"] = string(body)

	// === Required
	return res, nil

}
