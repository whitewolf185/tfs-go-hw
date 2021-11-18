package hendlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
)

type GetterCan interface {
	GetCandles(ws *websocket.Conn, ticket []string) (chan Candles, error)
}

type WSMsg struct {
	Event   string   `json:"event"`
	Feed    string   `json:"feed"`
	Tickets []string `json:"product_ids"`
}

type Connection struct {
	SubMessage WSMsg
	Candle     Candles
	errHandler MyErrors

	ws *websocket.Conn
}

func (obj *Connection) candleStream() (chan Candles, error) {
	canChan := make(chan Candles)

	go func() {
		for {
			_, data, err := obj.ws.ReadMessage()
			if err != nil {
				obj.errHandler.WSReadMsgErr(errors.New("In cansleStream" + err.Error()))
				close(canChan)
				return
			}
			var can Candles
			err = json.Unmarshal(data, &can)
			if err != nil {
				_ = obj.errHandler.UnmarshalErr(errors.New("In cansleStream" + err.Error()))
				close(canChan)
				return
			}
			canChan <- can
		}
	}()

	return canChan, nil
}

func (obj Connection) GetCandles(ws *websocket.Conn, ticket []string) (chan Candles, error) {
	obj.SubMessage.Event = "subscribe"
	obj.SubMessage.Feed = "candles_trade_1h"
	obj.SubMessage.Tickets = ticket
	obj.ws = ws

	msg, err := json.Marshal(obj.SubMessage)
	if err != nil {
		err = obj.errHandler.MarshalErr(err)
		return nil, err
	}

	fmt.Println(string(msg))

	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		return nil, err
	}

	_, data, err := ws.ReadMessage()
	if err != nil {
		obj.errHandler.WSReadMsgErr(err)
	}

	jsonData := make(map[string]interface{})

	_ = json.Unmarshal(data, &jsonData)

	fmt.Println(jsonData)

	return obj.candleStream()
}
