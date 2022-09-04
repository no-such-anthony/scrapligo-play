package main


import (
	"fmt"
	"time"
	"main/play/app"
	"main/play/plugins/sshscrapli"
	"main/play/plugins/netconfscrapli"
	"main/play/plugins/sshgomiko"
	"main/play/plugins/sshnetrasp"
	"main/play/plugins/other"
	"main/play/plugins/inventory"
)


func main() {
	// To time this process
	defer timeTrack(time.Now())

	hosts := inventory.GetHostsByYAML()

	//fmt.Println(hosts)

	//test global filter, not used in runner
	i := map[string][]string{"hostname": []string{"192.168.204.101","no.suchdomain"},
							 "model": []string{"C3560CX"}}
	x := map[string][]string{"name": []string{"sandbox"}}
	f := app.Filt(hosts, i, x)
	fmt.Println("Filter test:", f)

	//attempt at a simple playbook/runbook/taskbook in code
	command := "show version"
	textfsm := "../textfsm_templates/cisco_iosxe_show_version.textfsm"

	task1 := sshscrapli.SendCommand{
		TaskBase: app.TaskBase{
			Name: "my first show version",
			Include: map[string][]string{"hostname": []string{"192.168.204.101","no.suchdomain"},
									 "model": []string{"C3560CX"}},
			Exclude: map[string][]string{"name": []string{"sandbox"}},
		},
		Command: command,
		Textfsm: textfsm,
	}
	wtask1 := sshscrapli.Wrap{&task1}

	task2 := sshscrapli.SendCommand{
		TaskBase: app.TaskBase{
			Name: "my second show version",
			Exclude: map[string][]string{"name": []string{"192.168.204.101"}},
		},
		Command: command,
		Textfsm: textfsm,
	}
	wtask2 := sshscrapli.Wrap{&task2}

	ncFilter := "" +
	"<interfaces xmlns=\"urn:ietf:params:xml:ns:yang:ietf-interfaces\">\n" +
	"  <interface>\n" +
	"    <name>\n" +
	"      GigabitEthernet1\n" +
	"    </name>\n" +
	"  </interface>\n" +
	"</interfaces>"

	task3 := netconfscrapli.GetConfig{
		TaskBase: app.TaskBase{
			Name: "my netconf show run",
			Include: map[string][]string{"name": []string{"sandbox","r1"}},
		},
		Type: "running",
		Filter: ncFilter,
	}
	wtask3 := netconfscrapli.Wrap{&task3}

	task4 := other.TaskTest{
		TaskBase: app.TaskBase{
			Name: "my default wrappered task test",
		},
		Kwargs: map[string]interface{} { "hello": "defaultwrapped"},
	}
	//tasks.Wrap is default wrapper for tasks not requiring one of the pre-configured connections
	//but nothing stopping you from adding to your task.
	wtask4 := app.Wrap{&task4}

	task5 := sshgomiko.SendCommand{
		TaskBase: app.TaskBase{
			Name: "my gomiko show version",
			Exclude: map[string][]string{"name": []string{"192.168.204.101"}},
		},
		Command: command,
		Textfsm: textfsm,
	}
	wtask5 := sshgomiko.Wrap{&task5}

	task6 := other.TestRestConf{
		TaskBase: app.TaskBase{
			Name: "my restconf test",
			Include: map[string][]string{"name": []string{"sandbox"}},
		},
		Filter: "interface",
	}
	wtask6 := app.Wrap{&task6}

	task7 := sshnetrasp.SendCommand{
		TaskBase: app.TaskBase{
			Name: "my first show version in netrasp",
			Exclude: map[string][]string{"name": []string{"sandbox"}},
		},
		Command: command,
		Textfsm: textfsm,
	}
	wtask7 := sshnetrasp.Wrap{&task7}

	t := []app.Play{&wtask1, &wtask2, &wtask3, &wtask4, &wtask5, &wtask6, &wtask7}
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