package MyErrors

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

var (
	OrderNotSuccess  = errors.New("request do not have result success")
	StatusNotPlaced  = errors.New("cannot do this operation with ticket right now")
	ErrUnknownPeriod = errors.New("unknown period")
	ErrUnknownTicket = errors.New("unknown Ticket")
	NoMatches        = errors.New("No matches\n")
)

func WSConnectErr(err error) {
	log.Panicf("Problem with WS connect.  Error: %s", err.Error())
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

func ReadFileErr(Type string, err error) {
	log.Fatalf("Problem with reading file %s.  Error: %s", Type, err.Error())
}

func MarshalErr(err error) {
	log.Errorf("Cannot Marshal json. Error: %s", err.Error())
}

func UnmarshalErr(err error) {
	log.Errorf("Cannot Unarshal json. Error: %s", err.Error())
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

func TgBotMsgErr(err error) {
	log.Errorf("TgBot Message error.  Error: %s", err)
}

func TgBotErr(err error) {
	log.Fatalf("Problem with tgbotapi.  Error: %s", err)
}

func TgBotUpdateErr(err error) {
	log.Errorf("TgBot Update error.  Error: %s", err)
}

func DBConnectionErr(err error) {
	log.Fatalf("Cannot connect to DB.  Error: %s", err)
}

func DBCloseConnErr(err error) {
	log.Fatalf("Cannot close DB connection.  Error: %s", err)
}

func DBExecErr(err error) {
	log.Errorf("Data base exec error.  Error: %s", err)
}

func RegexpErr(err error) {
	log.Errorf("Something wrong with RegExp.  Error: %s", err)
}

func SendMsgErr(err error) {
	log.Errorf("Message has not sent.  Error: %s", err)
}

func ConvertErr(err error) {
	log.Errorf("Something wrong with convert.  Error: %s", err)
}
