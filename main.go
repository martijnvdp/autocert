package main

import (
	"flag"
	"fmt"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	domain := flag.String("domain", "", "Domain/Common name for certificate")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello HTTP/2")
	})

	fmt.Println("TLS domain", domain)
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(*domain),
		Cache:      autocert.DirCache("certs"),
	}

	tlsConfig := certManager.TLSConfig()
	server := http.Server{
		Addr:      ":443",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	fmt.Println("Server listening on", server.Addr)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		fmt.Println(err)
	}
}
