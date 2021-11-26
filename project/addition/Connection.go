package addition

import (
	"main.go/project/MyErrors"
	"os"
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

func CreateOptions(ticket []string, canPer string) (Options, error) {
	var (
		option Options
		err    error
	)
	option.Ticket = ticket
	option.CanPer, err = strToCanPer(canPer)
	if err != nil {
		return Options{}, err
	}

	return option, nil
}

func strToCanPer(period string) (CandlePeriod, error) {
	switch period {
	case "1m":
		return CandlePeriod1m, nil
	case "2m":
		return CandlePeriod2m, nil
	case "1h":
		return CandlePeriod1h, nil
	default:
		return "", MyErrors.ErrUnknownPeriod
	}
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
