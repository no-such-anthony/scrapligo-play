package tester1_test

import (
	"testing"
	//"log"
)


func testVersion(host string, res interface{}) func(*testing.T) {

	return func(t *testing.T) {

		hr := res.(map[string]interface{})

		xhr, ok := hr["result"].(map[string]interface{})
		if !ok {
			t.Skip("Skipped due to no results")
			return
		}

		v, ok := xhr["VERSION"]
		if !ok {
			t.Error("No version information found in results")
			return
		}

		if v != "15.2(4)M11" {
			t.Error("Incorrect Software Version")
		}

	}
	
}