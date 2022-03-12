package sshgomiko

import (
	"main/play/app"
	"github.com/Ali-aqrabawi/gomiko/pkg/types"
)

type Tasker interface {
	Run(*app.Host, types.Device, []map[string]interface{}) (map[string]interface{}, error)
	Task() app.TaskBase
}

type Wrap struct {
	Tasker Tasker
}

func (r *Wrap) Run(h *app.Host, prev_res []map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	task := r.Tasker.Task()

	if app.Skip(h, task.Include, task.Exclude) {
		res["task"] = task.Name
		res["skipped"] = true
		return res, nil
	}

	conn, err := GetConn(h)
	if err != nil {
		res["task"] = task.Name
		res["result"] = err
		res["failed"] = true
		return res, &app.ConnectionError{h.Name, err}
	}

	c := conn.(*GomikoSsh).C
	res, err = r.Tasker.Run(h, c, prev_res)
	if err != nil {
		return res, &app.TaskError{task.Name, h.Name, err}
	}

	return res, nil

}