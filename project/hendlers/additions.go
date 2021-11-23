package hendlers

import (
	"errors"
	"io/ioutil"
	"os"
)

type CandlePeriod string

const (
	CandlePeriod1m CandlePeriod = "1m"
	CandlePeriod2m CandlePeriod = "2m"
	CandlePeriod1h CandlePeriod = "1h"
)

var ErrUnknownPeriod = errors.New("unknown period")

// Options -- структура для записи необходимых для подписки на свечки настроек.
// Например, для того, чтобы настроить период свечей и какой тикет слушать.
type Options struct {
	ticket []string
	canPer CandlePeriod

	errHandler MyErrors
}

func CreateOptions(ticket []string, canPer string) (Options, error) {
	var (
		option Options
		err    error
	)
	option.ticket = ticket
	option.canPer, err = strToCanPer(canPer)
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
		return "", ErrUnknownPeriod
	}
}

// takeAPITokens функция, которая парсит ENV переменные.
func takeAPITokens(privateENV, publicENV, urlENV string) (string, string, string) {
	var (
		errHandler MyErrors
		ok         bool
		FilePath   string
		url        string
	)

	// private APIkey parsing
	FilePath, ok = os.LookupEnv(privateENV)
	if !ok {
		errHandler.APITokensReadErr("private")
	}

	apiKeyPrivate, err := ioutil.ReadFile(FilePath)
	if err != nil {
		errHandler.ReadFileErr(FilePath, err)
	}

	// public APIkey parsing
	FilePath, ok = os.LookupEnv(publicENV)
	if !ok {
		errHandler.APITokensReadErr("public")
	}

	apiKeyPublic, err := ioutil.ReadFile(FilePath)
	if err != nil {
		errHandler.ReadFileErr(FilePath, err)
	}

	// URL WB parsing
	url, ok = os.LookupEnv(urlENV)
	if !ok {
		errHandler.APITokensReadErr("WB url")
	}

	return string(apiKeyPrivate), string(apiKeyPublic), url
}
