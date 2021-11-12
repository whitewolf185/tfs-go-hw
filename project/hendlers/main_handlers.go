package hendlers

import (
	"github.com/gorilla/websocket"
)

func Start() {
	api := MakeAPI()
	var errHandler MyErrors

	ws, err := api.service.WebsocketConnect()
	if err != nil {
		errHandler.handler.WBConnectErr(err)
	}

	_, data, err := ws.ReadMessage()
	if err != nil {
		return
	}
}
