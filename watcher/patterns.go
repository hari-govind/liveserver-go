package watcher

import (
	"github.com/hari-govind/liveserver-go/config"
	"log"
	"path"
	"path/filepath"
)

func ShouldWatchFile(filename string) bool {
	ignorePatterns := config.GetConfig().IgnorePatterns

	for _, pattern := range ignorePatterns {
		match, err := path.Match(pattern, filepath.Base(filename))
		if err != nil {
			log.Println("Error while comparing file patterns")
			log.Println(err)
		}
		if match {
			return false
		}
	}
	return true
}
