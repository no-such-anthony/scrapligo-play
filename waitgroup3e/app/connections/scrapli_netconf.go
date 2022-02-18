package connections

import (
	"fmt"
	"main/app/inventory"
	"github.com/scrapli/scrapligo/driver/base"
	"github.com/scrapli/scrapligo/netconf"
)

type ScrapligoNetconf struct {
	C *netconf.Driver	
}


func (s ScrapligoNetconf) Close() {
	s.C.Close()
}


func (s *ScrapligoNetconf) Open(h *inventory.Host) (error) {

	c, err := netconf.NewNetconfDriver(
		h.Hostname,
		//h.Platform,
		base.WithAuthStrictKey(h.StrictKey),
		base.WithAuthUsername(h.Username),
		base.WithAuthPassword(h.Password),
		//base.WithAuthSecondary(h.Enable),
		//base.WithTransportType("standard"),
		//base.WithSSHConfigFile("ssh_config"),
	)

	if err != nil {
		return fmt.Errorf("failed to create driver for %s: %+v", h.Hostname, err)
	}

	err = c.Open()
	if err != nil {
		return fmt.Errorf("failed to open driver for %s: %+v", h.Hostname, err)
	}

	s.C = c
	return nil 

}