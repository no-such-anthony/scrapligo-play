package main

import (
	"database/sql"
	"log"
	"sync"
	//"os"
	"fmt"
	"encoding/json"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/scrapli/scrapligo/driver/network"
	"github.com/scrapli/scrapligo/driver/options"
	"github.com/scrapli/scrapligo/platform"
)

type Credentials struct {
	Username  string
	Password  string
}

type Host struct {
	Name      string
	Hostname  string
	Platform  string
	Creds	  Credentials
	Connection *network.Driver
}

type Hosts map[string]*Host


func main() {
	db, _ := sql.Open("sqlite3", "./network.db") 
	defer db.Close() 

	hosts := GetHosts(db)
	theGathering(hosts, db)

}


func theGathering(hosts Hosts, db *sql.DB) {

	var wg sync.WaitGroup
	num_workers := 5
	guard := make(chan bool, num_workers)

	//Note: Combining Waitgroup with a channel to restrict number of goroutines.
	wg.Add(len(hosts))
	for _, host := range hosts {
		guard <- true
		go func(h Host) {
			defer wg.Done()
			getStuff(&h, db)
			<-guard
		}(*host)
	}
	wg.Wait()

}

func getStuff(h *Host, db *sql.DB) {

	data := make(map[string][]map[string]interface{})
	config := ""

	err := h.Connect()
	if err != nil {
		fmt.Printf("failed to connect to %s: %+v\n\n", h.Hostname, err)
		return
	}
	defer h.Disconnect()
	d := h.Connection

	iosData := []string{ "show version",
						 "show cdp neighbors detail",
						 "show inventory",
						 "show ip interface brief",
						}

	for _, cmd := range iosData {
		rs, err := d.SendCommand(cmd)
		if err != nil {
			fmt.Printf("failed to send command for %s: %+v\n\n", h.Hostname, err)
			return
		}

		cmdKey := strings.ReplaceAll(cmd, " ", "_")
		parsedOut, err := rs.TextFsmParse("textfsm/cisco_ios_" + cmdKey + ".textfsm")
		if err != nil {
			fmt.Printf("failed to parse command for %s: %+v\n\n", h.Hostname, err)
			return
		}

		data[cmdKey] = parsedOut
	}
	jsonDataStr, err := json.Marshal(data)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//fmt.Println(string(jsonDataStr))

	rs, err := d.SendCommand("show running-config")
	if err != nil {
		fmt.Printf("failed to send command for %s: %+v\n\n", h.Hostname, err)
		return
	}
	config = rs.Result
	//fmt.Println(config)

	insertData(db, h.Name, config, jsonDataStr)


}

func insertData(db *sql.DB, name string, config string, data []byte) {
	log.Println("Inserting data record ...")
	insertDataSQL := `INSERT OR REPLACE INTO data(name, config, data) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertDataSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(name, config, data)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func GetHosts(db *sql.DB) Hosts {

	row, err := db.Query("SELECT * FROM cred ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}

	c := make(map[string]Credentials)
	for row.Next() {
		var name string
		var username string
		var password string
		row.Scan(&name, &username, &password)
		c[name] = Credentials{Username: username, Password: password}
	}
	row.Close()


	row, err = db.Query("SELECT * FROM host ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	i := make(Hosts)
	
	for row.Next() {
		var h Host
		var name string
		var hostname string
		var platform string
		var creds string
		row.Scan(&name, &hostname, &platform, &creds)
		h.Name = name
		h.Hostname = hostname
		h.Platform = platform
		h.Creds = c[creds]
		i[name] = &h
		log.Println("Host: ", name, " ", hostname, " ", platform, " ", creds)
	}

	return i
}

func (h *Host) Connect() error {

	p, err := platform.NewPlatform(
		h.Platform,
		h.Hostname,
		options.WithAuthNoStrictKey(),
		options.WithAuthUsername(h.Creds.Username),
		options.WithAuthPassword(h.Creds.Password),
		//options.WithTransportType("standard"),
		//options.WithStandardTransportExtraCiphers([]string{"3des-cbc"}),
		options.WithSSHConfigFile("../inventory/ssh_config"),
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
