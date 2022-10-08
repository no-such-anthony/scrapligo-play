package main

import (
	"fmt"
	"main/configparser"
)

func main() {

    cp := configparser.Scan("R1.txt")

	for _, n := range cp.FindNodes(`^interface`) {
		fmt.Println(n.Line)
	}
	
}