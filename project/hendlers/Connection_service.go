package hendlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"main.go/project/MyErrors"
	"main.go/project/addition"
	"sync"
	"time"
)

const (
	writeWait  = 2 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type GetterCan interface {
	GetCandles(ws *websocket.Conn, wg *sync.WaitGroup, ctx context.Context, options addition.Options) chan EventMsg
}

type Unsubscriber interface {
	Unsubscribe(ws *websocket.Conn) error
}

type WSMsg struct {
	Event   string   `json:"event"`
	Feed    string   `json:"feed"`
	Tickets []string `json:"product_ids"`
}

type Connection struct {
	SubMessage WSMsg
	Candle     EventMsg

	ws *websocket.Conn
}

func (obj *Connection) candleStream(wg *sync.WaitGroup, ctx context.Context) chan EventMsg {
	canChan := make(chan EventMsg)

	wg.Add(1)
	go func() {
		prev := make(map[string]int64)
		for _, ticket := range obj.SubMessage.Tickets {
			prev[ticket] = 0
		}

		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				close(canChan)
				return
			default:
				var event EventMsg
				err := obj.ws.ReadJSON(&event)
				if err != nil {
					MyErrors.WSReadMsgErr(errors.New("In cansleStream" + err.Error()))
					close(canChan)
					return
				}
				if event.Result == "error" { // обработчик ошибки отправки order
					log.Errorln(event)
				} else if prev[event.ProductId] != event.Candle.Time && event.Candle.Volume != 0 {
					canChan <- event
					prev[event.ProductId] = event.Candle.Time
				}
			}
		}
	}()

	return canChan
}

func (obj Connection) GetCandles(ws *websocket.Conn, wg *sync.WaitGroup, ctx context.Context, options addition.Options) chan EventMsg {
	obj.SubMessage.Event = "subscribe"
	obj.SubMessage.Feed = "candles_trade_" + string(options.CanPer)
	obj.SubMessage.Tickets = options.Ticket
	obj.ws = ws

	msg, err := json.Marshal(obj.SubMessage)
	if err != nil {
		MyErrors.MarshalErr(err)
	}

	fmt.Println(string(msg))

	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		MyErrors.WSWriteMsgErr(err)
	}

	_, data, err := ws.ReadMessage()
	if err != nil {
		MyErrors.WSReadMsgErr(err)
	}

	jsonData := make(map[string]interface{})

	_ = json.Unmarshal(data, &jsonData)

	fmt.Println(jsonData)

	if jsonData["event"] != "subscribed" {
		MyErrors.SubErr()
	}

	return obj.candleStream(wg, ctx)
}

func (obj Connection) Unsubscribe(ws *websocket.Conn) error {
	obj.SubMessage.Event = "unsubscribe"
	err := ws.WriteJSON(obj.SubMessage)
	if err != nil {
		return err
	}

	time.Sleep(writeWait)

	return nil
}
