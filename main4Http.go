package main

import (
	"io"
	"log"
	"net/http"
)

func HelloServer(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello ,this is from HelloServer func ")
}

func main() {
	http.HandleFunc("/hello", HelloServer)
	log.Fatal(http.ListenAndServe(":12345", nil))
}
