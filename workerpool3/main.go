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
	Data      map[string]interface{}
	Result   map[string]interface{}
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
		fmt.Println(result["error"])
		h.Result = result
		return
	}

	rs, err := c.SendCommand("show version")
	if err != nil {
		result["error"] = fmt.Sprintf("failed to send command for %s: %+v", h.Hostname, err)
		fmt.Println(result["error"])
		h.Result = result
		return
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil {
		result["error"] = fmt.Sprintf("failed to parse command for %s: %+v", h.Hostname, err)
		fmt.Println(result["error"])
		h.Result = result
		return
	}

	result["result"] = parsedOut[0]
	h.Result = result
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
		return nil, fmt.Errorf("failed to create platform for %s: %+v\n\n", h.Hostname, err)
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


func worker(host_jobs <-chan *Host, wg *sync.WaitGroup) {

	for h := range host_jobs {
		getVersion(h)
		wg.Done()
	}
}

func runner(hosts Hosts) {
	//waitgroup/channel workerpool combo, storing results in the host pointer
	var wg sync.WaitGroup
	const num_workers = 7
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
			continue
		}
		fmt.Println(host.Result["result"])
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