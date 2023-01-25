package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Default().Println("Starting...")
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	log.Default().Printf("Request: %s\n", r.URL.Path[1:])
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
