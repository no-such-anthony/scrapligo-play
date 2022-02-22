package sshscrapli

import (
	"fmt"
	"main/app/inventory"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
)

type ScrapligoSsh struct {
	C *network.Driver	
}


func (s ScrapligoSsh) Close() {
	s.C.Close()
}


func (s *ScrapligoSsh) Open(h *inventory.Host) (error) {

	c, err := core.NewCoreDriver(
		h.Hostname,
		h.Platform,
		base.WithAuthStrictKey(h.StrictKey),
		base.WithAuthUsername(h.Username),
		base.WithAuthPassword(h.Password),
		//base.WithAuthSecondary(h.Enable),
		//base.WithTransportType("standard"),
		//base.WithSSHConfigFile("ssh_config"),
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

	cc = inventory.Connector(&ScrapligoSsh{})
	err = cc.Open(h)

	if err != nil {
		return cc, err
	}

	h.SetConnection("scrapli_ssh", cc)
	return cc, nil
	
}
