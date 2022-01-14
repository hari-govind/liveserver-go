package watcher

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type Watcher interface {
	WatchPath(path string)
	ListenForChanges() <-chan string
}

type FileWatcher struct {
	watcher *fsnotify.Watcher
	events  chan string
}

func NewFileWatcher() *FileWatcher {
	fw := new(FileWatcher)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Cannot get watcher: ", err)
	}
	fw.watcher = watcher
	return fw
}

func (fw *FileWatcher) WatchPath(path string, depth int) {
	if ShouldWatchFile(path) {
		fw.watcher.Add(path)
	}

	if depth > 1 {
		dirs, err := ioutil.ReadDir(path)
		if err != nil {
			log.Print("Cannot read directory", err)
		}
		for _, dir := range dirs {
			if dir.IsDir() {
				fw.WatchPath(filepath.Join(path, dir.Name()), depth-1) //recursively watch subdirs
			}
		}
	}
}

func (fw *FileWatcher) ListenForChanges() <-chan string {
	if fw.events == nil {
		fw.events = make(chan string, 1024)
	}
	go func() {
		for {
			select {
			case event, ok := <-fw.watcher.Events:
				if !ok {
					log.Fatal("Cannot get event: ", event)
				}
				if ShouldWatchFile(event.Name) && (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) {
					fw.events <- event.Name
				}
			case err, ok := <-fw.watcher.Errors:
				if !ok {
					log.Fatal("Cannot get event: ")
				}
				log.Fatal(err)
			}
		}
	}()
	return fw.events
}
