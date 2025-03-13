package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	port := "8080"
	log.Printf("server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprint(w, "Olá, mundo!")
	if err != nil {
		http.Error(w, "error printing hello world", http.StatusInternalServerError)
		return
	}
}
