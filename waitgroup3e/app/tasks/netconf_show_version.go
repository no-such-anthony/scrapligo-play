package tasks

import (
	"fmt"
	"main/app/inventory"
)


type NetconfShowVersion struct {
	Name string
	Kwargs map[string]interface{}
}


func (s *NetconfShowVersion) Named() string {
	return fmt.Sprint(s.Name)
}


func (s *NetconfShowVersion) Run(h *inventory.Host, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	res["name"] = s.Name

	c := h.Connection

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
		return res, fmt.Errorf("failed to send command for %s: %+v", h.Hostname, err)
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil || len(parsedOut) == 0 {
		return res, fmt.Errorf("failed to parse command for %s: %+v", h.Hostname, err)
	}

	res["result"] = fmt.Sprintf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s",
				h.Hostname, parsedOut[0]["HARDWARE"],
				parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])

	return res, nil

}
