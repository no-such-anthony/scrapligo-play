package sshnetrasp

import (
	"fmt"
	"io/ioutil"
	"main/play/app"
	"github.com/networklore/netrasp/pkg/netrasp"
	"github.com/sirikothe/gotextfsm"
	"context"
	//"time"
)



type SendCommand struct {
	Name string
	Command string
	Textfsm string
	Include map[string][]string
	Exclude map[string][]string
}

func (s *SendCommand) Task() app.TaskBase {
	return app.TaskBase{
		Name: s.Name,
		Include: s.Include,
		Exclude: s.Exclude,
	}
}

func (s *SendCommand) Run(h *app.Host, c netrasp.Platform, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name

	// ==== Custom

	if s.Command == "" {
		res["result"] = "SendCommand: no command to run"
		res["failed"] = true
		return res, fmt.Errorf("SendCommand: no command to run")
	}


	//ctx, cancelRun := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	//defer cancelRun()
	rs, err := c.Run(context.Background(), s.Command)
	if err != nil {
		return res, fmt.Errorf("unable to run command: %v\n", err)
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

