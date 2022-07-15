package sshscrapli

import (
	"fmt"
	"main/play/app"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
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

	p, err := platform.NewPlatform(
		h.Platform,
		h.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(h.Username),
		options.WithAuthPassword(h.Password),
		options.WithSSHConfigFile("../inventory/ssh_config"),
	)
	if err != nil {
		return fmt.Errorf("ssh: failed to create platform for %s: %+v\n\n", h.Hostname, err)
	}

	c, err := p.GetNetworkDriver()
	if err != nil {
        return fmt.Errorf("ssh: failed to fetch network driver for %s: %+v\n\n", h.Hostname, err)
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
