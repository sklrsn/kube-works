package main

import (
	"log"
	"net/http"
)

func init() {
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Success"))
	})
	log.Fatalf("%v", http.ListenAndServe(":8080", nil))
}
