package tasks

import (
	"main/app/inventory"
)

type Tasker interface {
	Run(*inventory.Host, []map[string]interface{}) (map[string]interface{}, error)
	Named() string
}