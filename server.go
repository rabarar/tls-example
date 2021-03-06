package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	RAW_EXTENSION = ".raw"
	PEM_EXTENSION = ".pem"
)

func main() {

	var useTLS, useCACert bool
	flgCACert := flag.String("ca_cert", "ca_cert", "certificate auth cert file prefix (without .raw or .pem)")
	flgCAKey := flag.String("ca_key", "ca_key", "certificate  key file prefix (without .raw or .pem)")

	flgServerCert := flag.String("s_cert", "s_cert", "certificate auth cert file prefix (without .raw or .pem)")
	flgServerKey := flag.String("s_key", "s_key", "certificate  key file prefix (without .raw or .pem)")
	flag.BoolVar(&useTLS, "tls", false, "use tls if true or https")
	flag.BoolVar(&useCACert, "ca", false, "use CA Cert for ListenAndServer https")

	flag.Parse()

	// Read RAW Cert
	ca_b, err := ioutil.ReadFile(*flgCACert + RAW_EXTENSION)
	if err != nil {
		log.Printf("failed to read %s: %s\n", *flgCACert+RAW_EXTENSION, err)
		os.Exit(1)
	}
	ca, err := x509.ParseCertificate(ca_b)
	if err != nil {
		log.Printf("failed to parse %s: %s\n", *flgCACert+RAW_EXTENSION, err)
		os.Exit(1)
	}

	// Read the RAW key
	priv_b, err := ioutil.ReadFile(*flgCAKey + RAW_EXTENSION)
	if err != nil {
		log.Printf("failed to read %s: %s\n", *flgCAKey+RAW_EXTENSION, err)
		os.Exit(1)
	}

	// Convert the key ...
	priv, err := x509.ParsePKCS1PrivateKey(priv_b)
	if err != nil {
		log.Printf("failed to parse %s: %s\n", *flgCAKey+RAW_EXTENSION, err)
		os.Exit(1)
	}

	// Create a new Cert Pool and add the CA cert
	pool := x509.NewCertPool()
	pool.AddCert(ca)

	// Read the Server Cert
	var cert tls.Certificate
	if useCACert {
		cert = tls.Certificate{
			Certificate: [][]byte{ca_b},
			PrivateKey:  priv,
		}
	} else {

		ca_b, err := ioutil.ReadFile(*flgServerCert + RAW_EXTENSION)
		if err != nil {
			log.Printf("failed to read %s: %s\n", *flgServerCert+RAW_EXTENSION, err)
			os.Exit(1)
		}

		// Read the RAW key
		priv_b, err := ioutil.ReadFile(*flgServerKey + RAW_EXTENSION)
		if err != nil {
			log.Printf("failed to read %s: %s\n", *flgServerKey+RAW_EXTENSION, err)
			os.Exit(1)
		}

		// Convert the key ...
		priv, err := x509.ParsePKCS1PrivateKey(priv_b)
		if err != nil {
			log.Printf("failed to parse %s: %s\n", *flgCAKey+RAW_EXTENSION, err)
			os.Exit(1)
		}

		cert = tls.Certificate{
			Certificate: [][]byte{ca_b},
			PrivateKey:  priv,
		}
	}

	config := tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
	}
	config.Rand = rand.Reader

	if useTLS {
		if useCACert {
			log.Printf("TLS server: listening with Cert:%s, Key:%s\n", *flgCACert+PEM_EXTENSION, *flgCAKey+PEM_EXTENSION)
		} else {
			log.Printf("TLS server: listening with Cert:%s, Key:%s\n", *flgServerCert+PEM_EXTENSION, *flgServerKey+PEM_EXTENSION)
		}

		service := "0.0.0.0:8443"
		listener, err := tls.Listen("tcp", service, &config)
		if err != nil {
			log.Fatalf("server: listen: %s", err)
		}
		log.Print("TLS server: listening")

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
	} else {
		s := &http.Server{
			Addr:           "0.0.0.0:8443",
			Handler:        nil,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			TLSConfig:      &config,
			ConnState:      stateMonitor,
		}

		if useCACert {
			log.Printf("HTTPS server: listening with Cert:%s, Key:%s\n", *flgCACert+PEM_EXTENSION, *flgCAKey+PEM_EXTENSION)
		} else {
			log.Printf("HTTPS server: listening with Cert:%s, Key:%s\n", *flgServerCert+PEM_EXTENSION, *flgServerKey+PEM_EXTENSION)
		}

		// add handler to default service mux
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})

		if useCACert {
			log.Printf("using CA Certificate...\n")
			err := s.ListenAndServeTLS(*flgCACert+PEM_EXTENSION, *flgCAKey+PEM_EXTENSION)
			if err != nil {
				log.Printf("failed to ListenAndServerTLS: %s\n", err)
			}
		} else {
			log.Printf("using Server Certificate...\n")
			err := s.ListenAndServeTLS(*flgServerCert+PEM_EXTENSION, *flgServerKey+PEM_EXTENSION)
			if err != nil {
				log.Printf("failed to ListenAndServerTLS: %s\n", err)
			}
		}

	}
}

var stateToLabel = map[http.ConnState]string{
	http.StateNew:      "New",
	http.StateActive:   "Active",
	http.StateIdle:     "Idle",
	http.StateHijacked: "Hijacked",
	http.StateClosed:   "Closed",
}

func stateMonitor(conn net.Conn, state http.ConnState) {

	label, ok := stateToLabel[state]
	if ok {
		log.Printf(">>> state change[%q]: %s\n", state, label)
	} else {
		log.Printf(">>> state change[%q]: UNKNOWN\n")
	}

	tlscon, ok := conn.(*tls.Conn)
	if ok {
		tlsState := tlscon.ConnectionState()

		/*
			if len(tlsState.PeerCertificates) > 0 {
				sub := tlsState.PeerCertificates[0].Subject
				log.Printf("Connection from CertSubject: %V\n", sub)
			}
		*/

		for i := 0; i < len(tlsState.VerifiedChains); i++ {
			if len(tlsState.VerifiedChains[i]) > 0 {
				for j := 0; j < len(tlsState.VerifiedChains[i]); j++ {
					log.Printf("VerifiedChain[%d] Cert[%d]: %s\n", i, j, tlsState.VerifiedChains[i][j].Subject.Names)
				}
			}
		}
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

		//n, err = conn.Write(buf[:n])
		n, err = conn.Write([]byte("<html><head><title>Done!</title></head><body><H1>Hi there!</H1></Body></HTML>"))
		log.Printf("server: conn: wrote %d bytes", n)

		if err != nil {
			log.Printf("server: write: %s", err)
			break
		}
	}
	log.Println("server: conn: closed")
}
