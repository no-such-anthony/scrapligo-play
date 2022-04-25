package main
// A netrasp only example

import (
	"fmt"
	"time"
	"main/play/app"
	"main/play/plugins/sshnetrasp"
	"main/play/plugins/inventory"
)


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := inventory.GetHostsByYAML()
	//fmt.Println(hosts)

	command := "show version"
	textfsm := "../textfsm_templates/cisco_iosxe_show_version.textfsm"

	task1 := sshnetrasp.SendCommand{
		Name: "a show version in netrasp",
		Command: command,
		Textfsm: textfsm,
		//Exclude: map[string][]string{"name": []string{"sandbox"}},
	}
	wtask1 := sshnetrasp.Wrap{Tasker: &task1}

	t := []app.Wrapper{&wtask1}
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