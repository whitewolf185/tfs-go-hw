package hendlers

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"main.go/project/MyErrors"
	"main.go/project/addition"
	tg_bot "main.go/project/tg-bot"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type connectionService interface {
	GetterCan
	Unsubscriber
}

type API struct {
	urlWebSocket  string
	apiKeyPrivate string
	apiKeyPublic  string

	Ws       *websocket.Conn
	Ctx      context.Context
	cancel   context.CancelFunc
	connServ connectionService
	orChan   chan tg_bot.Orders
}

func MakeAPI(service connectionService, ctx context.Context, orChan chan tg_bot.Orders) API {
	var api API
	tokens := addition.TakeAPITokens()
	api.apiKeyPrivate = tokens.Private
	api.apiKeyPublic = tokens.Public
	api.urlWebSocket = tokens.Url

	api.connServ = service
	api.Ctx, api.cancel = context.WithCancel(ctx)
	api.orChan = orChan

	return api
}

func (obj *API) generateAuthent(PostData, endpontPath string) (string, error) {
	// step 1 and 2
	sha := sha256.New()
	src := PostData + endpontPath
	sha.Write([]byte(src))

	// step 3
	apiDecode, err := base64.StdEncoding.DecodeString(obj.apiKeyPrivate)
	if err != nil {
		return "", err
	}

	// step 4
	h := hmac.New(sha512.New, apiDecode)
	h.Write(sha.Sum(nil))

	// step 5
	result := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return result, nil
}

func (obj *API) webConn() (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(obj.urlWebSocket, http.Header{
		"Sec-WebSocket-Extensions": []string{"permessage-deflate", "client_max_window_bits"}})
}

func (obj *API) WebsocketConnect(wg *sync.WaitGroup) {
	ws, res, err := obj.webConn()
	if err != nil {
		if res.StatusCode == 404 {
			for i := 0; i < 4 && err != nil; i++ { // 4 попытки подключиться к вебсокету с интервалом в 5 секунд
				log.Println("trying to connect to WebSocket. Try", i+1)
				time.Sleep(time.Second * 5)
				ws, res, err = obj.webConn()
			}
		}
	}

	if err != nil {
		MyErrors.WSConnectErr(err)
	}

	obj.Ws = ws

	log.Info("Connecting successful")

	obj.PingPong(wg)
}

// PingPong функция нужна для того, чтобы организовать heartbeat.
func (obj *API) PingPong(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		ping := time.NewTicker(pingPeriod)

		defer func() {
			ping.Stop()
			wg.Done()
		}()

		err := obj.Ws.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			// TODO тут надо делать реконнект
			panic(err)
		}
		obj.Ws.SetPongHandler(func(string) error {
			err := obj.Ws.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				// TODO тут надо делать реконнект
				panic(err)
			}
			return nil
		})

		for {
			select {
			case <-ping.C:
				if err := obj.Ws.WriteMessage(websocket.PingMessage, nil); err != nil {
					MyErrors.PingErr(err)
				}
			case <-obj.Ctx.Done():
				return
			}
		}
	}()

}

// OrderListener pipeline служащий для того, чтобы заниматься отправкой Orders по RESTAPI с целью покупки или продажи валюты
func (obj *API) OrderListener(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-obj.Ctx.Done():
				return

			case order := <-obj.orChan:
				log.Println("Got the order")
				u := url.URL{
					Scheme: "https",
					Host:   "demo-futures.kraken.com",
					Path:   "derivatives",
				}
				u.Path += order.Endpoint
				obj.makeQuery(&u, order.PostData)
				PostData := obj.makePostData(order.PostData)
				authent, err := obj.generateAuthent(PostData, order.Endpoint)
				fmt.Println(authent)
				if err != nil {
					MyErrors.APIGenerateErr(err)
				}

				client := &http.Client{}
				r, _ := http.NewRequest(http.MethodPost, u.String(), nil)
				r.Header.Add("APIKey", obj.apiKeyPublic)
				r.Header.Add("Authent", authent)

				req, err := client.Do(r)
				if err != nil {
					MyErrors.HTTPRequestErr(err)
				}

				var reqMsg ResponseMsg
				if err = json.NewDecoder(req.Body).Decode(&reqMsg); err != nil {
					panic(err)
				}
				err = req.Body.Close()
				if err != nil {
					MyErrors.BadBodyCloseErr(err)
				}

				err = nil
				if reqMsg.Result != "success" {
					MyErrors.OrderSentErr("request do not have result success")
					err = MyErrors.OrderNotSuccess
				} else if reqMsg.Status.Status != "placed" {
					MyErrors.OrderSentErr("cannot do this operation with ticket right now")
					err = MyErrors.StatusNotPlaced
				}
				if err != nil {
					// не записывать в БД
				}
				// записывать в БД
			}
		}
	}()
}

// TODO эту функцию можно протестировать
func (obj *API) makePostData(data map[string]string) string {
	values := url.Values{}

	for argument, value := range data {
		values.Add(argument, value)
	}

	return values.Encode()
}

func (obj *API) makeQuery(u *url.URL, data map[string]string) {
	q := u.Query()
	for argument, value := range data {
		q.Set(argument, value)
	}
	u.RawQuery = q.Encode()
}

func (obj *API) Close() error {
	obj.cancel()
	err := obj.connServ.Unsubscribe(obj.Ws)
	if err != nil {
		MyErrors.UnsubErr(err)
	}

	err = obj.Ws.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		return err
	}

	err = obj.Ws.Close()
	if err != nil {
		return err
	}

	return nil
}
