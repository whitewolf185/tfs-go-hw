package hendlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/whitewolf185/fs-go-hw/project/pkg/hendlers/add_Conn"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
)

const (
	writeWait  = 2 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type SubEvent struct {
	Event string `json:"event"`
}

type WSMsg struct {
	Event   string   `json:"event"`
	Feed    string   `json:"feed"`
	Tickets []string `json:"product_ids"`
}

type Connection struct {
	SubMessage WSMsg
	Candle     add_Conn.EventMsg

	ws *websocket.Conn
}

// candleStream -- функция-обработчик отправки Orders
func (obj *Connection) candleStream(wg *sync.WaitGroup, ctx context.Context) (chan add_Conn.EventMsg, error) {
	canChan := make(chan add_Conn.EventMsg)

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
				var event add_Conn.EventMsg
				err := obj.ws.ReadJSON(&event)
				if err != nil {
					_ = MyErrors.WSReadMsgErr(errors.New("In cansleStream" + err.Error()))
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

	return canChan, nil
}

func (obj Connection) PrepareCandles(ws *websocket.Conn, wg *sync.WaitGroup, ctx context.Context, options addition.Options) (chan add_Conn.EventMsg, error) {
	obj.SubMessage.Event = "subscribe"
	obj.SubMessage.Feed = "candles_trade_" + string(options.CanPer)
	obj.SubMessage.Tickets = options.Ticket
	obj.ws = ws

	msg, err := json.Marshal(obj.SubMessage)
	if err != nil {
		MyErrors.MarshalErr(err)
		return nil, err
	}

	fmt.Println(string(msg))

	err = ws.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		MyErrors.WSWriteMsgErr(err)
		return nil, err
	}

	_, data, err := ws.ReadMessage()
	if err != nil {
		return nil, MyErrors.WSReadMsgErr(err)
	}

	jsonData := SubEvent{}
	if err = json.Unmarshal(data, &jsonData); err != nil {
		return nil, MyErrors.UnmarshalErr(err)
	}
	fmt.Println(jsonData)
	if jsonData.Event != "subscribed" {
		return nil, MyErrors.SubErr
	}

	return obj.candleStream(wg, ctx)
}

func (obj Connection) WebConn(url string) (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(url, http.Header{
		"Sec-WebSocket-Extensions": []string{"permessage-deflate", "client_max_window_bits"}})
}

// PingPong функция нужна для того, чтобы организовать heartbeat.
func (obj Connection) PingPong(wg *sync.WaitGroup, ctx context.Context, ws *websocket.Conn) {
	wg.Add(1)
	go func() {
		ping := time.NewTicker(pingPeriod)

		defer func() {
			ping.Stop()
			wg.Done()
		}()

		err := ws.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			log.Info("Reconnecting to WS...")
			return
		}
		ws.SetPongHandler(func(string) error {
			err := ws.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				log.Info("Reconnecting to WS...")
				return err
			}
			return nil
		})

		for {
			select {
			case <-ping.C:
				if err := ws.WriteMessage(websocket.PingMessage, nil); err != nil {
					MyErrors.PingErr(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
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
