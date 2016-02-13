package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	secretMsg = `pE8ATwG3z9ISOYG0LgNc8BeCKSvEKk6Ic+3yWp87VeI=`
)

func main() {

	blockType := "RSA PRIVATE KEY"
	password := []byte("password")
	var headers map[string]string = map[string]string{
		"Proc-Type": "4,ENCRYPTED",
		"DEK-Info":  "AES-128-CBC,6acb08aae5e8ea919311556d7894f224",
	}

	sDec, err := base64.StdEncoding.DecodeString(secretMsg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("secret is\n%s\n\n", sDec)

	var b pem.Block = pem.Block{
		Type:    blockType,
		Headers: headers,
		Bytes:   sDec,
	}

	DecryptPEMBlock, err := x509.DecryptPEMBlock(&b, password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", DecryptPEMBlock)
}
