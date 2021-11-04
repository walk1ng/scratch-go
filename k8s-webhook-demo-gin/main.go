package main

import (
	"flag"
	"fmt"
	"log"

	"webhook/router"
)

var (
	port     int
	certFile string
	keyFile  string
)

func main() {
	flag.IntVar(&port, "port", 4443, "webhook listen port")
	flag.StringVar(&certFile, "tlsCertFile", "/etc/tls/tls.crt", "")
	flag.StringVar(&keyFile, "tlsKeyFile", "/etc/tls/tls.key", "")
	flag.Parse()

	r := router.NewRouter()
	if err := r.RunTLS(fmt.Sprintf(":%d", port), certFile, keyFile); err != nil {
		log.Fatal(err)
	}
}
