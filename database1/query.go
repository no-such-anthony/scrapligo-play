package main

import (
	"database/sql"
	"fmt"
	"log"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)


func main() {
	db, _ := sql.Open("sqlite3", "./network.db") 
	defer db.Close() 

	hosts := GetHosts(db)
	fmt.Println(hosts)
	
	fmt.Println(GetConfig(db, "R1"))
	
	dataJson := GetData(db, "R1")
	var data map[string][]map[string]interface{}
    json.Unmarshal([]byte(dataJson), &data)
	fmt.Println(data["show_version"])

}

func GetHosts(db *sql.DB) []string {

	row, err := db.Query("SELECT name FROM host ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	i := []string{}
	
	for row.Next() {
		var name string
		row.Scan(&name)
		i = append(i, name)
	}

	return i
}

func GetData(db *sql.DB, name string) string {

	var i string
	row := db.QueryRow("SELECT data FROM data ORDER BY name")
	err := row.Scan(&i)
	if err != nil {
		log.Fatal(err)
	}

	return i
}

func GetConfig(db *sql.DB, name string) string {

	var i string
	row := db.QueryRow("SELECT config FROM data ORDER BY name")
	err := row.Scan(&i)
	if err != nil {
		log.Fatal(err)
	}

	return i
}