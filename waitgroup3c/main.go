package main


import (
	"fmt"
	"time"
	"sync"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
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
	Result    []map[string]interface{}
}


type Hosts map[string]*Host

type Task struct {
	Name		string
	Function	func(*Host, map[string]interface{}) (map[string]interface{}, error)
	Args		map[string]interface{}
}

type Tasks []Task


func getVersion(h *Host, kwargs map[string]interface{}) (map[string]interface{}, error) {

	res := make(map[string]interface{})

	c := h.Connection

	fmt.Printf("%v - args: %+v\n",h.Name, kwargs)
	if len(h.Result)>=1 {
		fmt.Printf("%v - previous result: %+v\n",h.Name, h.Result[len(h.Result)-1])
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


func runTasks(h *Host, tasks Tasks) (error) {

	c, err := getConnection(*h)
	if err != nil {
		result := make(map[string]interface{})
		result["connection"] = err
		h.Result = append(h.Result, result)
		return err
	}
	h.Connection = c

	// task loop
	for _, task := range tasks {

		result := make(map[string]interface{})
		res, err := task.Function(h, task.Args)
		if err != nil {
			result[task.Name] = err
			h.Result = append(h.Result, result)
			return err
		}
		result[task.Name] = res
		h.Result = append(h.Result, result)
	}

	c.Close()
	return nil

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

	runner(hosts, tasks)


	fmt.Print("\n\n")
	fmt.Println("======================= RESULTS ========================================")
	for n, h := range hosts {
		fmt.Println("Name:", n)
		fmt.Println("Length:",len(h.Result))
		//fmt.Println(h.Result)
		for _, res := range h.Result {
			fmt.Println(res)
		}
		fmt.Print("\n\n")
	}

}


func runner(hosts Hosts, tasks Tasks) {

	var wg sync.WaitGroup

	num_workers := 7
	guard := make(chan bool, num_workers)
	wg.Add(len(hosts))

	//Combining Waitgroup with a channel to restrict number of goroutines.

	for _, host := range hosts {
	
		guard <- true
		go func(h *Host) {
			defer wg.Done()
			err := runTasks(h, tasks)
			//Print errors immediately but collate results for printing later.
			if err != nil {
				fmt.Println(err.Error())
			}
			<-guard
		}(host)
	}
	wg.Wait()
}


func getConnection(h Host) (*network.Driver, error) {

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


func getHosts() Hosts {

	devices := []string{"no.suchdomain","192.168.204.101","192.168.204.102","192.168.204.103","192.168.204.104"}

	hosts := make(Hosts)

	for _,value := range devices {
		var host Host
		host.Data = make(map[string]interface{})
		host.Result = make([]map[string]interface{},0)
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