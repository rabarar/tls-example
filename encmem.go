package main

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func main() {
	var blk pem.Block = pem.Block{
		Type:    "CERTIFICATE",
		Headers: nil,
	}

	ca_b, _ := ioutil.ReadFile("ca.pem")
	blk.Bytes = ca_b

	data := pem.EncodeToMemory(&blk)
	fmt.Printf("%s\n", data)
}
