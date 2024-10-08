package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Hello World from Hermione")

	http.ListenAndServe(":8080", nil)
}
