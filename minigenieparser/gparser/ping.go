package gparser

import (
	//"fmt"
	"strings"
	"strconv"
	"regexp"
)


// based on python Genie code
// https://github.com/CiscoTestAutomation/genieparser/blob/master/src/genie/libs/parser/iosxe/ping.py


type Round_Trip struct {
	Valid	bool
	Min_ms	int
	Avg_ms	int
	Max_ms	int
}

type Statistics struct {
	Send					int
	Received				int
	Success_rate_percent	float64
	Round_trip				Round_Trip
}

type Ping struct {
	Address			string
	Data_bytes		int
	Repeat			int
	Timeout_secs	int
	Source			string
	Result_per_line []string
	Statistics		Statistics
}


func ParsePing(r string) Ping {


	p1 := regexp.MustCompile(`Sending +(?P<repeat>\d+), +(?P<data_bytes>\d+)-byte +ICMP +Echos +to +(?P<address>[\S\s]+), +timeout +is +(?P<timeout>\d+) +seconds:`)
	p2 := regexp.MustCompile(`Packet +sent +with +a +source +address +of +(?P<source>[\S\s]+)`)
	p3 := regexp.MustCompile(`[!\.UQM\?&]+`)
	p4 := regexp.MustCompile(`Success +rate +is +(?P<success_percent>\d+) +percent +\((?P<received>\d+)\/(?P<send>\d+)\)(, +round-trip +min/avg/max *= *(?P<min>\d+)/(?P<avg>\d+)/(?P<max>\d+) +(?P<unit>\w+))?`)

	parsed := Ping{}

	for _, line := range strings.Split(r, "\n") {
		line = strings.TrimSpace(line)

		if x := groupmap(line,p1); len(x)!=0 {
			parsed.Address = x["address"]
			parsed.Data_bytes, _ = strconv.Atoi(x["data_bytes"])
			parsed.Repeat, _ = strconv.Atoi(x["repeat"])
			parsed.Timeout_secs, _ = strconv.Atoi(x["timeout"])
			continue
		}

		if x := groupmap(line,p2); len(x)!=0 {
			parsed.Source = x["source"]
			continue
		}

		if x := p3.FindString(line); x!="" {
			parsed.Result_per_line = append(parsed.Result_per_line, x)
			continue
		}

		if x := groupmap(line,p4); len(x)!=0 {
			parsed.Statistics.Success_rate_percent, _ = strconv.ParseFloat(x["success_percent"], 64)
			parsed.Statistics.Received, _ = strconv.Atoi(x["received"])
			parsed.Statistics.Send, _ = strconv.Atoi(x["send"])

			if x["min"] != "" {
				min_ms, _ := strconv.Atoi(x["min"])
				avg_ms, _ := strconv.Atoi(x["avg"])
				max_ms, _ := strconv.Atoi(x["max"])
				if x["unit"] == "s" {
					min_ms *= 1000
					avg_ms *= 1000
					max_ms *= 1000
				}
				parsed.Statistics.Round_trip.Min_ms = min_ms
				parsed.Statistics.Round_trip.Avg_ms = avg_ms
				parsed.Statistics.Round_trip.Max_ms = max_ms
				// not sure how to deal with default values so introduced valid boolean
				parsed.Statistics.Round_trip.Valid = true
			}
			continue
		}
		
	}
	return parsed
	
}