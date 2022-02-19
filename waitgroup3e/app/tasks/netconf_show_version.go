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
	res["name"] = s.Name

	if inventory.Skip(h, s.Include, s.Exclude) {
		res["skipped"] = true
		return res, nil
	}

	conn, ok := h.Connection.(*connections.ScrapligoNetconf)
	if !ok {
		return res, fmt.Errorf("no connection method for %s", h.Hostname)	
	}
	c := conn.C

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
