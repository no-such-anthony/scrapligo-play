package sshgoexpect

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"main/play/app"
)

type Tasker interface {
	Run(*app.Host, *ssh.Client, []map[string]interface{}) (map[string]interface{}, error)
	Task() app.TaskBase
}

type Wrap struct {
	Tasker Tasker
}

func (r *Wrap) Run(h *app.Host, prev_res []map[string]interface{}) (res map[string]interface{}, wrapErr error) {

	res = make(map[string]interface{})
	task := r.Tasker.Task()
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

	c := conn.(*GoExpectSsh).C
	res, err = r.Tasker.Run(h, c, prev_res)
	if err != nil {
		return res, &app.TaskError{task.Name, h.Name, err}
	}

	return res, nil

}