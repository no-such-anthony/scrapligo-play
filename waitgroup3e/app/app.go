package app

import (
	"fmt"
	"sync"
	"main/app/inventory"
	"main/app/tasks"
)


func runTasks(h *inventory.Host, t []tasks.Wrapper, rc chan<- []map[string]interface{}) {

	host_results := []map[string]interface{}{}

	// task loop
	for _, task := range t {

		res, err := task.Run(h, host_results)
		// don't continue on error
		if err != nil {
			host_results = append(host_results, res)
			rc <- host_results
			return
		}

		host_results = append(host_results, res)
	}

	for _, v := range h.Connections {
		v.Close()
	}
	rc <- host_results

}


func Runner(hosts inventory.Hosts, t []tasks.Wrapper) (map[string]interface{})  {

	var wg sync.WaitGroup

	num_workers := 20
	guard := make(chan bool, num_workers)
	rc := make(chan []map[string]interface{}, num_workers)
	results := map[string]interface{}{}
	wg.Add(len(hosts))
	mutex := &sync.Mutex{}

	//Combining Waitgroup with a channel to restrict number of goroutines.
	//Results returned in a channel.
	for _, host := range hosts {
		guard <- true
		go func(h *inventory.Host) {
			defer wg.Done()
			runTasks(h, t, rc)
			res := <-rc

			// Print errors immediately
			if _, ok := res[len(res)-1]["failed"]; ok {
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
