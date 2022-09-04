package sshgo

import (
	"fmt"
	"bytes"
	"golang.org/x/crypto/ssh"
	"main/play/app"
)


type SendCommand struct {
	app.TaskBase
	Command string
}

func (s *SendCommand) Info() app.TaskBase {
	return s.TaskBase
}

func (s *SendCommand) Run(h *app.Host, c *ssh.Client, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name

	// ==== Custom

	if s.Command == "" {
		res["result"] = "SendCommand: no command to run"
		res["failed"] = true
		return res, fmt.Errorf("SendCommand: no command to run")
	}

	session, err := c.NewSession()
	if err != nil {
		return res, fmt.Errorf("ssh: failed to create session: %+v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(s.Command); err != nil {
		res["result"] = err
		res["failed"] = true
		return res, fmt.Errorf("ssh: failed to execute command: %+v", err)
	}

	res["result"] = make(map[string]string)
	res["result"].(map[string]string)["stdout"] = string(stdout.Bytes())
	res["result"].(map[string]string)["stderr"] = string(stderr.Bytes())

	// ====== Required
	return res, nil

}

