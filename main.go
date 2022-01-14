package main

import (
	"fmt"

	"github.com/hari-govind/liveserver-go/config"
	"github.com/hari-govind/liveserver-go/watcher"
)

func main() {
	fmt.Println(config.GetConfig().Depth)
	for ev := range watcher.Listen() {
		fmt.Println("Event: ", ev)
	}
}
