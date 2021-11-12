package hendlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"

	"github.com/gorilla/websocket"
)

type APIService interface {
	generateAuthent(PostData, endpontPath string)
	WebsocketConnect() (*websocket.Conn, error)
	takeAPITokens()
}

type API struct {
	privateTokenPathENV string
	publicTokenPathENV  string
	urlWebSocket        string
	apiKeyPrivate       string
	apiKeyPublic        string
	service             APIService
}

func MakeAPI() API {
	var api API
	api.privateTokenPathENV = "TOKEN_PATH_PRIVATE"
	api.publicTokenPathENV = "TOKEN_PATH_PUBLIC"
	api.urlWebSocket = "WS_URL"
	api.service.takeAPITokens()

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

func (obj API) WebsocketConnect() (*websocket.Conn, error) {

}
