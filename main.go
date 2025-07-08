package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	dir := flag.String("dir", ".", "directory to serve files from")
	port := flag.Int("port", 8080, "port to listen on")

	flag.Parse()

	fs := http.FileServer(http.Dir(*dir))

	http.Handle("/", fs)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Serving %s on port %d", *dir, *port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
