package tester1_test

import (
	"testing"
	"os"
	"log"
	"flag"
	"strings"
	"tester1"
)


var (
	d  map[string]interface{}
	tasks []string
)

type arrayValue []string

func (i *arrayValue) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayValue) Set(s string) error {
   	*i = append(*i, strings.TrimSpace(s))
	return nil
}

func TestMain(m *testing.M) {

	// Stuff before tests
	tasks = []string{"version","bgp"}

	setName := arrayValue{}
	setTask := arrayValue{}
	flag.Var(&setName, "name", "device names")
	flag.Var(&setTask, "task", "task names")
	flag.Parse()

	log.Printf("Names: %v", setName.String())
	log.Printf("Tasks: %v", setTask.String())

	if len(setTask) > 0 {
		tasks = setTask
	}

	// Connection happens first, so prepend
	tasks = append([]string{"connection"}, tasks...)

	d = gathering(setName)

	// Tests
	exitVal := m.Run()

	// Stuff after tests

	// Exit
	os.Exit(exitVal)
}


func TestHost(t *testing.T) {

	for host, res := range d {
		t.Run(host, testHost(host, res))		
	}
}


func testHost(host string, res interface{}) func(*testing.T) {

	return func(t *testing.T) {
		xres := res.(map[string]interface{})
		for _, task := range tasks {
			t.Run(task, getTest(host, task, xres))
		}
	}

}

func getTest(host string, task string, res map[string]interface{}) func(*testing.T) {

	tr, ok := res[task]
	if !ok {
		//lets just have empty map for now
		tr = make(map[string]interface{})
	}

	switch task {
	case "connection":
		return testConnection(host, tr)
	case "version":
		return testVersion(host, tr)
	case "bgp":
		return testBgp(host, tr)
	}

	return func(t *testing.T) {
		t.Error("Unknown task")
	}
 
}


func gathering(filt []string) map[string]interface{} {

	hosts := tester1.GetHosts()

	if len(filt) > 0 {
		hosts = tester1.Filter(hosts, filt)
	}

	results := tester1.Runner(hosts, tasks)
	return results

}

