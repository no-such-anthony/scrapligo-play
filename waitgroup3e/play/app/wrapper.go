package app

import (
	"fmt"
)

type Tasker interface {
	Run(*Host, []map[string]interface{}) (map[string]interface{}, error)
	Task() TaskBase
}

type Wrap struct {
	Tasker Tasker
}

func (r *Wrap) Run(h *Host, prev_res []map[string]interface{}) (res map[string]interface{}, wrapErr error) {

	res = make(map[string]interface{})
	task := r.Tasker.Task()
	res["task"] = task.Name
	wrapErr = nil

	if Skip(h, task.Include, task.Exclude) {
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
				wrapErr = &TaskError{task.Name, h.Name, fmt.Errorf(x)}
			case error:
				wrapErr = &TaskError{task.Name, h.Name, x}
			default:
				wrapErr = &TaskError{task.Name, h.Name, fmt.Errorf("unknown panic")}
			}
		}
	}()

	res, err := r.Tasker.Run(h, prev_res)
	if err != nil {
		return res, &TaskError{task.Name, h.Name, err}
	}

	return res, nil

}