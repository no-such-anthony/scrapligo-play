package main


import (
	"fmt"
	"time"
	"sync"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
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

func getVersion(h Host) {
	p, err := platform.NewPlatform(
		h.Platform,
		h.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(h.Username),
		options.WithAuthPassword(h.Password),
		options.WithSSHConfigFile("../inventory/ssh_config"),
	)
	if err != nil {
		fmt.Printf("failed to create platform for %s: %+v\n", h.Hostname, err)
		return
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
        fmt.Printf("failed to fetch network driver for %s: %+v\n\n", h.Hostname, err)
        return
    }
	
	err = d.Open()
	if err != nil {
		fmt.Printf("failed to open driver for %s: %+v\n", h.Hostname, err)
		return
	}
	defer d.Close()

	rs, err := d.SendCommand("show version")
	if err != nil {
		fmt.Printf("failed to send command for %s: %+v\n", h.Hostname, err)
		return
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil {
		fmt.Printf("failed to parse command for %s: %+v\n", h.Hostname, err)
		return
	}

	fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
				h.Hostname, parsedOut[0]["HARDWARE"],
				parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])

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

func chunker(hosts map[string]Host, num_workers int) [][]string {
	
	keys := make([]string, 0, len(hosts))
	for k := range hosts {
		keys = append(keys, k)
	}

	size := num_workers
	total := len(keys)
	var chunks [][]string

	for i:=0;i<total;i+=size {
		end := i + size

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > total {
			end = total
		}

		chunks = append(chunks, keys[i:end])
	}
	return chunks
}


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := getHosts()
	//fmt.Println(hosts)

	var wg sync.WaitGroup

	//Waitgroup with chunks of data to restrict the number of goroutines.
	chunks := chunker(hosts, 5)

	for chunk := range chunks {

		for _, host_key := range chunks[chunk] {
			wg.Add(1)
		
			go func(h Host) {
				defer wg.Done()
				getVersion(h)
			}(hosts[host_key])
		
		}
		wg.Wait()
	}

}