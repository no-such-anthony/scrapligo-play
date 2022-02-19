package app

import (
	"fmt"
	"sync"
	"main/app/inventory"
	"main/app/connections"
	"main/app/tasks"
	//"github.com/scrapli/scrapligo/driver/network"
)


func runTasks(h *inventory.Host, t []tasks.Tasker, rc chan<- []map[string]interface{}) {

	host_results := []map[string]interface{}{}

	var cc inventory.Connector

	switch h.Method {
	case "scrapli_ssh": 
		cc = inventory.Connector(&connections.ScrapligoSsh{})
	case "scrapli_netconf": 
		cc = inventory.Connector(&connections.ScrapligoNetconf{})
	default:
		cc = inventory.Connector(&connections.ScrapligoSsh{})
	}

	err := cc.Open(h)
	if err != nil {
		result := make(map[string]interface{})
		result["task"] = "connection"
		result["result"] = err
		result["failed"] = true
		host_results = append(host_results, result)
		rc <- host_results
		return
	}
	h.Connection = cc

	// task loop
	for _, task := range t {
		result := make(map[string]interface{})
		res, err := task.Run(h, host_results)
		if err != nil {
			result["result"] = err
			result["failed"] = true
			result["task"] = task.Named()
			host_results = append(host_results, result)
			rc <- host_results
			return
		}
		host_results = append(host_results, res)
	}
	h.Connection.Close()
	rc <- host_results

}


func Runner(hosts inventory.Hosts, t []tasks.Tasker) (map[string]interface{})  {

	var wg sync.WaitGroup

	num_workers := 10
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
