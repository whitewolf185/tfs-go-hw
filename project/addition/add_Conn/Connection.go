package add_Conn

import (
	"os"
	"strconv"

	"github.com/whitewolf185/fs-go-hw/project/addition"
	"github.com/whitewolf185/fs-go-hw/project/addition/MyErrors"
)

type WSTokens struct {
	Private string
	Public  string
	Url     string
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
	result.Url, ok = os.LookupEnv(urlWebSocketENV)
	if !ok {
		MyErrors.TokensReadErr(urlWebSocketENV)
	}

	return result
}

func ConvertToFloat(price string) float32 {
	result, err := strconv.ParseFloat(price, 32)
	if err != nil {
		MyErrors.ConvertErr(err)
		return -1
	}

	return float32(result)
}
