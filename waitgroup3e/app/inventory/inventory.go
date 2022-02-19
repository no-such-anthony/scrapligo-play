package inventory

import "fmt"

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