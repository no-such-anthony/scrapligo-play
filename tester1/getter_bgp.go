package tester1


import (
	"fmt"
	//"log"
	"github.com/scrapli/scrapligo/driver/network"
)


func Bgp(h Host, c *network.Driver) (map[string]interface{}, error) {

	res := make(map[string]interface{})

	rs, err := c.SendCommand("show ip bgp summary")
	if err != nil {
		return res, fmt.Errorf("failed to send command for %s: %+v", h.Hostname, err)
	}

	parsedOut, err := rs.TextFsmParse("cisco_ios_show_ip_bgp_summary.textfsm")
	if err != nil {
		return res, fmt.Errorf("failed to parse command for %s: %+v", h.Hostname, err)
	}

	//TODO: check parsedOut len and return error
	res["result"] = parsedOut

	return res, nil

}