package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"

	"golang.org/x/net/websocket"
)

var originStr = "http://127.0.0.1/"
var urlStr = "wss://127.0.0.1:8080/echo"

func main() {

	locationURL, _ := url.Parse(urlStr)
	originURL, _ := url.Parse(originStr)

	// Load CA cert
	cert2_b, _ := ioutil.ReadFile("client_cert.raw")
	priv2_b, _ := ioutil.ReadFile("client_key.raw")
	priv2, _ := x509.ParsePKCS1PrivateKey(priv2_b)

	cert := tls.Certificate{
		Certificate: [][]byte{cert2_b},
		PrivateKey:  priv2,
	}

	roots := x509.NewCertPool()
	pem_ca, err := ioutil.ReadFile("ca_cert.raw")
	pem, _ := x509.ParseCertificate(pem_ca)
	if err != nil {
		panic("failed to read CA pem")
	}

	roots.AddCert(pem)

	var tlsConfig = tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
		RootCAs:            roots,
	}

	var wsConfig = websocket.Config{
		Location:  locationURL,
		Origin:    originURL,
		TlsConfig: &tlsConfig,
		Protocol:  []string{""},
		Version:   13,
	}

	ws, err := websocket.DialConfig(&wsConfig)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("hello, world!")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)

	var msg = make([]byte, 512)
	_, err = ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg)
}
