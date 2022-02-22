package netconfscrapli

import (
	"fmt"
	"main/app/inventory"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/netconf"
)

type ScrapligoNetconf struct {
	C *netconf.Driver	
}


func (s ScrapligoNetconf) Close() {
	s.C.Close()
}


func (s *ScrapligoNetconf) Open(h *inventory.Host) (error) {

	c, err := netconf.NewNetconfDriver(
		h.Hostname,
		base.WithPort(h.Port),
		base.WithAuthStrictKey(h.StrictKey),
		base.WithAuthUsername(h.Username),
		base.WithAuthPassword(h.Password),
	)

	if err != nil {
		return fmt.Errorf("failed to create driver for %s: %+v", h.Hostname, err)
	}

	err = c.Open()
	if err != nil {
		return fmt.Errorf("failed to open driver for %s: %+v", h.Hostname, err)
	}

	s.C = c
	return nil 

}


func GetConn(h *inventory.Host) (inventory.Connector, error) {

	var cc inventory.Connector

	conn, err := h.GetConnection("scrapli_ssh")
	if err == nil {
		return conn, nil
	}

	cc = inventory.Connector(&ScrapligoNetconf{})
	err = cc.Open(h)

	if err != nil {
		return cc, err
	}

	h.SetConnection("scrapli_ssh", cc)
	return cc, nil
	
}