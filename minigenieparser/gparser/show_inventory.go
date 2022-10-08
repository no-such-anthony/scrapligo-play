package gparser

import (
	//"fmt"
	"strings"
	//"strconv"
	"regexp"
)


// based on python Genie code
// https://github.com/CiscoTestAutomation/genieparser/blob/master/src/genie/libs/parser/iosxe/show_inventory.py


type InvData struct {
	Description   	string
	Pid				string
	Vid				string
	Sn				string
	Oid				string
}

type Show_Inventory map[string]*InvData

func ParseShowInventory(r string) Show_Inventory {

	p1 := regexp.MustCompile(`^NAME: "(?P<name>[\w\d\s(\/\-)?]+)", DESCR: "(?P<description>[\w\d\s(\-\.\/)?,]+)"$`)
	p2 := regexp.MustCompile(`^PID:\s+(?P<pid>[\w\d\-]+|Unknown PID)?\s*,\s+VID:\s+(?P<vid>[\d\w\.]+)?\s*,\s+SN:\s*(?P<sn>[\w\d\-]+)?$`)
	p3 := regexp.MustCompile(`^OID: +(?P<oid>[\d\.]+)$`)

	parsed := make(Show_Inventory)

	currentRecord := ""

	for _, line := range strings.Split(r, "\n") {
		line = strings.TrimSpace(line)
		
		if x := groupmap(line, p1); len(x)!=0 {
			currentRecord = x["name"]
			parsed[currentRecord] = &InvData{}
			parsed[currentRecord].Description = x["description"]
			continue
		}

		if x := groupmap(line, p2); len(x)!=0 {
			// will crash if p1 hasn't matched so let's try to get it to just continue
			if currentRecord == "" {
				continue
			}
			if parsed[currentRecord].Pid != "" {
				continue
			}
			parsed[currentRecord].Pid = x["pid"]
			parsed[currentRecord].Vid = x["vid"]
			parsed[currentRecord].Sn = x["sn"]
			continue
		}

		if x := groupmap(line, p3); len(x)!=0 {
			parsed[currentRecord].Oid = x["oid"]
			continue
		}
		
	}
	return parsed
	
}