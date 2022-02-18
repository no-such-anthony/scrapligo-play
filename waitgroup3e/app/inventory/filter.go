package inventory

import (
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"os"
)

func Skip(h *Host, include map[string][]string, exclude map[string][]string) bool {

	skip1 := false
	skip2 := false

	if len(include) > 0 {
		skip1 = !F_include(h, include)
	}

	if len(exclude) > 0 {
		skip2 = F_exclude(h, exclude)
	}

	return skip1 || skip2

}

func F_include(h *Host, when map[string][]string) bool {

	hostMatch := false

	for loc, includes := range when {

		loc = strings.ToLower(loc)
		if loc == "username" || loc == "password" || loc == "enable" || loc == "strictkey" {
			fmt.Println("I am not programmed to filter on " + loc + ".\n")
		}

		//make a set to hold matches
		
		for _, f_value := range includes {
			r := regexp.MustCompile(f_value)
			switch loc {
			case "name":
				if r.Match([]byte(h.Name)) {
					hostMatch = true
				} 

			case "hostname":
				if r.Match([]byte(h.Hostname)) {
					hostMatch = true
				} 

			case "platform":
				if r.Match([]byte(h.Platform)) {
					hostMatch = true
				} 

			case "groups":
				f := false
				for _, g := range h.Groups {
					if r.Match([]byte(g)) {
						f = true
						break
					}
				}
				if f {
					hostMatch = true
				} 

			default:
				switch x := h.Data[loc].(type) {
				case nil:
					//when the data key doesn't exist

				case string:
					if r.Match([]byte(h.Data[loc].(string))) {
						hostMatch = true
					} 

				case int:
					if r.Match([]byte(strconv.Itoa(h.Data[loc].(int)))) {
						hostMatch = true
					} 

				case []interface {}:
					f := false
					if _, ok := h.Data[loc]; ok {
						for _, g := range h.Data[loc].([]interface{}) {
							if r.Match([]byte(g.(string))) {
								f = true
								break
							}
						}
					}
					if f {
						hostMatch = true
					}

				default:
					//TODO
					fmt.Printf("I don't know how to filter on type %T\n", x)
					os.Exit(0)
				} 
			}
		}
	}
	return hostMatch
}


func F_exclude(h *Host, not_when map[string][]string) bool {

	for loc, excludes := range not_when {

		loc = strings.ToLower(loc)
		if loc == "username" || loc == "password" || loc == "enable" || loc == "strictkey" {
			fmt.Println("I am not programmed to filter on " + loc + ".\n")
		}

		for _, f_value := range excludes {
			r := regexp.MustCompile(f_value)
			switch loc {

			case "name":
				if r.Match([]byte(h.Name)) {
					return true
				}

			case "hostname":
				if r.Match([]byte(h.Hostname)) {
					return true
				}

			case "platform":
				if r.Match([]byte(h.Platform)) {
					return true
				}

			case "groups":
				f := false
				for _, g := range h.Groups {
					if r.Match([]byte(g)) {
						f = true
						break
					}
				}
				if f {
					return true
				}

			default:
				switch x:= h.Data[loc].(type) {
				case nil:
					//we don't care if the data key doesn't exist
				case string:
					if r.Match([]byte(h.Data[loc].(string))) {
						return true
					}
				case int:
					if r.Match([]byte(strconv.Itoa(h.Data[loc].(int)))) {
						return true
					}
				case []interface {}:
					f := false
					if _, ok := h.Data[loc]; ok {
						for _, g := range h.Data[loc].([]interface{}) {
							if r.Match([]byte(g.(string))) {
								f = true
								break
							}
						}
					}
					if f {
						return true
					}
				default:
					//TODO
					fmt.Printf("I don't know how to filter on type %T\n", x)
					os.Exit(0)
				} 
			}
		}
	}
	return false
}