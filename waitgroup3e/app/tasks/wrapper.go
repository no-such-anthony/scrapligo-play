package tasks

import (
	"main/app/inventory"
)

type Tasker interface {
	Run(*inventory.Host, []map[string]interface{}) (map[string]interface{}, error)
	Task() TaskBase
}

type Wrap struct {
	Tasker Tasker
}

func (r *Wrap) Run(h *inventory.Host, prev_res []map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	task := r.Tasker.Task()

	if inventory.Skip(h, task.Include, task.Exclude) {
		res["task"] = task.Name
		res["skipped"] = true
		return res, nil
	}

	res, err := r.Tasker.Run(h, prev_res)

	return res, err

}