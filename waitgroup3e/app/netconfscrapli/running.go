package netconfscrapli

import (
	"fmt"
	"main/app/inventory"
	"main/app/tasks"
	"github.com/scrapli/scrapligo/netconf"
	"github.com/go-xmlfmt/xmlfmt"
)

type Running struct {
	Name string
	NcFilter string
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
	r, err := c.GetConfig("running", netconf.WithNetconfFilter(s.NcFilter))
	if err != nil {
		return res, fmt.Errorf("failed to get config; error: %+v", err)
	}

	res["result"] = xmlfmt.FormatXML(r.Result, "", "  ")

	// === Required
	return res, nil

}
