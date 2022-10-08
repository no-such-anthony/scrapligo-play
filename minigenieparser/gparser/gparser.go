package gparser

import (
	"regexp"
	//"fmt"
)


func groupmap(s string, r *regexp.Regexp) map[string]string {
    values := r.FindStringSubmatch(s)
    keys := r.SubexpNames()
    // create map
    d := make(map[string]string)
	if values != nil {
		for i := 1; i < len(keys); i++ {
			d[keys[i]] = values[i]
		}
	}
    return d
}