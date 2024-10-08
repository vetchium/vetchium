package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Hello World from Granger")

	http.ListenAndServe(":8080", nil)
}
