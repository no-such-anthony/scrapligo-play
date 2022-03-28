package tester1


import (
	"fmt"
	//"log"
	"github.com/scrapli/scrapligo/driver/network"
)


func Version(h Host, c *network.Driver) (map[string]interface{}, error) {


	res := make(map[string]interface{})

	rs, err := c.SendCommand("show version")
	if err != nil {
		return res, fmt.Errorf("failed to send command for %s: %+v", h.Hostname, err)
	}

	parsedOut, err := rs.TextFsmParse("cisco_ios_show_version.textfsm")
	if err != nil {
		return res, fmt.Errorf("failed to parse command for %s: %+v", h.Hostname, err)
	}

	//TODO: check parsedOut len and return error

	res["result"] = parsedOut[0]

	return res, nil

}