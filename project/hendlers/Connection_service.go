package hendlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
)

type connectionService interface {
	GetCandles(ws *websocket.Conn, ticket []string) error
}

type WSMsg struct {
	Event   string   `json:"event"`
	Feed    string   `json:"feed"`
	Tickets []string `json:"product_ids"`
}

type Connection struct {
	Message WSMsg

	errHandler MyErrors
}

func (obj Connection) GetCandles(ws *websocket.Conn, ticket []string) error {
	obj.Message.Event = "subscribe"
	obj.Message.Feed = "candles_trade_1h"
	obj.Message.Tickets = ticket

	msg, err := json.Marshal(obj.Message)
	if err != nil {
		err = obj.errHandler.MarshalErr(err)
		return err
	}

	fmt.Println(string(msg))

	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return err
	}

	return nil
}
