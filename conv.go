package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func main() {

	ca_b, _ := ioutil.ReadFile("ca.pem")
	ca, _ := x509.ParseCertificate(ca_b)
	fmt.Printf("%s\n", ca)
}
