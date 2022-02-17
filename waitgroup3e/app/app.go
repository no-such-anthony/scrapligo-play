package app

import (
	"fmt"
	"sync"
	"main/app/inventory"
	"main/app/connections"
	"main/app/tasks"
	//"github.com/scrapli/scrapligo/driver/network"
)


func runTasks(h *inventory.Host, t []tasks.RunTask, rc chan<- []map[string]interface{}) {

	host_results := []map[string]interface{}{}

	conn := connections.ScrapligoSsh{}

	err := connections.Connectors.Open(&conn, h)
	if err != nil {
		result := make(map[string]interface{})
		result["task"] = "connection"
		result["result"] = err
		result["failed"] = true
		host_results = append(host_results, result)
		rc <- host_results
		return
	}
	
	h.Connection = conn.C

	// task loop
	for _, task := range t {
		result := make(map[string]interface{})
		result["task"] = task.Named()
		res, err := task.Run(h, host_results)
		if err != nil {
			result["result"] = err
			result["failed"] = true
			host_results = append(host_results, result)
			rc <- host_results
			return
		}
		result["result"] = res
		host_results = append(host_results, result)
	}
	h.Connection.Close()
	rc <- host_results

}


func Runner(hosts inventory.Hosts, t []tasks.RunTask) (map[string]interface{})  {

	var wg sync.WaitGroup

	num_workers := 10
	guard := make(chan bool, num_workers)
	rc := make(chan []map[string]interface{}, num_workers)
	results := map[string]interface{}{}
	wg.Add(len(hosts))

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

			results[h.Name] = res
			<-guard
		}(host)
	}
	wg.Wait()
	return results
}
