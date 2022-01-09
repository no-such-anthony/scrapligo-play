package main

// based on code found at https://github.com/PacktPublishing/Network-Automation-with-Go

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

type Inventory struct {
	Hosts map[string]Host
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}

func getVersion(h Host) {
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
		fmt.Printf("failed to create driver for %s: %+v\n", h.Hostname, err)
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

func getInventory() Inventory {

	var inventory Inventory
	devices := []string{"no.suchdomain","192.168.204.101","192.168.204.102","192.168.204.103","192.168.204.104"}

	inventory.Hosts = make(map[string]Host)

	for _,value := range devices {
		var host Host
		host.Data = make(map[string]interface{})

		host.Name = value
		host.Hostname = value
		host.Platform = "cisco_iosxe"
		host.Username = "fred"
		host.Password = "bedrock"
		host.Data["example_only"] = 100

		inventory.Hosts[host.Name] = host
		
	}

	return inventory
}

func chunker(inventory Inventory, num_workers int) [][]string {

	keys := make([]string, 0, len(inventory.Hosts))
    for k := range inventory.Hosts {
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

	inventory := getInventory()
	//fmt.Println(inventory)

	var wg sync.WaitGroup

	//Note: Not ideal, but here we are using chunks to restrict the number of goroutines
	chunks := chunker(inventory, 5)

	for chunk := range chunks {

		for _, host_key := range chunks[chunk] {
			wg.Add(1)
		
			go func(h Host) {
				defer wg.Done()
				getVersion(h)
			}(inventory.Hosts[host_key])
		
		}
		wg.Wait()
	}

}