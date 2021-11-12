package hendlers

import (
	log "github.com/sirupsen/logrus"
)

type errorer interface {
	WBConnectErr(error)
}

type MyErrors struct {
	handler errorer
}

func (obj MyErrors) WBConnectErr(err error) {
	log.Errorf("Problem with WS connect:\n%e", err)
}
