package hendlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type connectionService interface {
	GetterCan
}

type API struct {
	urlWebSocket  string
	apiKeyPrivate string
	apiKeyPublic  string

	Ws       *websocket.Conn
	connServ connectionService
}

func MakeAPI(service connectionService) API {
	var api API
	privateTokenPathENV := "TOKEN_PATH_PRIVATE"
	publicTokenPathENV := "TOKEN_PATH_PUBLIC"
	urlWebSocketENV := "WS_URL"
	api.apiKeyPrivate, api.apiKeyPublic, api.urlWebSocket = takeAPITokens(privateTokenPathENV, publicTokenPathENV, urlWebSocketENV)
	api.connServ = service

	return api
}

func (obj API) generateAuthent(PostData, endpontPath string) (string, error) {
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

func (obj API) webConn() (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(obj.urlWebSocket, http.Header{
		"Sec-WebSocket-Extensions": []string{"permessage-deflate", "client_max_window_bits"}})
}

func (obj *API) WebsocketConnect() {
	var errHandler MyErrors
	ws, res, err := obj.webConn()
	if err != nil {
		if res.StatusCode == 404 {
			for i := 0; i < 4 && err != nil; i++ { // 4 попытки подключиться к вебсокету с интервалом в 5 секунд
				log.Println("trying to connect to WebSocket. Try", i+1)
				ticker := time.NewTicker(5 * time.Second)
				<-ticker.C
				ticker.Stop()
				ws, res, err = obj.webConn()
			}
		}
	}

	if err != nil {
		errHandler.WSConnectErr(err)
	}

	obj.Ws = ws

	log.Info("Connecting successful")
}

func (obj *API) Close() error {
	err := obj.Ws.WriteMessage(websocket.CloseMessage, []byte{})
	if err != nil {
		return err
	}
	err = obj.Ws.Close()
	if err != nil {
		return err
	}

	return nil
}
