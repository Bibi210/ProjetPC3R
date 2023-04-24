package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var cleanDatabase = true

func bootServer(port uint16) {
	if port <= 1024 {
		log.Fatal("port is too small")
	}
	for key, value := range HandlersMap() {
		http.HandleFunc(key, value.ToHandler())
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("frontend/dist/assets"))))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Println("Server Open")
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server Closed\n")
	} else if err != nil {
		log.Printf("Error while starting server: %s\n", err)
		os.Exit(1)
	}
}

func shutdownListener() {
	log.Println("Shutdown Listener")
	// Catch SIGINT to shutdownDatabase()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Server Shutdown")
	shutdownDatabase(cleanDatabase)
	os.Exit(0)
}

func main() {
	createDatabase()
	go shutdownListener()
	bootServer(25565)
	shutdownDatabase(cleanDatabase)
}
