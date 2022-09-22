package main

import (
	//"bufio"
	"fmt"
	//"io"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
	"io/ioutil"
	"text/template"
	"bytes"

	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
	"github.com/scrapli/scrapligo/driver/network"
)

type Host struct {
	Name      string
	Hostname  string
	Platform  string
	Connection *network.Driver
	Facts     *Facts
}

type VLAN struct {
	Number	int		`yaml:"number"`
	Name	string	`yaml:"name"`
}

type Facts struct {
	VLAN	[]VLAN	`yaml:"VLAN"`
}

type TclContent struct {
	Filename string
	Content	string
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

	i := ReadFactsYaml("host_vars/R1.yaml")
	//fmt.Printf("%#v\n",i)
	h.Facts = i

	// template loaded from file
	tplates := []string{"templates/vlan.tmpl"}
	tmpl := template.Must(template.New("").ParseFiles(tplates...))

	var configVlan bytes.Buffer
	tmpl = tmpl.Lookup("vlan.tmpl")
	_ = tmpl.Execute(&configVlan, h.Facts)
	//fmt.Println(configVlan.String())

	// template loaded from inline text
	tclContent := TclContent{
		Filename: "disk0:vlan_config.cfg",
		Content: configVlan.String(),
	} 

	tclTemplate := `tclsh
puts [open "{{ .Filename }}" w+] {
{{ .Content -}}
}
tclquit`

	tcpTmpl := template.Must(template.New("TCL").Parse(tclTemplate))
	var tclBuffer bytes.Buffer
	_ = tcpTmpl.Execute(&tclBuffer, tclContent)

	// change return char to do multiline tcl stuff
	d := h.Connection
	d.Channel.ReturnChar = []byte("\r")

	// tidy string
	s := tclBuffer.String()
	s = strings.ReplaceAll(s, "\r\n","\n")
	s = strings.TrimRight(s, "\n")

	// deploy tcl
	tclResult := ""
	for _, line := range strings.Split(s, "\n") {
		r, err := d.SendCommand(line)

		if err != nil {
			fmt.Printf("failed to send command for %s: %+v\n", h.Hostname, err)
			os.Exit(1)
		}

		tclResult += r.Result
	}

	// in this case we don't expect to see any output if successful
	if tclResult != "" {
		fmt.Printf("tcl error for %s: %s\n", h.Hostname, tclResult)
		os.Exit(1)
	}
	
	d.Channel.ReturnChar = []byte("\n")
	r, err := d.SendCommand("more vlan_config.cfg")
	fmt.Println(r.Result)
	fmt.Println("===== EOF =====")

	// tidy tclContent for comparison
	a := tclContent.Content
	a = strings.ReplaceAll(a, "\r\n","\n")
	a = strings.TrimRight(a, "\n")

	// compare
	if strings.Compare(r.Result, a) == 0 {
		fmt.Println("File written to disk matches!")
	} else {
		fmt.Println("File written to disk does NOT match :(")
	}
	
}

func (h *Host) Connect() error {

	p, err := platform.NewPlatform(
		h.Platform,
		h.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername("fred"),
		options.WithAuthPassword("bedrock"),
		options.WithSSHConfigFile("../inventory/ssh_config"),
		// if working in windows with GNS3 and an old 7200 image...
		//options.WithTransportType("standard"),
		//options.WithStandardTransportExtraCiphers([]string{"3des-cbc"}),
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

func ReadFactsYaml(dataFile string) *Facts {

	obj := Facts{}

	yamlFile, err := ioutil.ReadFile(dataFile)
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.UnmarshalStrict(yamlFile, &obj)
	if err != nil {
		fmt.Println(err)
	}

	return &obj
}
