package app

import (
	"fmt"
	"sync"
	"main/app/inventory"
	"main/app/tasks"
)


func runTasks(h *inventory.Host, t []tasks.Wrapper) []map[string]interface{} {

	host_results := []map[string]interface{}{}

	// task loop
	for _, task := range t {

		res, err := task.Run(h, host_results)
		// don't continue on error
		if err != nil {

			switch err.(type) {
			case *tasks.ConnectionError:
				fmt.Println("connection error:", err)
			case *tasks.TaskError:
				fmt.Println("task error:", err)
				for _, v := range h.Connections {
					v.Close()
				}
			default:
				fmt.Println("unexpected error:", err)
			}
			host_results = append(host_results, res)
			return host_results
		}

		host_results = append(host_results, res)
	}

	for _, v := range h.Connections {
		v.Close()
	}

	return host_results

}


func Runner(hosts inventory.Hosts, t []tasks.Wrapper) (map[string]interface{})  {

	var wg sync.WaitGroup

	num_workers := 20
	guard := make(chan bool, num_workers)
	results := map[string]interface{}{}
	wg.Add(len(hosts))
	mutex := &sync.Mutex{}

	//Combining Waitgroup with a channel to restrict number of goroutines.
	for _, host := range hosts {
		guard <- true
		go func(h *inventory.Host) {
			defer wg.Done()
			res := runTasks(h, t)
			mutex.Lock()
			results[h.Name] = res
			mutex.Unlock()
			<-guard
		}(host)
	}
	wg.Wait()
	return results
}
