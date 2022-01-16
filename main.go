package main

import (
	"github.com/hari-govind/liveserver-go/server"

	"fmt"
)

func main() {
	fmt.Println("Liveserver")
	go server.StartWsServer()
	server.ServeRootDir()
}
