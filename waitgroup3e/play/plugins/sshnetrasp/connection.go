package sshnetrasp

import (
	"fmt"
	"main/play/app"
	"github.com/networklore/netrasp/pkg/netrasp"
	"context"
	"time"
)

type NetraspSsh struct {
	C netrasp.Platform
}


func (s NetraspSsh) Close() {
	s.C.Close(context.Background())
}


func (s *NetraspSsh) Open(h *app.Host) (error) {

	sshport := h.Port

	if sshport == 0 {
		sshport = 22
	}

	device, err := netrasp.New(h.Hostname,
		netrasp.WithUsernamePassword(h.Username, h.Password),
		netrasp.WithDriver("ios"),
		netrasp.WithSSHCipher("aes128-cbc"),
		netrasp.WithInsecureIgnoreHostKey(),
	)
	if err != nil {
		return fmt.Errorf("unable to create client: %v", err)
	}

	ctx, cancelOpen := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancelOpen()
	err = device.Dial(ctx)
	if err != nil {
		return fmt.Errorf("unable to connect: %v\n", err)
	}

	s.C = device
	return nil

}

func GetConn(h *app.Host) (app.Connector, error) {

	var cc app.Connector

	conn, err := h.GetConnection("netrasp_ssh")
	if err == nil {
		return conn, nil
	}

	cc = app.Connector(&NetraspSsh{})
	err = cc.Open(h)

	if err != nil {
		return cc, err
	}

	h.SetConnection("netrasp_ssh", cc)
	return cc, nil
	
}
