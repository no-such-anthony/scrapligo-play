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

