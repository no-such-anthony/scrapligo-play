package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"os"
	"text/template"
	"text/tabwriter"
	"bytes"
	"time"
	"github.com/sirikothe/gotextfsm"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/opoptions"
	"github.com/scrapli/scrapligo/response"
)

type Host struct {
	Name      string
	Hostname  string
	Platform  string
	Connection *network.Driver

}

func main() {

	h := Host{
		Name: "R1",
		Hostname: "192.168.204.101",
		Platform: "cisco_iosxe",
	}

	err := h.Connect()
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer h.Disconnect()
	d := h.Connection

	ipAddresses := []string{"10.10.0.11","8.8.8.8","192.168.201.1"}

	tclTemplate := `tclsh
foreach addr {
{{ range . -}}
{{ . }}
{{ end -}}
} { ping $addr repeat 3
}
tclquit`

	tcpTmpl := template.Must(template.New("TCL").Parse(tclTemplate))
	var tclBuffer bytes.Buffer
	_ = tcpTmpl.Execute(&tclBuffer, ipAddresses)

	fmt.Println(tclBuffer.String())

	//change return char to do multiline tcl stuff
	d.Channel.ReturnChar = []byte("\r")
	// tidy string
	s := tclBuffer.String()
	s = strings.ReplaceAll(s, "\r\n","\n")
	s = strings.TrimRight(s, "\n")

	// deploy tcl
	tclResult := ""
	var r *response.Response
	for _, line := range strings.Split(s, "\n") {

		if line == "}" {
			// should be when ping starts, so let us add timeout
			r, err = d.SendCommand(line, opoptions.WithTimeoutOps(120 * time.Second))
		} else {
			r, err = d.SendCommand(line)
		}

		if err != nil {
			fmt.Printf("failed to send command for %s: %+v\n", h.Hostname, err)
			os.Exit(1)
		}

		if line == "}" {
			tclResult += r.Result
		}
	}

	d.Channel.ReturnChar = []byte("\n")

	fmt.Println()
	fmt.Println(tclResult)
	fmt.Println()
	
	textfsm, err := ioutil.ReadFile("ping_ios.textfsm")
	if err != nil {
		fmt.Printf("SendCommand: Error opening template '%s'", "ping_ios.textfsm")
		os.Exit(1)
	}

	fsm := gotextfsm.TextFSM{}
	err = fsm.ParseString(string(textfsm))
	if err != nil {
		fmt.Printf("SendCommand: Error while parsing template '%s'","ping_ios.textfsm")
		os.Exit(1)
	}
	
	parser := gotextfsm.ParserOutput{}
	err = parser.ParseTextString(tclResult, fsm, true)
	if err != nil {
		fmt.Printf("SendCommand: Error while parsing input for template '%s'","ping_ios.textfsm")
		os.Exit(1)
	}
	
	// initialize tabwriter
	w := new(tabwriter.Writer)
	
	// minwidth, tabwidth, padding, padchar, flags
	w.Init(os.Stdout, 12, 8, 0, '\t', 0)
	
	fmt.Fprintf(w, "%s\t%s\t%s\t\n", "IP", "SUCCESS", "RTT")
	fmt.Fprintf(w, "%s\t%s\t%s\t\n", "", "%", "min/avg/max")
	fmt.Fprintf(w, "%s\n", strings.Repeat("-",50))
	
	for _, r := range parser.Dict {
        	fmt.Fprintf(w, "%s\t%s\t%s\t\n", r["IP"], r["SUCCESS"], r["RTT"])
    }
	w.Flush()

}

func (h *Host) Connect() error {

	p, err := platform.NewPlatform(
		h.Platform,
		h.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("fred"),
		options.WithAuthPassword("bedrock"),
		options.WithTransportType("standard"),
		options.WithStandardTransportExtraCiphers([]string{"3des-cbc"}),
		//options.WithReturnChar("\r"),
		//options.WithSSHConfigFile("ssh_config"),
	)
	if err != nil {
        return fmt.Errorf("failed to create platform; error: %+v\n", err)
    }

	d, err := p.GetNetworkDriver()
	if err != nil {
		return fmt.Errorf("failed to create driver for %s: %+v\n\n", h.Hostname, err)
	}

	err = d.Open()
	if err != nil {
		return fmt.Errorf("failed to open driver for %s: %+v\n\n", h.Hostname, err)
	}

	h.Connection = d

	return nil
}

func (h *Host) Disconnect() {

	if h.Connection.Driver.Transport.IsAlive() {
		h.Connection.Close()
	}
}
