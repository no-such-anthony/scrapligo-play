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
	Data      map[string]interface{}
	Result   map[string]interface{}
	Connection *network.Driver
}

type Hosts map[string]*Host


func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}

func getVersion(h *Host) {

	result := make(map[string]interface{})
	result["name"] = h.Name

	c, err := getConnection(*h)
	if err != nil {
		result["error"] = err.Error()
		fmt.Println(err)
		h.Result = result
		return
	}

	rs, err := c.SendCommand("show version")
	if err != nil {
		err := fmt.Errorf("failed to send command for %s: %+v", h.Hostname, err)
		result["error"] = err.Error()
		fmt.Println(err)
		h.Result = result
		return
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil {
		err := fmt.Errorf("failed to parse command for %s: %+v", h.Hostname, err)
		result["error"] = err.Error()
		fmt.Println(err)
		h.Result = result
		return
	}

	result["result"] = parsedOut[0]
	h.Result = result
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


func worker(host_jobs <-chan *Host, wg *sync.WaitGroup) {

	for h := range host_jobs {
		getVersion(h)
		wg.Done()
	}
}

func runner(hosts Hosts) {
	//waitgroup/channel workerpool combo, storing results in the host pointer
	var wg sync.WaitGroup
	const num_workers = 1
	host_jobs := make(chan *Host, len(hosts))
	wg.Add(len(hosts))
	
	//worker pools
	for w := 1; w <= num_workers; w++ {
		go worker(host_jobs, &wg)
	}

	for _, host := range hosts {
		host_jobs <- host
	}
	close(host_jobs)
	wg.Wait()
}


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := getHosts()
	//fmt.Println(hosts)

	runner(hosts)

	for name, host := range hosts {
		fmt.Println(name)
		if err, ok := host.Result["error"]; ok {
			fmt.Println(err)
		} else {
			fmt.Println(host.Result["result"])
		}
	}

}


func getHosts() Hosts {

	devices := []string{"no.suchdomain","192.168.204.101","192.168.204.102","192.168.204.103","192.168.204.104"}

	hosts := make(Hosts)

	for _,value := range devices {
		var host Host
		host.Data = make(map[string]interface{})
		host.Result = make(map[string]interface{})
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