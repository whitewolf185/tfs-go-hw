package addition

import (
	"io/ioutil"
	"os"

	"github.com/whitewolf185/fs-go-hw/project/addition/MyErrors"
)

type TakeProfitCh struct {
	TakeFl float32
	Size   int
}
type StopLossCh struct {
	StopFl float32
	Size   int
}

type CandlePeriod string

const (
	CandlePeriod1m CandlePeriod = "1m"
	CandlePeriod2m CandlePeriod = "2m"
	CandlePeriod1h CandlePeriod = "1h"
)

// Options -- структура для записи необходимых для подписки на свечки настроек.
// Например, для того, чтобы настроить период свечей и какой тикет слушать.
type Options struct {
	Ticket []string
	CanPer CandlePeriod
}

// Orders структура, используемая для того, чтобы отправлять запросы на покупку или продажу валюту
type Orders struct {
	Endpoint string
	PostData map[string]string
}

func ENVParser(ENV string) string {
	FilePath, ok := os.LookupEnv(ENV)
	if !ok {
		MyErrors.TokensReadErr(ENV)
	}

	token, err := ioutil.ReadFile(FilePath)
	if err != nil {
		MyErrors.ReadFileErr(FilePath, err)
	}

	return string(token)
}
