package MyErrors

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
)

var (
	ErrOrderNotSuccess = errors.New("request do not have result success")
	ErrStatusNotPlaced = errors.New("cannot do this operation with ticket right now")
	ErrUnknownPeriod   = errors.New("unknown period")
	ErrUnknownTicket   = errors.New("unknown Ticket")
	ErrNoMatches       = errors.New("no matches")
	ErrSub             = errors.New("subscribe failed")
	ErrOptionChanErr   = errors.New("options channel has closed")
)

func WSConnectErr(err error) error {
	str := fmt.Sprintf("Problem with WS connect.  Error: %s", err.Error())
	return errors.New(str)
}

func WSReadMsgErr(err error) error {
	tmp := fmt.Sprintf("Problem with WebSocket message read.  Error: %s", err.Error())
	return errors.New(tmp)
}

func WSWriteMsgErr(err error) {
	log.Fatalf("Problem with WebSocket message write.  Error: %s", err.Error())
}

func APIGenerateErr(err error) {
	log.Fatalf("Some problem with API generate func.  Error %s", err)
}

func TokensReadErr(tYpe string) {
	log.Fatalf("Cant see ENV value %s", tYpe)
}

func GetCandlesErr() {
	log.Fatalf("Something fatal was happend")
}

func ReadFileErr(tYpe string, err error) {
	log.Fatalf("Problem with reading file %s.  Error: %s", tYpe, err.Error())
}

func MarshalErr(err error) {
	log.Errorf("Cannot Marshal json. Error: %s", err.Error())
}

func UnmarshalErr(err error) error {
	tmp := fmt.Sprintf("Cannot Unarshal json. Error: %s", err.Error())
	return errors.New(tmp)
}

func PingErr(err error) {
	log.Errorf("Some troble with ping message.  Error: %s", err)
}

func UnsubErr(err error) {
	log.Errorf("Some troble with Unsubscribe.  Error: %s", err)
}

func BadAPIClose(err error) {
	log.Fatalf("bad Api close.  Error: %s", err)
}

func HTTPRequestErr(err error) {
	log.Errorf("Bad request.  Error: %s", err)
}

func OrderSentErr(tYpe string) {
	log.Errorf("Order sent failed because of %s", tYpe)
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
