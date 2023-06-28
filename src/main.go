package main

import "github.com/moxxteroxxte1/stafftime-backend/src/handlers"

func main() {
	server := handlers.NewAPIServer("0.0.0.0:3000")
	server.Start()
}
