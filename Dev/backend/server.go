package main

import (
	"Backend/Database"
	"Backend/Helpers"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var cleanDatabase = false

func bootServer(port uint16) {
	if port <= 1024 {
		log.Fatal("port is too small")
	}
	for key, value := range HandlersMap() {
		http.HandleFunc(key, value.ToHandler())
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../frontend/dist/assets"))))
	log.Printf("Server Open on http://localhost:%d\n", port)
	Helpers.ServerRuntimeError("Can't open server", http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func shutdownListener() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Server Shutdown")
	Database.ShutdownDatabase(cleanDatabase)
	os.Exit(0)
}

func main() {
	Database.CreateDatabase()
	go shutdownListener()
	db := Database.OpenDatabase()
	Database.ShowDatabase(db)
	Database.CleanCloser(db)
	bootServer(25565)
}
