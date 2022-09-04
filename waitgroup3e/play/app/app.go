package app

import (
	"fmt"
	"sync"
	"strings"
)


func runTasks(h *Host, t []Play) []map[string]interface{} {

	host_results := []map[string]interface{}{}

	// task loop
	for _, task := range t {

		res, err := task.Start(h, host_results)
		// don't continue on error
		if err != nil {
			fmt.Println("error:", strings.TrimSuffix(err.Error(), "\n"))
			for _, v := range h.Connections {
				v.Close()
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


func Runner(hosts Hosts, t []Play) (map[string]interface{})  {

	var wg sync.WaitGroup

	num_workers := 20
	guard := make(chan bool, num_workers)
	results := map[string]interface{}{}
	wg.Add(len(hosts))
	mutex := &sync.Mutex{}

	//Combining Waitgroup with a channel to restrict number of goroutines.
	for _, host := range hosts {
		guard <- true
		go func(h *Host) {
			defer wg.Done()
			res := runTasks(h, t)
			fmt.Println("runner: " + h.Name + " completed tasks.")
			mutex.Lock()
			results[h.Name] = res
			mutex.Unlock()
			<-guard
		}(host)
	}
	wg.Wait()
	return results
}
