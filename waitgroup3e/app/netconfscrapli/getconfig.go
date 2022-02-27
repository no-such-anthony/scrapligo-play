package netconfscrapli

import (
	"fmt"
	"main/app/inventory"
	"main/app/tasks"
	"github.com/scrapli/scrapligo/netconf"
	"github.com/go-xmlfmt/xmlfmt"
)

type GetConfig struct {
	Name string
	Type string  	// running, startup, candidate...
	Filter string	// netconf filter
	Include map[string][]string
	Exclude map[string][]string
}

func (s *GetConfig) Task() tasks.TaskBase {
	return tasks.TaskBase{
		Name: s.Name,
		Include: s.Include,
		Exclude: s.Exclude,
	}
}

func (s *GetConfig) Run(h *inventory.Host, c *netconf.Driver, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name
	
	// ==== Custom
	if s.Type == "" {
		s.Type = "running"
	}

	r, err := c.GetConfig(s.Type, netconf.WithNetconfFilter(s.Filter))
	if err != nil {
		return res, fmt.Errorf("failed to get config; error: %+v", err)
	}

	res["result"] = xmlfmt.FormatXML(r.Result, "", "  ")

	// === Required
	return res, nil

}