package netconfscrapli

import (
	"fmt"
	"main/play/app"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/netconf"
)

type ScrapligoNetconf struct {
	C *netconf.Driver	
}


func (s ScrapligoNetconf) Close() {
	s.C.Close()
}


func (s *ScrapligoNetconf) Open(h *app.Host) (error) {

	ncport, ok := h.Data["netconf_port"].(int)

	if !ok {
		ncport = 830
	}

	c, err := netconf.NewNetconfDriver(
		h.Hostname,
		base.WithPort(ncport),
		base.WithAuthStrictKey(h.StrictKey),
		base.WithAuthUsername(h.Username),
		base.WithAuthPassword(h.Password),
	)

	if err != nil {
		return fmt.Errorf("netconf: failed to create driver: %+v", err)
	}

	err = c.Open()
	if err != nil {
		return fmt.Errorf("netconf: failed to open driver: %+v", err)
	}

	s.C = c
	return nil 

}


func GetConn(h *app.Host) (app.Connector, error) {

	var cc app.Connector

	conn, err := h.GetConnection("scrapli_netconf")
	if err == nil {
		return conn, nil
	}

	cc = app.Connector(&ScrapligoNetconf{})
	err = cc.Open(h)

	if err != nil {
		return cc, err
	}

	h.SetConnection("scrapli_ssh", cc)
	return cc, nil
	
}