package inventory

import (
)

type Connector interface {
	Open(*Host)  (error)
	Close()
}

type Host struct {
	Name      string
	Hostname  string
	Platform  string
	Port      int
	Username  string
	Password  string
	Enable    string
	Method    string
	StrictKey bool
	Groups    []string
	Connection Connector
	Data      map[string]interface{}
}


type Hosts map[string]*Host


func GetHosts() Hosts {

	devices := []string{"no.suchdomain","192.168.204.101","192.168.204.102","192.168.204.103","192.168.204.104"}
	hosts := make(Hosts)

	for _,value := range devices {
		var host Host
		host.Data = make(map[string]interface{})
		host.Name = value
		host.Hostname = value
		host.Platform = "cisco_iosxe"
		host.Username = "fred"
		host.Password = "bedrock"
		host.StrictKey = false
		//host.Method = "scrapli_ssh"
		host.Data["example_only"] = 100
		hosts[host.Name] = &host
	}

	var host Host
	host.Data = make(map[string]interface{})
	host.Name = "sandbox"
	host.Hostname = "sandbox-iosxe-latest-1.cisco.com"
	host.Port = 830
	//host.Platform = "cisco_iosxe"
	host.Username = "developer"
	host.Password = "C1sco12345"
	host.Method = "scrapli_netconf"
	host.StrictKey = false
	host.Data["example_only"] = 100
	hosts["sandbox"] = &host



	return hosts
}