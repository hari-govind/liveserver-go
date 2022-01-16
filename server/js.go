package server

import (
	"github.com/hari-govind/liveserver-go/config"

	"fmt"
	"strconv"
)

// Will inject the following script tag to html files when serving
// Script will reload the page when any data is received through websocket
// Will wait  config.Wait ms before reloading, new data resets the timeout
var SCRIPT_TAG = []byte(`
<script>
(function websocket(){
	let ws = new WebSocket('ws://` + fmt.Sprintf("%s:%d", config.GetConfig().ListenAddress, config.GetConfig().WebsocketPort) + `');
	let tid;
	ws.onmessage = function(event) {
		console.log(event.data);
		clearTimeout(tid);
		tid = setTimeout(function(){
		location.reload();
	},` + strconv.Itoa(config.GetConfig().Wait) + `)

	};

	ws.onerror = function(err) {
		if (window.confirm("Websocket connection failed, reload page?")) {
				location.reload();
		}
	}

}())
</script>
`)
