package app

import (
	"fmt"
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

