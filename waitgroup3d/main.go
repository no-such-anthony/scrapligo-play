package main


import (
	"fmt"
	"time"
	"sync"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/driver/network"
)


type Host struct {
	Name      string
	Hostname  string
	Platform  string
	Username  string
	Password  string
	Enable    string
	StrictKey bool
	Connection *network.Driver
	Data      map[string]interface{}
}


type Hosts map[string]*Host

type Task struct {
	Name		string
	Function	func(*Host, map[string]interface{}, []map[string]interface{}) (map[string]interface{}, error)
	Args		map[string]interface{}
}

type Tasks []Task

// All tasks in Tasbook would need host,kwargs,result
func getVersion(h *Host, kwargs map[string]interface{}, prev_results []map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})
	c := h.Connection

	fmt.Printf("%v - args: %+v\n",h.Name, kwargs)
	if len(prev_results)>=1 {
		fmt.Printf("%v - previous result: %+v\n",h.Name, prev_results[len(prev_results)-1])
	}

	rs, err := c.SendCommand("show version")
	if err != nil {
		return res, fmt.Errorf("failed to send command for %s: %+v", h.Hostname, err)
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil {
		return res, fmt.Errorf("failed to parse command for %s: %+v", h.Hostname, err)
	}

	res["result"] = fmt.Sprintf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s",
				h.Hostname, parsedOut[0]["HARDWARE"],
				parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])

	return res, nil

}


func runTasks(h *Host, tasks Tasks) ([]map[string]interface{}, error) {

	host_results := []map[string]interface{}{}

	c, err := getConnection(*h)
	if err != nil {
		result := make(map[string]interface{})
		result["connection"] = err
		result["failed"] = true
		host_results = append(host_results, result)
		return host_results, err
	}
	h.Connection = c

	// task loop
	for _, task := range tasks {
		result := make(map[string]interface{})
		res, err := task.Function(h, task.Args, host_results)
		if err != nil {
			result[task.Name] = err
			result["failed"] = true
			host_results = append(host_results, result)
			return host_results, err
		}
		result[task.Name] = res
		host_results = append(host_results, result)
	}
	c.Close()
	return host_results, err


}


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := getHosts()
	//fmt.Println(hosts)

	//attempt at a simple playbook/runbook/taskbook in code
	var tasks Tasks
	var task1 Task
	var task2 Task

	task1.Name = "getVersion"
	task1.Function = getVersion
	task1.Args = make(map[string]interface{})
	
	task2.Name = "getVersion2"
	task2.Function = getVersion
	task2.Args = make(map[string]interface{})
	task2.Args["test"] = "my args test"

	tasks = append(tasks, task1)
	tasks = append(tasks, task2)
	//fmt.Printf("%+v\n", tasks)

	results := runner(hosts, tasks)

	fmt.Print("\n\n")
	fmt.Println("======================= RESULTS =================================")
	for n, h := range results {
		fmt.Println("Name:", n)
		for _, res := range h.([]map[string]interface{}) {
			fmt.Println(res)
		}
		fmt.Print("\n\n")
	}

}


func runner(hosts Hosts, tasks Tasks) (map[string]interface{})  {

	var wg sync.WaitGroup

	num_workers := 7
	guard := make(chan bool, num_workers)
	results := map[string]interface{}{}
	wg.Add(len(hosts))
	mutex := &sync.Mutex{}

	//Combining Waitgroup with a channel to restrict number of goroutines.
	//Results returned in a channel.
	for _, host := range hosts {
		guard <- true
		go func(h *Host) {
			defer wg.Done()
			res, err := runTasks(h, tasks)

			// Print errors immediately
			if err != nil {
				fmt.Println("error:", res)
			}

			mutex.Lock()
			results[h.Name] = res
			mutex.Unlock()

			<-guard
		}(host)
	}
	wg.Wait()
	return results
}


func getConnection(h Host) (*network.Driver, error) {

	p, err := platform.NewPlatform(
		h.Platform,
		h.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(h.Username),
		options.WithAuthPassword(h.Password),
		options.WithSSHConfigFile("../inventory/ssh_config"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform for %s: %+v", h.Hostname, err)
	}

	c, err := p.GetNetworkDriver()
	if err != nil {
        return nil, fmt.Errorf("failed to fetch network driver for %s: %+v\n\n", h.Hostname, err)
    }

	err = c.Open()
	if err != nil {
		return c, fmt.Errorf("failed to open driver for %s: %+v", h.Hostname, err)
	}

	return c, nil 

}


func getHosts() Hosts {

	devices := []string{"no.suchdomain","192.168.204.101","192.168.204.102","192.168.204.103","192.168.204.104"}

	hosts := make(Hosts)

	for _,value := range devices {
		var host Host
		host.Data = make(map[string]interface{})
		host.Name = value
		host.Hostname = value
		host.Platform = "cisco_iosxe"
		host.Username = "fred"
		host.Password = "bedrock"
		host.Data["example_only"] = 100

		hosts[host.Name] = &host
		
	}

	return hosts
}


func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}