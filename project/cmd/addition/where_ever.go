package addition

import (
	"io/ioutil"
	"os"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
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

func ENVParser(env string) string {
	FilePath, ok := os.LookupEnv(env)
	if !ok {
		MyErrors.TokensReadErr(env)
	}

	token, err := ioutil.ReadFile(FilePath)
	if err != nil {
		MyErrors.ReadFileErr(FilePath, err)
	}

	return string(token)
}
