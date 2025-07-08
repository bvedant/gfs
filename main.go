package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	dir := flag.String("dir", ".", "directory to serve files from")
	port := flag.Int("port", 8080, "port to listen on")

	flag.Parse()

	fs := http.FileServer(http.Dir(*dir))

	http.Handle("/", loggingMiddleware(fs))

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Serving %s on port %d", *dir, *port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
