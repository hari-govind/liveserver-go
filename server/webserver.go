package server

import (
	"github.com/hari-govind/liveserver-go/config"

	"fmt"
	"log"
	"net/http"
)

func ServeRootDir() {
	listenAddress := fmt.Sprintf("%s:%d",
		config.GetConfig().ListenAddress,
		config.GetConfig().WebserverPort)
	fmt.Printf("\nWebserver location: http://%s\n", listenAddress)
	fsys := injectScriptFileSystem{http.Dir(config.GetConfig().Root)}
	// We use http.ListenAndServe with injectScriptFileSystem, it will inject the reload script to html file when read
	err := http.ListenAndServe(listenAddress, http.FileServer(fsys))
	if err != nil {
		log.Println("Error starting webserver", err)
	}
}
