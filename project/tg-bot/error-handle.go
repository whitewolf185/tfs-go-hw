package tg_bot

import (
	log "github.com/sirupsen/logrus"
)

type MyErrors struct{}

func (obj MyErrors) TgBotTokensReadErr() {
	log.Fatalf("Cant see TgBot token ENV %s", TgTokenENV)
}

func (obj MyErrors) ReadFileErr(Type string, err error) {
	log.Fatalf("Problem with reading file %s.  Error: %s", Type, err.Error())
}

func (obj MyErrors) TgBotErr(err error) {
	log.Fatalf("Problem with tgbotapi.  Error: %s", err)
}

func (obj MyErrors) TgBotUpdateErr(err error) {
	log.Errorf("TgBot Update error.  Error: %s", err)
}
