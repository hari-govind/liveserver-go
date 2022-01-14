package watcher

import (
	"github.com/hari-govind/liveserver-go/config"
	"log"
)

func Listen() <-chan string {
	fileWatcher := NewFileWatcher()
	changes := make(chan string, 1024)

	log.Println("Watching for changes in", config.GetConfig().Root)
	fileWatcher.WatchPath(config.GetConfig().Root, config.GetConfig().Depth)

	go func() {
		for e := range fileWatcher.ListenForChanges() {
			changes <- e
		}
	}()

	return changes
}
