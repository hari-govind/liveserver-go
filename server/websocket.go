package server

import (
	"github.com/hari-govind/liveserver-go/config"
	"github.com/hari-govind/liveserver-go/watcher"

	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WsHandler struct {
	cond *sync.Cond
}

func (ws *WsHandler) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println("Cannot create websocket ", err)
	}

	ws.cond.L.Lock()
	for {
		ws.cond.Wait()
		err := conn.WriteMessage(websocket.TextMessage, []byte("file changed"))
		if err != nil {
			// Close websocket connection when messages are unreachable
			// Each page reload will result in the previous one being unused
			// These connections will remain open until next change is attempted to be written
			// log.Println("Cannot write to websocket,", err)
			conn.Close()
			ws.cond.L.Unlock()
			break
		}
	}
}

func StartWsServer() {
	events := watcher.Listen()
	m := sync.Mutex{}
	cond := sync.NewCond(&m)
	wsHandler := &WsHandler{cond}
	http.HandleFunc("/", wsHandler.serveWs)

	go func() {
		for range events {
			cond.Broadcast()
		}
	}()

	err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.GetConfig().ListenAddress, config.GetConfig().WebsocketPort), nil)
	if err != nil {
		log.Panicln("Cannot start websocket server", err)
	}
}
