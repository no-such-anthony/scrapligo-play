package main


import (
	"fmt"
	"time"
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


func getVersion(h Host) map[string]interface{} {

	result := make(map[string]interface{})
	result["name"] = h.Name

	p, err := platform.NewPlatform(
		h.Platform,
		h.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(h.Username),
		options.WithAuthPassword(h.Password),
		options.WithSSHConfigFile("../inventory/ssh_config"),
	)
	if err != nil {
		result["error"] = fmt.Sprintf("failed to create platform for %s: %+v\n\n", h.Hostname, err)
		return result
	}

	d, err := p.GetNetworkDriver()
	if err != nil {
        result["error"] = fmt.Sprintf("failed to fetch network driver for %s: %+v\n\n", h.Hostname, err)
        return result
    }

	err = d.Open()
	if err != nil {
		result["error"] = fmt.Sprintf("failed to open driver for %s: %+v", h.Hostname, err)
		return result
	}
	defer d.Close()

	rs, err := d.SendCommand("show version")
	if err != nil {
		result["error"] = fmt.Sprintf("failed to send command for %s: %+v", h.Hostname, err)
		return result
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil {
		result["error"] = fmt.Sprintf("failed to parse command for %s: %+v", h.Hostname, err)
		return result
	}

	result["result"] = fmt.Sprintf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s",
							parsedOut[0]["HOSTNAME"], parsedOut[0]["HARDWARE"],
							parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])
	return result

}

func worker(host_jobs <-chan Host, host_results chan<- map[string]interface{}) {
	
	for h := range host_jobs {
		result := getVersion(h)
		host_results <- result
	}

}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := getHosts()
	//fmt.Println(hosts)

	// In/Out buffered channels with a results returned channel and num_workers.
	const num_workers = 5
	host_jobs := make(chan Host, len(hosts))	// room to drop all hosts into the channel at once.
	host_results := make(chan map[string]interface{}, len(hosts)) // make sure enough buffer or could end up with deadlock.
	agg_results := make(map[string]interface{})

	//worker pools
	for w := 1; w <= num_workers; w++ {
		go worker(host_jobs, host_results)
	}

	for _, host := range hosts {
		host_jobs <- host
	}
	close(host_jobs)

	fmt.Println("Printing worker results as they arrive...\n")
	for r := 1; r <= len(hosts); r++ {
		result := <-host_results
		agg_results[result["name"].(string)] = result
		if err, ok := result["error"]; ok {
			fmt.Printf("Host: %s had error %s\n\n", result["name"], err)
			continue
		}
		fmt.Printf("Host: %s results =>\n%s\n\n", result["name"], result["result"])
	}
	fmt.Println("\n\n")
	fmt.Println("And again, as we stored the results such that we can use outside of the return channel loop.\n")
	for name, results := range agg_results {
		result := results.(map[string]interface{})
		if err, ok := result["error"]; ok {
			fmt.Printf("Host: %s had error %s\n\n", name, err)
			continue
		}
		fmt.Printf("Host: %s results =>\n%s\n\n", name, result["result"])
	}
	
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