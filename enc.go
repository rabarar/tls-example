package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"
)

const (
	secretMsg = "happy days are here again!!!"
)

func main() {
	blockType := "RSA PRIVATE KEY"
	password := []byte("password")

	// see http://golang.org/pkg/crypto/x509/#pkg-constants
	cipherType := x509.PEMCipherAES128

	EncryptedPEMBlock, err := x509.EncryptPEMBlock(rand.Reader,
		blockType,
		[]byte(secretMsg),
		password,
		cipherType)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sDec := base64.StdEncoding.EncodeToString(EncryptedPEMBlock.Bytes)
	bs := len(sDec)
	fmt.Printf("-----BEGIN %s-----\n", blockType)
	for k, v := range EncryptedPEMBlock.Headers {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Printf("\n")

	for i := 0; i < bs; i++ {
		if i+64 < bs {
			for j := i; j < i+64; j++ {
				fmt.Printf("%c", sDec[j])
			}
			fmt.Printf("\n")
			i += 64
		} else {
			//fmt.Printf("i = %d\n", i)
			for j := 0; j < bs-i; j++ {
				fmt.Printf("%c", sDec[i+j])
			}
			fmt.Printf("\n")

			break
		}
	}

	fmt.Printf("-----END %s-----\n", blockType)
}
