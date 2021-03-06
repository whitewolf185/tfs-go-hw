package addConn

import (
	"os"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
)

type WSTokens struct {
	Private string
	Public  string
	URL     string
}

const (
	privateTokenPathENV = "TOKEN_PATH_PRIVATE"
	publicTokenPathENV  = "TOKEN_PATH_PUBLIC"
	urlWebSocketENV     = "WS_URL"
)

// TakeAPITokens функция, которая выдает API токены.
func TakeAPITokens() WSTokens {
	var (
		result WSTokens
		ok     bool
	)

	// APIkey parsing
	result.Private = addition.ENVParser(privateTokenPathENV)
	result.Public = addition.ENVParser(publicTokenPathENV)

	// URL WB parsing
	result.URL, ok = os.LookupEnv(urlWebSocketENV)
	if !ok {
		MyErrors.TokensReadErr(urlWebSocketENV)
	}

	return result
}
