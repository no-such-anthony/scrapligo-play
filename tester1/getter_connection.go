package tester1

import (
	"fmt"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/driver/core"
	"github.com/scrapli/scrapligo/driver/network"
)


func Connection(h Host) (*network.Driver, error) {

	c, err := core.NewCoreDriver(
		h.Hostname,
		h.Platform,
		base.WithAuthStrictKey(h.StrictKey),
		base.WithAuthUsername(h.Username),
		base.WithAuthPassword(h.Password),
		//base.WithAuthSecondary(h.Enable),
		//base.WithTransportType("standard"),
		//base.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		return c, fmt.Errorf("failed to create driver: %+v", err)
	}

	err = c.Open()
	if err != nil {
		return c, fmt.Errorf("failed to open driver: %+v", err)
	}

	return c, nil 

}
