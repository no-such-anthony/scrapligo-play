package inventory

import (
	"github.com/scrapli/scrapligo/driver/network"
)


type Host struct {
	Name      string
	Hostname  string
	Platform  string
	Username  string
	Password  string
	Enable    string
	StrictKey bool
	Groups    []string
	Connection *network.Driver
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
		host.Data["example_only"] = 100
		hosts[host.Name] = &host
	}
	return hosts
}