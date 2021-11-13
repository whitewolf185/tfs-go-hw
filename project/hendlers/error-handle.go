package hendlers

import (
	log "github.com/sirupsen/logrus"
)

type MyErrors struct{}

func (obj MyErrors) WBConnectErr(err error) {
	log.Fatalf("Problem with WS connect.  Error: %s", err.Error())
}

func (obj MyErrors) WBReadMsgErr(err error) {
	log.Errorf("Problem with WebSocket message read.  Error: %s", err.Error())
}

func (obj MyErrors) APITokensReadErr(Type string) {
	log.Fatalf("Problem with %s API read", Type)
}

func (obj MyErrors) ReadFileErr(Type string, err error) {
	log.Fatalf("Problem with reading file in %s.  Error: %s", Type, err.Error())
}
