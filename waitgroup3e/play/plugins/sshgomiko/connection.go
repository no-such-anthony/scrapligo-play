package sshgomiko

import (
	"fmt"
	"main/play/app"
	"github.com/Ali-aqrabawi/gomiko/pkg"
	"github.com/Ali-aqrabawi/gomiko/pkg/types"
)

type GomikoSsh struct {
	C types.Device	
}


func (s GomikoSsh) Close() {
	s.C.Disconnect()
}


func (s *GomikoSsh) Open(h *app.Host) (error) {

	sshport := h.Port

	if sshport == 0 {
		sshport = 22
	}

	c, err := gomiko.NewDevice(h.Hostname, h.Username, h.Password, "cisco_ios", 22)

	if err != nil {
		return fmt.Errorf("ssh: failed to create driver: %+v", err)
	}

	err = c.Connect()
	if err != nil {
		return fmt.Errorf("ssh: failed to open driver: %+v", err)
	}

	s.C = c
	return nil 

}

func GetConn(h *app.Host) (app.Connector, error) {

	var cc app.Connector

	conn, err := h.GetConnection("gomiko_ssh")
	if err == nil {
		return conn, nil
	}

	cc = app.Connector(&GomikoSsh{})
	err = cc.Open(h)

	if err != nil {
		return cc, err
	}

	h.SetConnection("gomiko_ssh", cc)
	return cc, nil
	
}
