package tasks

import (
	"fmt"
	"main/app/inventory"
	"main/app/connections"
)


type ShowVersion struct {
	Name string
	Method string
	Kwargs map[string]interface{}
	Include map[string][]string
	Exclude map[string][]string
}


func (s *ShowVersion) Named() string {
	return fmt.Sprint(s.Name)
}


func (s *ShowVersion) Run(h *inventory.Host, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required

	res := make(map[string]interface{})
	res["task"] = s.Name

	if inventory.Skip(h, s.Include, s.Exclude) {
		res["skipped"] = true
		return res, nil
	}

	conn, err := connections.GetConn(h, "scrapli_ssh")
	if err != nil {
		res["result"] = err
		res["failed"] = true
		return res, err	
	}

	c := conn.(*connections.ScrapligoSsh).C

	// ==== Custom

	fmt.Printf("%v - args: %+v\n",h.Name, s.Kwargs)
	if len(prev_results)>=1 {
		fmt.Printf("%v - previous result: %+v\n",h.Name, prev_results[len(prev_results)-1])
	}

	cmd := "show version"
	if h.Name == "192.168.204.103" {
		cmd = "show dodgy command"
	}

	rs, err := c.SendCommand(cmd)
	if err != nil {
		res["result"] = err
		res["failed"] = true
		return res, fmt.Errorf("failed to send command for %s: %+v", h.Hostname, err)
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil || len(parsedOut) == 0 {
		msg := fmt.Sprintf("failed to parse command for %s: %+v", h.Hostname, err)
		res["result"] = msg
		res["failed"] = true
		return res, fmt.Errorf(msg)
	}

	if len(parsedOut) == 0 {
		msg := fmt.Sprintf("no output from textfsm parser for %s", h.Hostname)
		res["result"] = msg
		res["failed"] = true
		return res, fmt.Errorf(msg)
	}


	res["result"] = fmt.Sprintf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s",
				h.Hostname, parsedOut[0]["HARDWARE"],
				parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])

	// ====== Required

	return res, nil

}

