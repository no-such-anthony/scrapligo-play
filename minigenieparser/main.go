package main

import (
	"fmt"
	"main/gparser"
)

// based on python Genie code
// https://github.com/CiscoTestAutomation/genieparser/blob/master/src/genie/libs/parser/iosxe

func main() {

	r := `Sending 10, 100-byte ICMP Echos to 10.229.1.1, timeout is 2 seconds:
Packet sent with a source address of 10.229.1.2
!!!!!!!!!!!!!!
!.UQM?&!.UQM?&
Success rate is 100 percent (100/100), round-trip min/avg/max = 1/2/14 ms`
//Success rate is 0 percent (0/10)`

	parsePing := gparser.ParsePing(r)
	fmt.Println(parsePing.Address, parsePing.Statistics.Success_rate_percent)

	r = `NAME: "Chassis", DESCR: "Cisco 7206VXR, 6-slot chassis"
PID: CISCO7206VXR      , VID:    , SN: 4279256517
	
NAME: "NPE400 0", DESCR: "Cisco 7200VXR Network Processing Engine NPE-400"
PID: NPE-400           , VID:    , SN: 11111111`

	parseShowInventory := gparser.ParseShowInventory(r)
	fmt.Println(parseShowInventory["Chassis"].Pid)

}