package hendlers

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

var SubErr = errors.New("subscribe error")

type MyErrors struct{}

func (obj MyErrors) WSConnectErr(err error) {
	log.Fatalf("Problem with WS connect.  Error: %s", err.Error())
}

func (obj MyErrors) WSReadMsgErr(err error) {
	log.Panicf("Problem with WebSocket message read.  Error: %s", err.Error())
}

func (obj MyErrors) APITokensReadErr(Type string) {
	log.Fatalf("Problem with %s API read", Type)
}

func (obj MyErrors) GetCandlesErr(err error) {
	log.Errorf("Problem with GetCandles.  Error: %s", err.Error())
}

func (obj MyErrors) ReadFileErr(Type string, err error) {
	log.Fatalf("Problem with reading file in %s.  Error: %s", Type, err.Error())
}

func (obj MyErrors) MarshalErr(err error) error {
	log.Errorf("Cannot Marshal json. Error: %s", err.Error())
	return errors.New("See error above\n")
}

func (obj MyErrors) UnmarshalErr(err error) error {
	log.Errorf("Cannot Unarshal json. Error: %s", err.Error())
	return errors.New("See error above\n")
}

func (obj MyErrors) UnknownPeriod(err error) {
	log.Errorln(err)
}

func (obj MyErrors) SubErr() {
	log.Errorln("subscribe failed")
}

func (obj MyErrors) PingErr(err error) {
	log.Errorf("Some troble with ping message.  Error: %s", err)
}

func (obj MyErrors) UnsubErr(err error) {
	log.Errorf("Some troble with Unsubscribe.  Error: %s", err)
}

func (obj MyErrors) BadApiClose(err error) {
	log.Fatalf("bad Api close.  Error: %s", err)
}
