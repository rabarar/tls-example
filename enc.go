package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
)

const ()

func main() {

	secretMsg, err := ioutil.ReadFile("cert2.pem")
	if err != nil {
		fmt.Printf("ReadFile: %s\n", err)
		os.Exit(1)
	}

	blockType := "ENCRYPTED PRIVATE KEY"
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

	// fmt.Printf("raw[%d]:\n%q\n", bs, sDec)
	fmt.Printf("-----BEGIN %s-----\n", blockType)
	for k, v := range EncryptedPEMBlock.Headers {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Printf("\n")

	nblks := bs / 64
	rem := bs % 64

	// fmt.Printf("nBlks = %d, rem = %d\n", nblks, rem)

	for i := 0; i < nblks; i++ {
		fmt.Printf("%s\n", sDec[i*64:(i+1)*64])
	}

	// write the remaining bs-((nblks)*64)
	fmt.Printf("%s\n", sDec[bs-rem:])

	fmt.Printf("-----END %s-----\n", blockType)
}
