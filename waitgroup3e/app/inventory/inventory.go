package inventory

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
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
	StrictKey bool
	Groups    []string
	Connections map[string]Connector
	Data      map[string]interface{}
}


type Hosts map[string]*Host

// SetConnections stores a connection
func (h *Host) SetConnection(name string, conn Connector) {
	if h.Connections == nil {
		h.Connections = make(map[string]Connector)
	}
	h.Connections[name] = conn
}

// GetConnection retrieves a connection that was previously set
func (h *Host) GetConnection(name string) (Connector, error) {
	if h.Connections == nil {
		h.Connections = make(map[string]Connector)
	}
	if c, ok := h.Connections[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("couldn't find connection")
}


func GetHostsByYAML() Hosts {

	var i Hosts
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
				v.Username = "developer"
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


func GetHostsByCode() Hosts {

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
		host.Data["example_only"] = 100
		hosts[host.Name] = &host
	}

	var host Host
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