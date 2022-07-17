package netconfscrapli

import (
	"fmt"
	"main/play/app"
	"github.com/scrapli/scrapligo/driver/netconf"
	"github.com/scrapli/scrapligo/driver/opoptions"
	"github.com/go-xmlfmt/xmlfmt"
)

type GetConfig struct {
	Name string
	Type string  	// running, startup, candidate...
	Filter string	// netconf filter
	Include map[string][]string
	Exclude map[string][]string
}

func (s *GetConfig) Task() app.TaskBase {
	return app.TaskBase{
		Name: s.Name,
		Include: s.Include,
		Exclude: s.Exclude,
	}
}

func (s *GetConfig) Run(h *app.Host, c *netconf.Driver, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name
	
	// ==== Custom
	if s.Type == "" {
		s.Type = "running"
	}
	
	r, err := c.GetConfig(s.Type, opoptions.WithFilter(s.Filter))
	if err != nil {
		return res, fmt.Errorf("failed to get config; error: %+v", err)
	}

	res["result"] = xmlfmt.FormatXML(r.Result, "", "  ")

	// === Required
	return res, nil

}