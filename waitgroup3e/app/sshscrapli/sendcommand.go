package sshscrapli

import (
	"fmt"
	"main/app/inventory"
	"main/app/tasks"
	"github.com/scrapli/scrapligo/driver/network"
)


type SendCommand struct {
	Name string
	Command string
	Textfsm string
	Include map[string][]string
	Exclude map[string][]string
}

func (s *SendCommand) Task() tasks.TaskBase {
	return tasks.TaskBase{
		Name: s.Name,
		Include: s.Include,
		Exclude: s.Exclude,
	}
}

func (s *SendCommand) Run(h *inventory.Host, c *network.Driver, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name

	// ==== Custom

	if s.Command == "" {
		res["result"] = "SendCommand: no command to run"
		res["failed"] = true
		return res, fmt.Errorf("SendCommand: no command to run")
	}

	rs, err := c.SendCommand(s.Command)
	if err != nil {
		res["result"] = err
		res["failed"] = true
		return res, fmt.Errorf("failed to send command: %+v", err)
	}

	if s.Textfsm != "" {
		parsedOut, err := rs.TextFsmParse(s.Textfsm)
		if err != nil {
			msg := fmt.Sprintf("failed to parse command: %+v", err)
			res["result"] = msg
			res["failed"] = true
			return res, fmt.Errorf(msg)
		}

		if len(parsedOut) == 0 {
			msg := "no output from textfsm parser"
			res["result"] = msg
			res["failed"] = true
			return res, fmt.Errorf(msg)
		}

		res["result"] = parsedOut
	} else {
		res["result"] = rs.Result
	}

	// ====== Required
	return res, nil

}

