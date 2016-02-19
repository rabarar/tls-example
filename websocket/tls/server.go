package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

const (
	CERT_NAME     = "client_cert.pem"
	KEY_NAME      = "client_key.pem"
	BASIC_HANDLER = false
)

func EchoHandler(ws *websocket.Conn) {
	log.Printf("handler")
	io.Copy(ws, ws)
}

func echoHandler(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:n])

	m, err := ws.Write(msg[:n])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", msg[:m])
}

func main() {

	if BASIC_HANDLER {
		http.Handle("/echo", websocket.Handler(echoHandler))
		err := http.ListenAndServeTLS(":8080", "client_cert.pem", "client_key.pem", nil)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	} else {

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			switch r.URL.Path {
			case "/":
				u := func() string {
					if r.TLS == nil {
						return "ws://localhost:4080/echo"
					} else {
						return "wss://localhost:8080/echo"
					}
				}()
				http.ServeContent(w, r, "index.html", time.Now(), bytes.NewReader([]byte(Content(u))))
			case "/echo":
				websocket.Handler(EchoHandler).ServeHTTP(w, r)
			default:
				http.NotFound(w, r)
			}
		})

		done := make(chan bool)

		go func() {

			var tlsConfig tls.Config = tls.Config{
				InsecureSkipVerify: true,
			}

			var server http.Server = http.Server{
				Addr:      ":8080",
				TLSConfig: &tlsConfig,
			}

			fmt.Println(server.ListenAndServeTLS(CERT_NAME, KEY_NAME))
			done <- true
		}()

		<-done
	}

}

func Content(ws string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang='en'>
  <head>
    <meta charset='utf-8'/>
    <title>websocket test</title>
    <script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.2/jquery.min.js"></script>  
  </head>
  <body>
    <p> go-lang / chrome websocket test at %s: </p>
    <ol/>
    <script>

var log = function(e) {
    $('ol').append("<li>"+JSON.stringify(e)+"</li>");
}

var sock = new WebSocket("%s")

sock.onerror = function(e) { log({ERROR:e}) };    
sock.onclose = function(e) { log({CLOSE:e}) };    
sock.onmessage = function(e) { log({MESSAGE:e.data}) };    

sock.onopen = function(e) { 
    log({OPEN:e}) 
    setInterval(function() {
	sock.send("howdy at " + new Date());
    },1000);
};    

    </script>
  </body>
</html>
`, ws, ws)
}
