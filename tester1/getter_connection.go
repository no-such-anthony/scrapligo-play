package tester1

import (
	"fmt"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/driver/network"
)


func Connection(h Host) (*network.Driver, error) {

	p, err := platform.NewPlatform(
		h.Platform,
		h.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(h.Username),
		options.WithAuthPassword(h.Password),
		options.WithSSHConfigFile("../inventory/ssh_config"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create platform for %s: %+v\n\n", h.Hostname, err)
	}

	c, err := p.GetNetworkDriver()
	if err != nil {
        return nil, fmt.Errorf("failed to fetch network driver for %s: %+v\n\n", h.Hostname, err)
	}   

	err = c.Open()
	if err != nil {
		return c, fmt.Errorf("failed to open driver: %+v", err)
	}

	return c, nil 

}
