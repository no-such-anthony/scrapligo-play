package connections

import (
	"main/app/inventory"
)


func GetConn(h *inventory.Host, m string) (inventory.Connector, error) {

	var cc inventory.Connector

	conn, err := h.GetConnection(m)
	if err == nil {
		return conn, nil
	}

	switch m {
	case "scrapli_ssh": 
		cc = inventory.Connector(&ScrapligoSsh{})
	case "scrapli_netconf": 
		cc = inventory.Connector(&ScrapligoNetconf{})
	default:
		cc = inventory.Connector(&ScrapligoSsh{})
	}

	err = cc.Open(h)

	if err != nil {
		return cc, err
	}

	h.SetConnection(m, cc)
	return cc, nil
	
}