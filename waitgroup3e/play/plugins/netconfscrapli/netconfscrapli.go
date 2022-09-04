package netconfscrapli

import (
	"fmt"
	"main/play/app"
	"github.com/scrapli/scrapligo/driver/netconf"
)

type Task interface {
	Run(*app.Host, *netconf.Driver, []map[string]interface{}) (map[string]interface{}, error)
	Info() app.TaskBase
}

type Wrap struct {
	Task
}

func (r *Wrap) Start(h *app.Host, prev_res []map[string]interface{}) (res map[string]interface{}, wrapErr error) {

	res = make(map[string]interface{})
	task := r.Task.Info()
	res["task"] = task.Name
	wrapErr = nil

	if app.Skip(h, task.Include, task.Exclude) {
		res["skipped"] = true
		return res, nil
	}

	// try to handle someone using panic instead of returning error from wrapped functions
	defer func() {
		if err := recover(); err != nil {
			res["result"] = err
			res["failed"] = true
	
			// find out what the panic was and set wrapErr
			switch x := err.(type) {
			case string:
				wrapErr = &app.TaskError{task.Name, h.Name, fmt.Errorf(x)}
			case error:
				wrapErr = &app.TaskError{task.Name, h.Name, x}
			default:
				wrapErr = &app.TaskError{task.Name, h.Name, fmt.Errorf("unknown panic")}
			}
		}
	}()

	conn, err := GetConn(h)
	if err != nil {
		res["result"] = err
		res["failed"] = true
		return res, &app.TaskError{task.Name, h.Name, err}
	}

	c := conn.(*ScrapligoNetconf).C
	res, err = r.Task.Run(h, c, prev_res)
	if err != nil {
		return res, &app.TaskError{task.Name, h.Name, err}
	}

	return res, nil

}