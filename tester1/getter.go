package tester1


import (
	"time"
	"sync"
	"log"
	"github.com/scrapli/scrapligo/driver/network"
)


func runTasks(h Host, tasks []string) (map[string]interface{}, error) {

	results := map[string]interface{}{}
	res := make(map[string]interface{})

	var c *network.Driver
	var err error

	for _, task := range tasks {
		switch task {
		case "connection":
			results["connection"] = res
			c, err = Connection(h)
		case "version":
			res, err = Version(h, c)
		case "bgp":
			res, err = Bgp(h, c)
		}

		if err != nil {
			res["failed"] = true
			res["exception"] = err
		}

		results[task] = res

		//stop processing tasks on a connection task error
		if err != nil && task == "connection" {
			return results, err
		}
	}

	if c.Driver.Transport.IsAlive() {
		c.Close()
	}

	return results, nil

}


func Runner(hosts Hosts, tasks []string) map[string]interface{} {

	var wg sync.WaitGroup

	num_workers := 5
	guard := make(chan bool, num_workers)
	results := make(map[string]interface{})
	mux := &sync.Mutex{}
	wg.Add(len(hosts))
	
	// A waitgroup, a mutex, and a channel go the races
	log.Println("Gathering ...")
	defer timeTrack(time.Now())
	for _, host := range hosts {
	
		guard <- true
		go func(h Host) {
			defer wg.Done()
			res, err := runTasks(h, tasks)
			// Print errors immediately
			if err != nil {
				log.Println(err.Error())
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




func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	log.Printf("Runner took %s\n", elapsed)
}