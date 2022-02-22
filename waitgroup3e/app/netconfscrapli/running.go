package netconfscrapli

import (
	"fmt"
	"main/app/inventory"
	"main/app/tasks"
	"github.com/scrapli/scrapligo/netconf"
)

type Running struct {
	Name string
	Kwargs map[string]interface{}
	Include map[string][]string
	Exclude map[string][]string
}

func (s *Running) Task() tasks.TaskBase {
	return tasks.TaskBase{
		Name: s.Name,
		Include: s.Include,
		Exclude: s.Exclude,
	}
}

func (s *Running) Run(h *inventory.Host, c *netconf.Driver, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name
	
	// ==== Custom
	fmt.Printf("%v - args: %+v\n",h.Name, s.Kwargs)
	if len(prev_results)>=1 {
		fmt.Printf("%v - previous result: %+v\n",h.Name, prev_results[len(prev_results)-1])
	}

	r, err := c.GetConfig("running")
	if err != nil {
		return res, fmt.Errorf("failed to get config; error: %+v", err)
	}

	res["result"] = r.Result

	// === Required
	return res, nil

}
