package tasks

import (
	"fmt"
	"main/app/inventory"
)

type TaskTest struct {
	Name string
	Kwargs map[string]interface{}
	Include map[string][]string
	Exclude map[string][]string
}

func (s *TaskTest) Task() TaskBase {
	return TaskBase{
		Name: s.Name,
		Include: s.Include,
		Exclude: s.Exclude,
	}
}

func (s *TaskTest) Run(h *inventory.Host, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name
	
	// ==== Custom
	fmt.Printf("%v - args: %+v\n",h.Name, s.Kwargs)
	if len(prev_results)>=1 {
			fmt.Printf("%v - previous result: %+v\n",h.Name, prev_results[len(prev_results)-1])
	}

	res["result"] = fmt.Sprintf("host %s, just chillin'", h.Name)

	// === Required
	return res, nil

}
