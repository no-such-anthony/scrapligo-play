package tasks

import (
	"main/app/inventory"
)

type TaskBase struct {
	Name string
	Include map[string][]string
	Exclude map[string][]string
}

type Wrapper interface {
	Run(*inventory.Host, []map[string]interface{}) (map[string]interface{}, error)
}

type ConnectionError struct {
	Name string
	Err error
}

func (e *ConnectionError) Error() string {
	return e.Name + ": " + e.Err.Error()
}

type TaskError struct {
	Task string
	Name string
	Err error
}

func (e *TaskError) Error() string {
	return e.Name + ": " + e.Task + ": " + e.Err.Error()
}
