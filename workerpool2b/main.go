package main


import (
	"fmt"
	"time"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
	"errors"
	"regexp"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"strings"
	"os"
	"strconv"
)

type Host struct {
	Name      string
	Hostname  string
	Username  string
	Password  string
	Enable    string
	StrictKey bool
	Platform  string
	Groups    []string
	Data      map[string]interface{}
}

type Hosts map[string]*Host


func getVersion(h Host, conn *network.Driver) map[string]interface{} {

	result := make(map[string]interface{})
	result["name"] = h.Name

	c := conn

	rs, err := c.SendCommand("show version")
	if err != nil {
		result["error"] = fmt.Sprintf("failed to send command for %s: %+v", h.Hostname, err)
		return result
	}

	parsedOut, err := rs.TextFsmParse("../textfsm_templates/" + h.Platform + "_show_version.textfsm")
	if err != nil {
		result["error"] = fmt.Sprintf("failed to parse command for %s: %+v", h.Hostname, err)
		return result
	}

	//backup running config

	rs, err = c.SendCommand("show run")
	if err != nil {
		result["error"] = fmt.Sprintf("failed to send command for %s: %+v", h.Hostname, err)
		return result
	}

	if err = ioutil.WriteFile("../../" + h.Name + ".txt", []byte(rs.Result),0777); err != nil {
		result["error"] = fmt.Sprintf("failed to write config to disk for %s: %+v", h.Hostname, err)
		return result
	}

	// update host data if we want
	h.Data["SW version"] = parsedOut[0]["VERSION"]

	result["result"] = fmt.Sprintf("Hostname: %s\nHardware: %s\nSW Version: %s\nUptime: %s",
							parsedOut[0]["HOSTNAME"], parsedOut[0]["HARDWARE"],
							parsedOut[0]["VERSION"], parsedOut[0]["UPTIME"])
	return result
}



func worker(host_jobs <-chan Host, host_results chan<- map[string]interface{}) {

	for h := range host_jobs {
		conn, err := getConnection(h)
		if err != nil {
			result := make(map[string]interface{})
			result["name"] = h.Name
			result["error"] = err.Error()
			host_results <- result
			continue
		}
		// put your tasks here
		result := getVersion(h, conn)
		host_results <- result
		conn.Close()
		
	}
}

func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := getHostsbyYAML()

	// Filtering methods via regex slice, F_include and F_exclude, examples
	//hosts = F_include(hosts, "name", []string{"r1","r2"})
	//hosts = F_include(hosts, "name", []string{"r1|r2"})
	//hosts = F_include(hosts, "hostname", []string{"192.168"})
	//hosts = F_include(hosts, "platform", []string{"group2"})
	//hosts = F_include(hosts, "groups", []string{"group2","group3"})
	//hosts = F_include(hosts, "model", []string{"C35|Dumm"})  // <- from data map
	//hosts = F_include(hosts, "a_number", []string{"5"})  // <- from data map
	//hosts = F_exclude(hosts, "a_role", []string{"stuff2"})  // <- from data map
	//fmt.Println(hosts)

	// If no hosts then exit
	if len(hosts) == 0 {
		fmt.Println("No hosts.  Nothing to do.")
		os.Exit(0)
	}

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
		host_jobs <- *host
	}
	close(host_jobs)

	fmt.Println("Printing worker results as they arrive...\n")
	for r := 1; r <= len(hosts); r++ {
		result := <-host_results
		if err, ok := result["error"]; ok {
			fmt.Printf("Host: %s had error %s\n\n", result["name"], err)
		} else {
			fmt.Printf("Host: %s results =>\n%s\n\n", result["name"], result["result"])
		}
		agg_results[result["name"].(string)] = result
	}
	fmt.Println("\n\n")
	fmt.Println("And again, as we stored the results such that we can use outside of the return channel loop.\n")
	for name, results := range agg_results {
		result := results.(map[string]interface{})
		if err, ok := result["error"]; ok {
			fmt.Printf("Host: %s had error %s\n\n", name, err)
		} else {
			fmt.Printf("Host: %s results =>\n%s\n\n", name, result["result"])
		}
	}
	
	//verify host Data updated...
	fmt.Println("\n\n")
	fmt.Println("And lastly verify host data was updated during goroutines.\n")
	for _, host := range hosts {
		fmt.Println(host.Data)
	}
	fmt.Println("\n\n")
}

