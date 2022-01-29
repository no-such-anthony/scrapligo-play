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
}


type Hosts map[string]Host


func getVersion(h Host, c *network.Driver) (string, error) {


	rs, err := c.SendCommand("show version")
	if err != nil {
		return "", fmt.Errorf("failed to send command for %s: %+v", h.Hostname, err)
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil {
		return "", fmt.Errorf("failed to parse command for %s: %+v", h.Hostname, err)
	}

	res := fmt.Sprintf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s",
				h.Hostname, parsedOut[0]["HARDWARE"],
				parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])

	return res, nil

}


func runTasks(h Host) (string, error) {

	res := ""
	c, err := getConnection(h)
	if err != nil {
		return err.Error(), err
	}
	// put all your tasks here
	res, err = getVersion(h, c)
	if err != nil {
		return res, err
	}

	c.Close()
	return res, nil

}


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := getHosts()
	//fmt.Println(hosts)

	results := runner(hosts)

	fmt.Print("\n\n")
	for host, res := range results {
		fmt.Println("Name:", host)
		fmt.Println(res)
		fmt.Print("\n\n")
	}

}


func runner(hosts Hosts) map[string]string {

	var wg sync.WaitGroup

	num_workers := 5
	guard := make(chan bool, num_workers)
	results := make(map[string]string)
	wg.Add(len(hosts))
	mux := &sync.Mutex{}
	//Combining Waitgroup with a channel to restrict number of goroutines.

	for _, host := range hosts {
	
		guard <- true
		go func(h Host) {
			defer wg.Done()
			res, err := runTasks(h)
			//Print errors immediately but collate results for printing later.
			if err != nil {
				fmt.Println(err.Error())
			}
			mux.Lock()
			results[h.Name] = res
			mux.Unlock()
			<-guard
		}(host)
    
	}
	wg.Wait()

	return results
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

		host.Name = value
		host.Hostname = value
		host.Platform = "cisco_iosxe"
		host.Username = "fred"
		host.Password = "bedrock"
		host.Data["example_only"] = 100

		hosts[host.Name] = host
		
	}

	return hosts
}


func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}