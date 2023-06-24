package inventory

import (
	"main/play/app"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"fmt"
)

func GetHostsByYAML() app.Hosts {

	var i app.Hosts
	yamlFile, err := ioutil.ReadFile("hosts.yaml")
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.UnmarshalStrict(yamlFile, &i)
	if err != nil {
		fmt.Println(err)
	}

	// use group creds and if empty Data map, create.
	for h,v := range(i) {

		v.Name = h

		for _, g := range v.Groups {
			if g == "gns" {
				v.Username = "fred"
				v.Password = "bedrock"
				v.Platform = "cisco_iosxe"
				v.StrictKey = false
				break
			} else if g == "devnet" {
				v.Username = "admin"
				v.Password = "C1sco12345"
				v.Platform = "cisco_iosxe"
				v.StrictKey = false
				break
			}
		}

		if len(v.Data) == 0 {
			v.Data = make(map[string]interface{})
		}
	}
	return i

}


func GetHostsByCode() app.Hosts {

	devices := []string{"no.suchdomain","192.168.204.101","192.168.204.102","192.168.204.103","192.168.204.104"}
	hosts := make(app.Hosts)

	for _,value := range devices {
		var host app.Host
		host.Data = make(map[string]interface{})
		host.Name = value
		host.Hostname = value
		host.Platform = "cisco_iosxe"
		host.Username = "fred"
		host.Password = "bedrock"
		host.StrictKey = false
		host.Data["example_only"] = 100
		hosts[host.Name] = &host
	}

	var host app.Host
	host.Data = make(map[string]interface{})
	host.Name = "sandbox"
	host.Hostname = "sandbox-iosxe-latest-1.cisco.com"
	host.Port = 830
	host.Platform = "cisco_iosxr"
	host.Username = "developer"
	host.Password = "C1sco12345"
	host.StrictKey = false
	host.Data["example_only"] = 100
	hosts["sandbox"] = &host



	return hosts
}