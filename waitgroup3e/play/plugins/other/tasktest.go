package other

import (
	"fmt"
	"main/play/app"
)

type TaskTest struct {
	app.TaskBase
	Kwargs map[string]interface{}
}

func (s *TaskTest) Info() app.TaskBase {
	return s.TaskBase
}

func (s *TaskTest) Run(h *app.Host, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	// === Required
	res := make(map[string]interface{})
	res["task"] = s.Name
	
	// ==== Custom
	//fmt.Printf("%v - args: %+v\n",h.Name, s.Kwargs)
	//if len(prev_results)>=1 {
	//		fmt.Printf("%v - previous result: %+v\n",h.Name, prev_results[len(prev_results)-1])
	//}

	res["result"] = fmt.Sprintf("host %s, just chillin'", h.Name)

	// === Required
	return res, nil

}
