package hendlers

import (
	"errors"
	log "github.com/sirupsen/logrus"
)

type MyErrors struct{}

func (obj MyErrors) WSConnectErr(err error) {
	log.Fatalf("Problem with WS connect.  Error: %s", err.Error())
}

func (obj MyErrors) WSReadMsgErr(err error) {
	log.Errorf("Problem with WebSocket message read.  Error: %s", err.Error())
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
