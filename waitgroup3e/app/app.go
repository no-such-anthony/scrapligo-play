package app

import (
	"fmt"
	"sync"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
	"main/app/inventory"

)


type RunTask interface {
	Run(*inventory.Host, []map[string]interface{}) (map[string]interface{}, error)
	Named() string
}


func runTasks(h *inventory.Host, t []RunTask, rc chan<- []map[string]interface{}) {

	host_results := []map[string]interface{}{}

	c, err := getConnection(*h)
	if err != nil {
		result := make(map[string]interface{})
		result["name"] = "connection"
		result["result"] = err
		result["failed"] = true
		host_results = append(host_results, result)
		rc <- host_results
		return
	}
	h.Connection = c

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
	c.Close()
	rc <- host_results

}


func Runner(hosts inventory.Hosts, t []RunTask) (map[string]interface{})  {

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


func getConnection(h inventory.Host) (*network.Driver, error) {

	c, err := core.NewCoreDriver(
		h.Hostname,
		h.Platform,
		base.WithAuthStrictKey(h.StrictKey),
		base.WithAuthUsername(h.Username),
		base.WithAuthPassword(h.Password),
		//base.WithAuthSecondary(h.Enable),
		//base.WithTransportType("standard"),
		//base.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		return c, fmt.Errorf("failed to create driver for %s: %+v", h.Hostname, err)
	}

	err = c.Open()
	if err != nil {
		return c, fmt.Errorf("failed to open driver for %s: %+v", h.Hostname, err)
	}

	return c, nil 

}
