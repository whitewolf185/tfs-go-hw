package addition

import (
	"fmt"
	"regexp"

	"main.go/project/addition/MyErrors"
)

const (
	TgTokenENV = "TG_BOT_TOKEN"
)

type MsgType int

type TakeProfitCh struct {
	TakeFl float32
	Size   int
}
type StopLossCh struct {
	StopFl float32
	Size   int
}

var commands = [...]string{
	"/buy",        // 0
	"/sell",       // 1
	"/option",     // 2
	"/start",      // 3
	"/stoploss",   // 4
	"/takeprofit", // 5
}

const (
	BuyNow     MsgType = iota // 0
	SellNow                   // 1
	OptionMsg                 // 2
	Start                     // 3
	StopLoss                  // 4
	TakeProfit                // 5
)

func TakeTgBotToken() string {

	// TgBot token parser
	token := ENVParser(TgTokenENV)

	return token
}

func MessageType(msg string) (MsgType, error) {

	for i := 0; i < len(commands); i++ {
		command := fmt.Sprintf("(?i)%s", commands[i])
		match, err := regexp.MatchString(command, msg)
		if err != nil {
			return -1, err
		}
		if match {
			return MsgType(i), nil
		}
	}

	return -1, MyErrors.NoMatches
}

func CreateOptions(ticket string, canPer string) (Options, error) {
	var (
		option Options
		err    error
	)
	tick, err := strToTicket(ticket)
	if err != nil {
		return Options{}, err
	}
	option.Ticket = tick
	option.CanPer, err = strToCanPer(canPer)
	if err != nil {
		return Options{}, err
	}

	return option, nil
}

func strToCanPer(period string) (CandlePeriod, error) {
	switch period {
	case "1m":
		return CandlePeriod1m, nil
	case "2m":
		return CandlePeriod2m, nil
	case "1h":
		return CandlePeriod1h, nil
	default:
		return "", MyErrors.ErrUnknownPeriod
	}
}

func strToTicket(ticket string) ([]string, error) {
	switch ticket {
	case "PI_XBTUSD":
		return []string{"PI_XBTUSD"}, nil

	case "PI_ETHUSD":
		return []string{"PI_ETHUSD"}, nil
	}

	return nil, MyErrors.ErrUnknownTicket
}
