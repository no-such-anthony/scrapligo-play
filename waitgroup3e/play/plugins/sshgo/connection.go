package sshgo

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"main/play/app"
	
)

type GoSsh struct {
	C *ssh.Client
}


func (s GoSsh) Close() {
	s.C.Close()
}


func (s *GoSsh) Open(h *app.Host) (error) {

	sshport := h.Port

	if sshport == 0 {
		sshport = 22
	}

	c, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", h.Hostname, sshport), &ssh.ClientConfig{
		User: h.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes128-cbc"},
		},
	})

	if err != nil {
		return fmt.Errorf("ssh: failed to dial: %+v", err)
	}

	s.C = c
	return nil 

}

func GetConn(h *app.Host) (app.Connector, error) {

	var cc app.Connector

	conn, err := h.GetConnection("go_ssh")
	if err == nil {
		return conn, nil
	}

	cc = app.Connector(&GoSsh{})
	err = cc.Open(h)

	if err != nil {
		return cc, err
	}

	h.SetConnection("go_ssh", cc)
	return cc, nil
	
}
