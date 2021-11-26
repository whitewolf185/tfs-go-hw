package MyErrors

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

var (
	OrderNotSuccess  = errors.New("request do not have result success")
	StatusNotPlaced  = errors.New("cannot do this operation with ticket right now")
	ErrUnknownPeriod = errors.New("unknown period")
)

func WSConnectErr(err error) {
	log.Fatalf("Problem with WS connect.  Error: %s", err.Error())
}

func WSReadMsgErr(err error) {
	log.Fatalf("Problem with WebSocket message read.  Error: %s", err.Error())
}

func WSWriteMsgErr(err error) {
	log.Fatalf("Problem with WebSocket message write.  Error: %s", err.Error())
}

func APIGenerateErr(err error) {
	log.Fatalf("Some problem with API generate func.  Error %s", err)
}

func TokensReadErr(Type string) {
	log.Fatalf("Cant see ENV value %s", Type)
}

func GetCandlesErr(err error) {
	log.Errorf("Problem with GetCandles.  Error: %s", err.Error())
}

func ReadFileErr(Type string, err error) {
	log.Fatalf("Problem with reading file %s.  Error: %s", Type, err.Error())
}

func MarshalErr(err error) {
	log.Errorf("Cannot Marshal json. Error: %s", err.Error())
}

func UnmarshalErr(err error) {
	log.Errorf("Cannot Unarshal json. Error: %s", err.Error())
}

func UnknownPeriod(err error) {
	log.Errorln(err)
}

func SubErr() {
	log.Errorln("subscribe failed")
}

func PingErr(err error) {
	log.Errorf("Some troble with ping message.  Error: %s", err)
}

func UnsubErr(err error) {
	log.Errorf("Some troble with Unsubscribe.  Error: %s", err)
}

func BadApiClose(err error) {
	log.Fatalf("bad Api close.  Error: %s", err)
}

func HTTPRequestErr(err error) {
	log.Errorf("Bad request.  Error: %s", err)
}

func OrderSentErr(Type string) {
	log.Errorf("Order sent failed because of %s", Type)
}

func BadBodyCloseErr(err error) {
	log.Errorf("bad close request body.  Errror: %s", err)
}

func TgBotErr(err error) {
	log.Fatalf("Problem with tgbotapi.  Error: %s", err)
}

func TgBotUpdateErr(err error) {
	log.Errorf("TgBot Update error.  Error: %s", err)
}
