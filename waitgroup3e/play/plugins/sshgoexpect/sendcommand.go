package sshgoexpect

import (
	"fmt"
	expect "github.com/google/goexpect"
	"golang.org/x/crypto/ssh"
	"log"
	"time"
	"regexp"
	"main/play/app"
)

const (
	timeout = 10 * time.Minute
)

type SendCommand struct {
	app.TaskBase
	Command string
}

//func (s *SendCommand) Info() app.TaskBase {
//	return s.TaskBase
//}

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

	x, _, err := expect.SpawnSSH(c, timeout)
	if err != nil {
			log.Fatal(err)
	}
	defer x.Close()

	prompt := regexp.MustCompile(".*#")

	x.Expect(prompt, timeout)
	x.Send("term len 0\n")
	x.Expect(prompt, timeout)
	x.Send(s.Command+"\n")
	r, _, _ := x.Expect(prompt, timeout)
	//c.Send("exit\n")
	res["result"] = r

	// ====== Required
	return res, nil

}

