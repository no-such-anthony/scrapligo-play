package main
// a netrasp plugin only example
// creating your own task within your script

import (
	"fmt"
	"time"
	"main/play/app"
	"main/play/plugins/sshnetrasp"
	"main/play/plugins/inventory"
	//extra imports for loading your own, or modifying existing
	"github.com/networklore/netrasp/pkg/netrasp"
	"context"
	"io/ioutil"
	"github.com/sirikothe/gotextfsm"
)

func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := inventory.GetHostsByYAML()
	//fmt.Println(hosts)

	command := "show version"
	textfsm := "../textfsm_templates/cisco_iosxe_show_version.textfsm"

	//notice the use of local sendcommand attached to the wrapper
	task1 := sendCommand{
		TaskBase: app.TaskBase{
			Name: "a show version in netrasp",
			//Exclude: map[string][]string{"name": []string{"sandbox"}},
		},
		Command: command,
		Textfsm: textfsm,
	}
	wtask1 := sshnetrasp.Wrap{&task1}

	t := []app.Play{&wtask1}
	//fmt.Printf("%+v\n", t)

	results := app.Runner(hosts, t)

	fmt.Print("\n\n")
	fmt.Println("======================= RESULTS =================================")
	for n, h := range results {
		fmt.Println("Name:", n)
		for _, res := range h.([]map[string]interface{}) {
			fmt.Println(res)
		}
		fmt.Print("\n\n")
	}

}


func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}

// this example is a creating a new sendcommand local to this file

type sendCommand struct {
	app.TaskBase
	Command string
	Textfsm string
}

func (s *sendCommand) Info() app.TaskBase {
	return s.TaskBase
}

func (s *sendCommand) Run(h *app.Host, c netrasp.Platform, prev_results []map[string]interface{}) (map[string]interface{}, error) {

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