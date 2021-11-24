package tg_bot

import (
	"io/ioutil"
	"os"
)

const (
	TgTokenENV = "TG_BOT_TOKEN"
)

func takeTgBotToken() string {
	var (
		errHandler MyErrors
		ok         bool
		FilePath   string
	)

	// TgBot token parser
	FilePath, ok = os.LookupEnv(TgTokenENV)
	if !ok {
		errHandler.TgBotTokensReadErr()
	}

	token, err := ioutil.ReadFile(FilePath)
	if err != nil {
		errHandler.ReadFileErr(FilePath, err)
	}

	return string(token)
}
