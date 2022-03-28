package tester1_test

import (
	"testing"
	//"log"
)


func testConnection(host string, res interface{}) func(*testing.T) {

	return func(t *testing.T) {

		hr := res.(map[string]interface{})
		e, ok := hr["failed"].(bool)
		if !ok {
			// key doesn't exist so can't have failed...right?
			return
		}

		if e {
			t.Error(hr["exception"])
		}

	}
}