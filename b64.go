package main

import (
	"encoding/base64"
	"fmt"
	"log"
)

const (
	str = `
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
helo helo helo helo helo helo
`
)

func main() {
	str := base64.StdEncoding.EncodeToString([]byte(str))
	fmt.Println(str)

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Fatal("error:", err)
	}

	fmt.Printf("%s\n", data)
}
