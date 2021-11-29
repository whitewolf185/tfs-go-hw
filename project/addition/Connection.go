package addition

import (
	"os"
	"strconv"

	"main.go/project/addition/MyErrors"
)

type WSTokens struct {
	Private string
	Public  string
	Url     string
}

type CandlePeriod string

const (
	CandlePeriod1m CandlePeriod = "1m"
	CandlePeriod2m CandlePeriod = "2m"
	CandlePeriod1h CandlePeriod = "1h"
)

const (
	privateTokenPathENV = "TOKEN_PATH_PRIVATE"
	publicTokenPathENV  = "TOKEN_PATH_PUBLIC"
	urlWebSocketENV     = "WS_URL"
)

// Options -- структура для записи необходимых для подписки на свечки настроек.
// Например, для того, чтобы настроить период свечей и какой тикет слушать.
type Options struct {
	Ticket []string
	CanPer CandlePeriod
}

// TakeAPITokens функция, которая выдает API токены.
func TakeAPITokens() WSTokens {
	var (
		result WSTokens
		ok     bool
	)

	// APIkey parsing
	result.Private = ENVParser(privateTokenPathENV)
	result.Public = ENVParser(publicTokenPathENV)

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
	}

	return float32(result)
}
