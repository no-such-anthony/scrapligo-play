package tasks

import (
	"fmt"
	"main/app/inventory"
	"main/app/connections"
)


type NetconfShowVersion struct {
	Name string
	Kwargs map[string]interface{}
	Include map[string][]string
	Exclude map[string][]string
}


func (s *NetconfShowVersion) Named() string {
	return fmt.Sprint(s.Name)
}


func (s *NetconfShowVersion) Run(h *inventory.Host, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	res["task"] = s.Name

	if inventory.Skip(h, s.Include, s.Exclude) {
		res["skipped"] = true
		return res, nil
	}

	conn, err := connections.GetConn(h, "scrapli_netconf")
	if err != nil {
		res["result"] = err
		res["failed"] = true
		return res, err	
	}

	c := conn.(*connections.ScrapligoNetconf).C
	

	fmt.Printf("%v - args: %+v\n",h.Name, s.Kwargs)
	if len(prev_results)>=1 {
		fmt.Printf("%v - previous result: %+v\n",h.Name, prev_results[len(prev_results)-1])
	}

	r, err := c.GetConfig("running")
	if err != nil {
		return res, fmt.Errorf("failed to get config; error: %+v", err)
	}

	res["result"] = r.Result

	return res, nil

}