func F_include(i Hosts, loc string, includes []string) Hosts {

	loc = strings.ToLower(loc)
	if loc == "username" || loc == "password" || loc == "enable" || loc == "strictkey" {
		fmt.Println("I am not programmed to filter on " + loc + ".\n")
	}

	//make a set to hold matches
	hostMatches := make(map[string]bool)

	for _, f_value := range includes {
		r := regexp.MustCompile(f_value)
		switch loc {
		case "name":
			for h,v := range i {
				if r.Match([]byte(v.Name)) {
					hostMatches[h] = true
				}
			}
		case "hostname":
			for h,v := range i {
				if r.Match([]byte(v.Hostname)) {
					hostMatches[h] = true
				}
			}
		case "platform":
			for h,v := range i {
				if r.Match([]byte(v.Platform)) {
					hostMatches[h] = true
				}
			}
		case "groups":
			for h,v := range i {
				f := false
				for _, g := range v.Groups {
					if r.Match([]byte(g)) {
						f = true
						break
					}
				}
				if f {
					hostMatches[h] = true
				}
			}
		default:
			for h,v := range i {
				switch x := v.Data[loc].(type) {
				case nil:
					//when the data key doesn't exist
				case string:
					if r.Match([]byte(v.Data[loc].(string))) {
						hostMatches[h] = true
					}
				case int:
					if r.Match([]byte(strconv.Itoa(v.Data[loc].(int)))) {
						hostMatches[h] = true
					}
				case []interface {}:
					for h,v := range i {
						f := false
						if _, ok := v.Data[loc]; ok {
							for _, g := range v.Data[loc].([]interface{}) {
								if r.Match([]byte(g.(string))) {
									f = true
									break
								}
							}
						}
						if f {
							hostMatches[h] = true
						}
					}
				default:
					//TODO
					fmt.Printf("I don't know how to filter on type %T\n", x)
					os.Exit(0)
				} 
			}
		}
	}

	for h := range(i) {
		if !hostMatches[h] {
			delete(i, h)
		}
	}

	return i
}

func F_exclude(i Hosts, loc string, includes []string) Hosts {

	loc = strings.ToLower(loc)
	if loc == "username" || loc == "password" || loc == "enable" || loc == "strictkey" {
		fmt.Println("I am not programmed to filter on " + loc + ".\n")
	}

	for _, f_value := range includes {
		r := regexp.MustCompile(f_value)
		switch loc {
		case "name":
			for h,v := range i {
				if r.Match([]byte(v.Name)) {
					delete(i,h)
				}
			}
		case "hostname":
			for h,v := range i {
				if r.Match([]byte(v.Hostname)) {
					delete(i,h)
				}
			}
		case "platform":
			for h,v := range i {
				if r.Match([]byte(v.Platform)) {
					delete(i,h)
				}
			}
		case "groups":
			for h,v := range i {
				f := false
				for _, g := range v.Groups {
					if r.Match([]byte(g)) {
						f = true
						break
					}
				}
				if f {
					delete(i,h)
				}
			}
		default:
			for h,v := range i {
				switch x:= v.Data[loc].(type) {
				case nil:
					//we don't care if the data key doesn't exist
				case string:
					if r.Match([]byte(v.Data[loc].(string))) {
						delete(i,h)
					}
				case int:
					if r.Match([]byte(strconv.Itoa(v.Data[loc].(int)))) {
						delete(i,h)
					}
				case []interface {}:
					for h,v := range i {
						f := false
						if _, ok := v.Data[loc]; ok {
							for _, g := range v.Data[loc].([]interface{}) {
								if r.Match([]byte(g.(string))) {
									f = true
									break
								}
							}
						}
						if f {
							delete(i,h)
						}
					}
				default:
					//TODO
					fmt.Printf("I don't know how to filter on type %T\n", x)
					os.Exit(0)
				} 
			}
		}
	}
	return i
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
		msg := fmt.Sprintf("failed to create driver for %s: %+v", h.Hostname, err)
		return c, errors.New(msg)
	}

	err = c.Open()
	if err != nil {
		msg := fmt.Sprintf("failed to open driver for %s: %+v", h.Hostname, err)
		return c,errors.New(msg)
	}

	return c, nil 

}

func getHostsbyYAML() Hosts {

	var i Hosts
	yamlFile, err := ioutil.ReadFile("../inventory/hosts.yaml")
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.UnmarshalStrict(yamlFile, &i)
	if err != nil {
		fmt.Println(err)
	}

	// if any variable doesn't exist use default or create.
	for h,v := range(i) {
		v.Name = h		// copy host key to become host Name
		v.Username = "fred"
		v.Password = "bedrock"
		if v.Platform == "" {
			v.Platform = "cisco_iosxe"
		}
		if len(v.Data) == 0 {
			v.Data = make(map[string]interface{})
		}
	}
	return i

}


func getHostsByCode() Hosts {

	// unused now that YAML working

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
		host.Groups = []string{"group2"}
		host.Data["example_only"] = 100

		hosts[host.Name] = &host
		
	}

	return hosts
}


func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}