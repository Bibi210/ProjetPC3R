package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

func bootServer(port uint16) {
	if port <= 1024 {
		log.Fatal("port is too small")
	}
	for key, value := range HandlersMap() {
		http.HandleFunc(key, value.Handle())
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Println("Server Open")
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server Closed\n")
	} else if err != nil {
		log.Printf("Error while starting server: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	printDatabase()
	bootServer(4242)
}
