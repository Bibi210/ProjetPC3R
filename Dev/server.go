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
	for key, value := range getHandlersMap() {
		http.HandleFunc(key, value)
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server Closed\n")
	} else if err != nil {
		fmt.Printf("Error while starting server: %s\n", err)
		os.Exit(1)
	}
}

func main() {

	bootServer(4242)
}
