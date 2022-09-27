package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var HostTableSQL = `CREATE TABLE host (
	"name" TEXT NOT NULL PRIMARY KEY ,
	"hostname" TEXT,
	"platform" TEXT,
	"creds" TEXT		
);`

var DataTableSQL = `CREATE TABLE data (
	"name" TEXT NOT NULL PRIMARY KEY ,
	"config" TEXT,
	"data" TEXT		
);`

var CredsTableSQL = `CREATE TABLE cred (
	"name" TEXT NOT NULL PRIMARY KEY ,
	"username" TEXT,
	"password" TEXT
);`

func main() {
	os.Remove("network.db")

	log.Println("Creating network.db...")
	file, err := os.Create("network.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("network.db created")

	db, _ := sql.Open("sqlite3", "./network.db")
	defer db.Close()

	createHostTable(db)
	createDataTable(db)
	createCredsTable(db)

	insertHost(db, "R1", "192.168.204.101","cisco_iosxe","default")
	insertHost(db, "R2", "192.168.204.102","cisco_iosxe","default")
	insertHost(db, "R3", "192.168.204.103","cisco_iosxe","default")
	insertHost(db, "R4", "192.168.204.104","cisco_iosxe","default")
	insertCred(db, "default", "fred", "bedrock")
	displayHosts(db)
}

func createHostTable(db *sql.DB) {
	
	log.Println("Create host table...")
	statement, err := db.Prepare(HostTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("host table created")
}

func createDataTable(db *sql.DB) {
	
	log.Println("Create data table...")
	statement, err := db.Prepare(DataTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("data table created")
}

func createCredsTable(db *sql.DB) {
	
	log.Println("Create credentials table...")
	statement, err := db.Prepare(CredsTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("credentials table created")
}

func insertHost(db *sql.DB, name string, hostname string, platform string, creds string) {
	log.Println("Inserting host record ...")
	insertHostSQL := `INSERT INTO host(name, hostname, platform, creds) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertHostSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(name, hostname, platform, creds)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayHosts(db *sql.DB) {
	row, err := db.Query("SELECT * FROM host ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() {
		var name string
		var hostname string
		var platform string
		var creds string
		row.Scan(&name, &hostname, &platform, &creds)
		log.Println("Host: ", name, " ", hostname, " ", platform, " ", creds)
	}
}

func insertCred(db *sql.DB, name string, username string, password string) {
	log.Println("Inserting credentials record ...")
	insertHostSQL := `INSERT INTO cred(name, username, password) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertHostSQL)

	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(name, username, password)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
