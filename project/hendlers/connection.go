package hendlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"github.com/gorilla/websocket"
	"net/http"
)

type API struct {
	urlWebSocket  string
	apiKeyPrivate string
	apiKeyPublic  string
}

func MakeAPI() API {
	var api API
	privateTokenPathENV := "TOKEN_PATH_PRIVATE"
	publicTokenPathENV := "TOKEN_PATH_PUBLIC"
	urlWebSocketENV := "WS_URL"
	api.apiKeyPrivate, api.apiKeyPublic, api.urlWebSocket = takeAPITokens(privateTokenPathENV, publicTokenPathENV, urlWebSocketENV)

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

func (obj API) WebsocketConnect() (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(obj.urlWebSocket, http.Header{
		"Sec-WebSocket-Extensions": []string{"permessage-deflate", "client_max_window_bits"}})
}
