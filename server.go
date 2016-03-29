package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func main() {

	flgCACert := flag.String("ca", "ca_cert.raw", "certificate auth raw cert file")
	flgCAKey := flag.String("key", "ca_key.raw", "certificate  raw key file")

	flag.Parse()

	ca_b, err := ioutil.ReadFile(*flgCACert)
	if err != nil {
		log.Printf("failed to read %s: %s\n", *flgCACert, err)
		os.Exit(1)
	}
	ca, err := x509.ParseCertificate(ca_b)
	if err != nil {
		log.Printf("failed to parse %s: %s\n", *flgCACert, err)
		os.Exit(1)
	}
	priv_b, err := ioutil.ReadFile(*flgCAKey)
	if err != nil {
		log.Printf("failed to read %s: %s\n", *flgCAKey, err)
		os.Exit(1)
	}
	priv, err := x509.ParsePKCS1PrivateKey(priv_b)
	if err != nil {
		log.Printf("failed to parse %s: %s\n", *flgCAKey, err)
		os.Exit(1)
	}

	pool := x509.NewCertPool()
	pool.AddCert(ca)

	cert := tls.Certificate{
		Certificate: [][]byte{ca_b},
		PrivateKey:  priv,
	}

	config := tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
	}
	config.Rand = rand.Reader
	service := "0.0.0.0:8443"
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}

		tlscon, ok := conn.(*tls.Conn)
		if ok {
			state := tlscon.ConnectionState()
			sub := state.PeerCertificates[0].Subject
			log.Println(sub)
		}

		log.Printf("server: conn: echo %q\n", string(buf[:n]))
		n, err = conn.Write(buf[:n])

		n, err = conn.Write(buf[:n])
		log.Printf("server: conn: wrote %d bytes", n)

		if err != nil {
			log.Printf("server: write: %s", err)
			break
		}
	}
	log.Println("server: conn: closed")
}
