package main


import (
	"fmt"
	"time"
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

func getVersion(h Host) map[string]interface{} {

	result := make(map[string]interface{})
	result["name"] = h.Name

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
		msg := fmt.Sprintf("failed to create driver for %s: %+v", h.Hostname, err)
		result["error"] = msg
		return result
	}

	err = d.Open()
	if err != nil {
		msg := fmt.Sprintf("failed to open driver for %s: %+v", h.Hostname, err)
		result["error"] = msg
		return result
	}
	defer d.Close()

	rs, err := d.SendCommand("show version")
	if err != nil {
		msg := fmt.Sprintf("failed to send command for %s: %+v", h.Hostname, err)
		result["error"] = msg
		return result
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil {
		msg := fmt.Sprintf("failed to parse command for %s: %+v", h.Hostname, err)
		result["error"] = msg
		return result
	}

	//fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
	//			h.Hostname, parsedOut[0]["HARDWARE"],
	//			parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])

	result["result"] = parsedOut[0]

	return result
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

	//In/Out buffered channels with a results returned channel and num_workers.
	const num_workers = 5
	host_jobs := make(chan Host, 100)
	host_results := make(chan map[string]interface{}, 100)
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
		results := <-host_results
		//fmt.Println(results)
		if err, ok := results["error"]; ok {
			fmt.Printf("Host: %s had error %s\n\n", results["name"], err)
		} else {
			result := results["result"].(map[string]interface{})
			fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
					result["HOSTNAME"], result["HARDWARE"],
					result["VERSION"], result["UPTIME"])
		}
		agg_results[results["name"].(string)] = results
	}
	fmt.Println("\n\n")
	fmt.Println("And again, as we stored the results such that we can use outside of the return channel loop.\n")
	for name, results := range agg_results {
		//fmt.Println(name, results)
		result := results.(map[string]interface{})
		if err, ok := result["error"]; ok {
			fmt.Printf("Host: %s had error %s\n\n", name, err)
		} else {
			result = result["result"].(map[string]interface{})
			fmt.Printf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s\n\n",
					result["HOSTNAME"], result["HARDWARE"],
					result["VERSION"], result["UPTIME"])
		}
	}
	
}