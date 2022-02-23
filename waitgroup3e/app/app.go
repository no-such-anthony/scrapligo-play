package app

import (
	"fmt"
	"sync"
	"main/app/inventory"
	"main/app/tasks"
)


func runTasks(h *inventory.Host, t []tasks.Wrapper) ([]map[string]interface{}, error) {

	host_results := []map[string]interface{}{}

	// task loop
	for _, task := range t {

		res, err := task.Run(h, host_results)
		// don't continue on error
		if err != nil {
			host_results = append(host_results, res)
			return host_results, err
		}

		host_results = append(host_results, res)
	}

	for _, v := range h.Connections {
		v.Close()
	}

	return host_results, nil

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
			res, err := runTasks(h, t)
			// Print errors immediately
			if err != nil {
				fmt.Println("error:", res)
			}
			//fmt.Println(res)
			mutex.Lock()
			results[h.Name] = res
			mutex.Unlock()
			<-guard
		}(host)
	}
	wg.Wait()
	return results
}
