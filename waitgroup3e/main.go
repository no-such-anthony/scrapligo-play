package main


import (
	"fmt"
	"time"
	"main/app"
	"main/app/tasks"
	"main/app/sshscrapli"
	"main/app/netconfscrapli"
	"main/app/inventory"
)


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := inventory.GetHosts()
	//fmt.Println(hosts)

	//test global filter, not used in runner
	i := map[string][]string{"name": []string{"192.168.204.101","no.suchdomain"}}
	x := map[string][]string{"name": []string{"sandbox"}}
	f := inventory.Filt(hosts, i, x)
	fmt.Println(f)

	//attempt at a simple playbook/runbook/taskbook in code
	task1 := sshscrapli.ShowVersion{
		Name: "my first show version",
		Kwargs: map[string]interface{} { "hello": "first"},
		Include: map[string][]string{"name": []string{"192.168.204.101","no.suchdomain"}},
		Exclude: map[string][]string{"name": []string{"sandbox"}},
	}
	wtask1 := sshscrapli.ScrapliSSHWrap{Tasker: &task1}

	task2 := sshscrapli.ShowVersion{
		Name: "my second show version",
		Kwargs: map[string]interface{} { "hello": "second"},
		Exclude: map[string][]string{"name": []string{"192.168.204.101","sandbox"}},
	}
	wtask2 := sshscrapli.ScrapliSSHWrap{Tasker: &task2}

	task3 := netconfscrapli.Running{
		Name: "my netconf show run",
		Kwargs: map[string]interface{} { "hello": "netconf"},
		Include: map[string][]string{"name": []string{"sandbox"}},
	}
	wtask3 := netconfscrapli.ScrapliNetconfWrap{Tasker: &task3}

	task4 := tasks.DefaultTaskTest{
		Name: "my default wrappered task test",
		Kwargs: map[string]interface{} { "hello": "defaultwrapped"},
	}
	//default wrapper for tasks not requiring a connection
	wtask4 := tasks.DefaultWrap{Tasker: &task4}

	t := []tasks.Wrapper{&wtask1, &wtask2, &wtask3, &wtask4}
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