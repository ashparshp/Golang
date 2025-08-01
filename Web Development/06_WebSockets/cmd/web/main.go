package main

import (
	"log"
	"net/http"
)

func main() {
	routes := routes()

	log.Println("Starting web server on port 8081")

	_ = http.ListenAndServe(":8081", routes)
}
