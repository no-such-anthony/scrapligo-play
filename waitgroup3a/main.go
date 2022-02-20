package main


import (
	"fmt"
	"time"
	"sync"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
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

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}

func getVersion(h Host) (string, error) {
	d, err := core.NewCoreDriver(
		h.Hostname,
		h.Platform,
		base.WithAuthStrictKey(h.StrictKey),
		base.WithAuthUsername(h.Username),
		base.WithAuthPassword(h.Password),
		//base.WithTransportType("standard"),
		//base.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create driver for %s: %+v", h.Hostname, err)
	}

	err = d.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open driver for %s: %+v", h.Hostname, err)
	}
	defer d.Close()

	rs, err := d.SendCommand("show version")
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


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := getHosts()
	//fmt.Println(hosts)

	var wg sync.WaitGroup

	num_workers := 2
	guard := make(chan bool, num_workers)
	results := make(map[string]string)
	mux := &sync.Mutex{}

	//Combining Waitgroup with a channel to restrict number of goroutines.
	wg.Add(len(hosts))
	for _, host := range hosts {
	
		guard <- true
		go func(h Host) {
			defer wg.Done()
			res, err := getVersion(h)
			//Print errors immediately but collate results for printing later.
			if err != nil {
				fmt.Println(err.Error())
				res = err.Error()
			}
			mux.Lock()
			results[h.Name] = res
			mux.Unlock()
			<-guard
		}(host)
    
	}
	wg.Wait()

	fmt.Print("\n\n")
	for host, res := range results {
		fmt.Println("Name:", host)
		fmt.Println(res)
		fmt.Print("\n\n")
	}

}