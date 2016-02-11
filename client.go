package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

func main() {
	cert2_b, _ := ioutil.ReadFile("cert2.pem")
	priv2_b, _ := ioutil.ReadFile("cert2.key")
	priv2, _ := x509.ParsePKCS1PrivateKey(priv2_b)

	cert := tls.Certificate{
		Certificate: [][]byte{cert2_b},
		PrivateKey:  priv2,
	}

	roots := x509.NewCertPool()
	pem_ca, err := ioutil.ReadFile("ca.pem")
	pem, _ := x509.ParseCertificate(pem_ca)
	if err != nil {
		panic("failed to read CA pem")
	}

	roots.AddCert(pem)

	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: false,
		RootCAs:            roots,
	}
	conn, err := tls.Dial("tcp", "127.0.0.1:8443", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Println(v.Subject)
	}
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	message := "Hello\n"
	n, err := io.WriteString(conn, message)
	if err != nil {
		log.Fatalf("client: write: %s", err)
	}
	log.Printf("client: wrote %q (%d bytes)", message, n)

	reply := make([]byte, 256)
	n, err = conn.Read(reply)
	log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
	log.Print("client: exiting")
}
