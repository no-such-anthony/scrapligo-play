package app


type Tasker interface {
	Run(*Host, []map[string]interface{}) (map[string]interface{}, error)
	Task() TaskBase
}

type Wrap struct {
	Tasker Tasker
}

func (r *Wrap) Run(h *Host, prev_res []map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	task := r.Tasker.Task()

	if Skip(h, task.Include, task.Exclude) {
		res["task"] = task.Name
		res["skipped"] = true
		return res, nil
	}

	res, err := r.Tasker.Run(h, prev_res)
	if err != nil {
		return res, &TaskError{task.Name, h.Name, err}
	}

	return res, nil

}