package main


import (
	"fmt"
	"time"
	"main/app"
	"main/app/tasks"
	"main/app/inventory"
)


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := inventory.GetHosts()
	//fmt.Println(hosts)

	//attempt at a simple playbook/runbook/taskbook in code
	task1 := tasks.ShowVersion{
		Name: "my first show version",
		Kwargs: map[string]interface{} { "hello": "first"},
		Include: map[string][]string{"name": []string{"192.168.204.101","no.suchdomain"}},
		Exclude: map[string][]string{"name": []string{"sandbox"}},
	}

	task2 := tasks.ShowVersion{
		Name: "my second show version",
		Kwargs: map[string]interface{} { "hello": "second"},
		Exclude: map[string][]string{"name": []string{"192.168.204.101","sandbox"}},
	}

	task3 := tasks.NetconfShowVersion{
		Name: "my netconf show version",
		Kwargs: map[string]interface{} { "hello": "netconf"},
		Include: map[string][]string{"name": []string{"sandbox"}},
	}


	t := []tasks.Tasker{&task1, &task2, &task3}
	//fmt.Printf("%+v\n", t)

	results := app.Runner(hosts, t)

	fmt.Print("\n\n")
	fmt.Println("======================= RESULTS =================================")
	for n, h := range results {
		fmt.Println("Name:", n)
		for _, res := range h.([]map[string]interface{}) {
			fmt.Println(res)
		}
		fmt.Print("\n\n")
	}

}


func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("This process took %s\n", elapsed)
}