package sshscrapli

import (
	"main/app/inventory"
	//"main/app/connections"
	"main/app/tasks"
	"github.com/scrapli/scrapligo/driver/network"
)

type ScrapliSSHTasker interface {
	Run(*inventory.Host, *network.Driver, []map[string]interface{}) (map[string]interface{}, error)
	Task() tasks.TaskBase
}

type ScrapliSSHWrap struct {
	Tasker ScrapliSSHTasker
}

func (r *ScrapliSSHWrap) Run(h *inventory.Host, prev_res []map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	task := r.Tasker.Task()

	if inventory.Skip(h, task.Include, task.Exclude) {
		res["task"] = task.Name
		res["skipped"] = true
		return res, nil
	}

	conn, err := GetConn(h)
	if err != nil {
		res["task"] = task.Name
		res["result"] = err
		res["failed"] = true
		return res, err	
	}

	c := conn.(*ScrapligoSsh).C
	res, err = r.Tasker.Run(h, c, prev_res)

	return res, err

}