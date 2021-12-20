package hendlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
	"github.com/whitewolf185/fs-go-hw/project/pkg/hendlers/addConn"
	"github.com/whitewolf185/fs-go-hw/project/pkg/tg-bot/addTGbot"
	"github.com/whitewolf185/fs-go-hw/project/repository/DB/addDB"
)

//go:generate mockgen -source=API.go -destination=mock_hendlers/mock.go
type ConnectionService interface {
	PrepareCandles(ctx context.Context, ws *websocket.Conn, wg *sync.WaitGroup,
		options addition.Options) (chan addConn.EventMsg, error)
	Unsubscribe(ws *websocket.Conn) error
	WebConn(url string) (*websocket.Conn, *http.Response, error)
	PingPong(ctx context.Context, wg *sync.WaitGroup, ws *websocket.Conn)
}

type API struct {
	urlWebSocket  string
	apiKeyPrivate string
	apiKeyPublic  string

	Ws       *websocket.Conn
	Ctx      context.Context
	cancel   context.CancelFunc
	connServ ConnectionService

	orChan    chan addition.Orders
	dBQueChan chan addDB.Query
	tgQueChan chan addDB.Query
}

func MakeAPI(ctx context.Context, service ConnectionService,
	orChan chan addition.Orders, dbChan, tgChan chan addDB.Query) *API {
	var api API
	tokens := addConn.TakeAPITokens()
	api.apiKeyPrivate = tokens.Private
	api.apiKeyPublic = tokens.Public
	api.urlWebSocket = tokens.URL

	api.connServ = service
	api.Ctx, api.cancel = context.WithCancel(ctx)
	api.orChan = orChan
	api.dBQueChan = dbChan
	api.tgQueChan = tgChan

	return &api
}

func (obj *API) WebsocketConnect(wg *sync.WaitGroup) error {
	ws, res, err := obj.connServ.WebConn(obj.urlWebSocket)
	if err != nil {
		if res.StatusCode == 404 {
			for i := 0; i < 4 && err != nil; i++ { // 4 попытки подключиться к вебсокету с интервалом в 5 секунд
				log.Println("trying to connect to WebSocket. Try", i+1)
				time.Sleep(time.Second * 5)
				ws, _, err = obj.connServ.WebConn(obj.urlWebSocket)
			}
		}
	}

	if err != nil {
		return MyErrors.WSConnectErr(err)
	}

	obj.Ws = ws

	log.Info("Connecting successful")

	obj.connServ.PingPong(obj.Ctx, wg, obj.Ws)
	return nil
}

// SendOrder формирует заказ на покупку валюты, чтобы затем отправить по каналу addition.Orders
func (obj *API) SendOrder(orderType addTGbot.OrdersTypes, ticket string, size int) {
	var order addition.Orders

	order.Endpoint = "/api/v3/sendorder"

	dataP := make(map[string]string)
	dataP["orderType"] = "mkt"
	dataP["symbol"] = ticket
	dataP["size"] = strconv.Itoa(size)

	switch orderType {
	case addTGbot.BuyOrder:
		dataP["side"] = "buy"

	case addTGbot.SellOrder:
		dataP["side"] = "sell"
	}

	order.PostData = dataP

	obj.orChan <- order
	log.Println("Order was sent")
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
				addConn.MakeQuery(&u, order.PostData)
				PostData := addConn.MakePostData(order.PostData)
				authent, err := addConn.GenerateAuthent(PostData, order.Endpoint, obj.apiKeyPrivate)

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

				var reqMsg addConn.ResponseMsg
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
					err = MyErrors.ErrOrderNotSuccess
				} else if reqMsg.Status.OrEvent[0].Type != "EXECUTION" {
					MyErrors.OrderSentErr("cannot do this operation with ticket right now")
					err = MyErrors.ErrStatusNotPlaced
				}

				// отправка заказа в БД
				if err == nil {
					var query addDB.Query
					query.Ticket = reqMsg.Status.OrEvent[0].Executed.Ticket
					query.LimitPrice = reqMsg.Status.OrEvent[0].Executed.LimitPrice
					query.Size = reqMsg.Status.OrEvent[0].Executed.Quantity
					query.Type = reqMsg.Status.OrEvent[0].Executed.Side
					obj.dBQueChan <- query
					obj.tgQueChan <- query
				}
			}
		}
	}()
}

func (obj *API) GetCandles(wg *sync.WaitGroup, optionChan chan addition.Options) (chan addConn.EventMsg, addition.Options, error) {
	log.Info("Waiting fot incoming options")
	option, ok := <-optionChan
	if !ok {
		return nil, option, MyErrors.ErrOptionChanErr
	}
	log.Info("Handler caught options")

	canChan, err := obj.connServ.PrepareCandles(obj.Ctx, obj.Ws, wg, option)
	for i := 1; err != nil && i < 11; i++ {
		log.Errorln(err)
		log.Info("Some err was caught. Waiting fot another incoming options... Try", i)
		option, ok = <-optionChan
		if !ok {
			return nil, option, MyErrors.ErrOptionChanErr
		}
		log.Info("Handler caught options")
		canChan, err = obj.connServ.PrepareCandles(obj.Ctx, obj.Ws, wg, option)
	}

	return canChan, option, err
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
