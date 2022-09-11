package app


type TaskBase struct {
	Name string
	Include map[string][]string
	Exclude map[string][]string
}

func (s *TaskBase) Info() TaskBase {
	return TaskBase{
		Name: s.Name,
		Include: s.Include,
		Exclude: s.Exclude,
	}
}


type Play interface {
	Start(*Host, []map[string]interface{}) (map[string]interface{}, error)
}

type TaskError struct {
	Task string
	Name string
	Err error
}

func (e *TaskError) Error() string {
	return e.Name + ": " + e.Task + ": " + e.Err.Error()
}
