package sshscrapli

import (
	"fmt"
	"main/play/app"
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


func (s *ScrapligoSsh) Open(h *app.Host) (error) {

	sshport := h.Port

	if sshport == 0 {
		sshport = 22
	}

	c, err := core.NewCoreDriver(
		h.Hostname,
		h.Platform,
		base.WithAuthStrictKey(h.StrictKey),
		base.WithAuthUsername(h.Username),
		base.WithAuthPassword(h.Password),
		base.WithPort(sshport),
		//base.WithAuthSecondary(h.Enable),
		//base.WithTransportType("standard"),
		//base.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		return fmt.Errorf("ssh: failed to create driver: %+v", err)
	}

	err = c.Open()
	if err != nil {
		return fmt.Errorf("ssh: failed to open driver: %+v", err)
	}

	s.C = c
	return nil 

}

func GetConn(h *app.Host) (app.Connector, error) {

	var cc app.Connector

	conn, err := h.GetConnection("scrapli_ssh")
	if err == nil {
		return conn, nil
	}

	cc = app.Connector(&ScrapligoSsh{})
	err = cc.Open(h)

	if err != nil {
		return cc, err
	}

	h.SetConnection("scrapli_ssh", cc)
	return cc, nil
	
}
