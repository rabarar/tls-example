package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var useSkip bool
	flag.BoolVar(&useSkip, "skip", false, "Skip Insecure Verification")
	flgCACert := flag.String("ca", "ca_cert.raw", "client raw cert file")
	flgClientCert := flag.String("cert", "client_cert.raw", "client raw cert file")
	flgClientKey := flag.String("key", "client_key.raw", "client key file")

	flag.Parse()

	cert2_b, err := ioutil.ReadFile(*flgClientCert)
	if err != nil {
		log.Printf("failed to read %s: %s\n", *flgClientCert, err)
		os.Exit(1)
	}
	priv2_b, err := ioutil.ReadFile(*flgClientKey)
	if err != nil {
		log.Printf("failed to read %s: %s\n", *flgClientKey, err)
		os.Exit(1)
	}
	priv2, err := x509.ParsePKCS1PrivateKey(priv2_b)
	if err != nil {
		log.Printf("failed to parse %s: %s\n", *flgClientKey, err)
		os.Exit(1)
	}

	cert := tls.Certificate{
		Certificate: [][]byte{cert2_b},
		PrivateKey:  priv2,
	}

	roots := x509.NewCertPool()
	pem_ca, err := ioutil.ReadFile(*flgCACert)
	if err != nil {
		log.Printf("failed to read %s: %s\n", *flgCACert, err)
		os.Exit(1)
	}
	pem, err := x509.ParseCertificate(pem_ca)
	if err != nil {
		log.Printf("failed to Parse Cert %s: %s\n", *flgCACert, err)
		os.Exit(1)
	}

	roots.AddCert(pem)

	config := tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: useSkip,
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

	message := "Hello from the client\n"
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
