package sshgomiko

import (
	"fmt"
	"io/ioutil"
	"main/play/app"
	"github.com/Ali-aqrabawi/gomiko/pkg/types"
	"github.com/sirikothe/gotextfsm"
)


type SendCommand struct {
	app.TaskBase
	Command string
	Textfsm string
}

func (s *SendCommand) Info() app.TaskBase {
	return s.TaskBase
}

func (s *SendCommand) Run(h *app.Host, c types.Device, prev_results []map[string]interface{}) (map[string]interface{}, error) {

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

		textfsm, err := ioutil.ReadFile(s.Textfsm)
		if err != nil {
			res["result"] = fmt.Sprintf("SendCommand: Error opening template '%s'",s.Textfsm)
			res["failed"] = true
			res["error"] = err.Error()
			return res, fmt.Errorf("error opening template '%s'", err.Error())
		}

		fsm := gotextfsm.TextFSM{}
		err = fsm.ParseString(string(textfsm))
		if err != nil {
			res["result"] = fmt.Sprintf("SendCommand: Error while parsing template '%s'",s.Textfsm)
			res["failed"] = true
			res["error"] = err.Error()
			return res, fmt.Errorf("error while parsing template '%s'", err.Error())
		}
		
		parser := gotextfsm.ParserOutput{}
		err = parser.ParseTextString(rs, fsm, true)
		if err != nil {
			res["result"] = fmt.Sprintf("SendCommand: Error while parsing input for template '%s'",s.Textfsm)
			res["failed"] = true
			res["error"] = err.Error()
			return res, fmt.Errorf("error while parsing input '%s'", err.Error())
		}
		res["result"] = parser.Dict
	} else {
		res["result"] = rs
	}

	// ====== Required
	return res, nil

}

