package tasks

import (
	"fmt"
	"main/app/inventory"
	"crypto/tls"
	"io/ioutil"
	"net/http"
)

type TestRestConf struct {
	Name string
	Include map[string][]string
	Exclude map[string][]string
}

func (s *TestRestConf) Task() TaskBase {
	return TaskBase{
		Name: s.Name,
		Include: s.Include,
		Exclude: s.Exclude,
	}
}

func (s *TestRestConf) Run(h *inventory.Host, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name
	
	// ==== Custom
	certInfo := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: certInfo}
	url := "https://" + h.Hostname + ":443/restconf/data/Cisco-IOS-XE-native:native?content=config&depth=65535"

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
