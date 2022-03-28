package tester1_test

import (
	"testing"
	//"log"
)

func testBgp(host string, res interface{}) func(*testing.T) {

	return func(t *testing.T) {

		hr := res.(map[string]interface{})

		xhr, ok := hr["result"].([]map[string]interface{})
		if !ok {
			t.Skip("Skipped due to no results")
			return
		}

		if len(xhr) == 0 {
			t.Error("No BGP neighbor information in results")
		}

		for _, b := range xhr {

			n, ok := b["BGP_NEIGH"]
			if !ok {
				t.Error("No bgp neighbor key found in results")
				return
			}

			t.Run(n.(string), func(t *testing.T) {

				v, ok := b["STATE_PFXRCD"]
				if !ok {
					t.Error("No bgp state key found in results")
					return
				}

				switch v {
				case "0":
					t.Error("BGP Estatblished but zero learned routes")
				case "Idle", "Active", "Idle (Admin)":
					t.Error("BGP in Idle/Active state")
				}
			})
		}
	}
}