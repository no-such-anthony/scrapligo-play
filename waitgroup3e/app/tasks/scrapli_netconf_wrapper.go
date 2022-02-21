package tasks

import (
	"main/app/inventory"
	"main/app/connections"
	"github.com/scrapli/scrapligo/netconf"
)

type ScrapliNetconfTasker interface {
	Run(*inventory.Host, *netconf.Driver, []map[string]interface{}) (map[string]interface{}, error)
	Task() TaskBase
}

type ScrapliNetconfWrap struct {
	Tasker ScrapliNetconfTasker
}

func (r *ScrapliNetconfWrap) Run(h *inventory.Host, prev_res []map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	task := r.Tasker.Task()

	if inventory.Skip(h, task.Include, task.Exclude) {
		res["task"] = task.Name
		res["skipped"] = true
		return res, nil
	}

	conn, err := connections.GetConn(h, "scrapli_netconf")
	if err != nil {
		res["task"] = task.Name
		res["result"] = err
		res["failed"] = true
		return res, err	
	}

	c := conn.(*connections.ScrapligoNetconf).C
	res, err = r.Tasker.Run(h, c, prev_res)

	return res, err

}