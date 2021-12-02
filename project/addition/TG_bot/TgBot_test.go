package TG_bot

import (
	"github.com/whitewolf185/fs-go-hw/project/addition/MyErrors"
	"testing"
)

type Result struct {
	MsgType MsgType
	err     error
}

type TestMessages struct {
	Name     string
	In       string
	Expected Result
}

func TestMessageType(t *testing.T) {
	tests := []TestMessages{
		{
			"OK",
			"/buy",
			Result{
				BuyNow,
				nil,
			},
		},
		{
			"OK",
			"/sell",
			Result{
				SellNow,
				nil,
			},
		},
		{
			"OK",
			"/option",
			Result{
				OptionMsg,
				nil,
			},
		},
		{
			"OK",
			"/start",
			Result{
				Start,
				nil,
			},
		},
		{
			"OK",
			"/stoploss",
			Result{
				StopLoss,
				nil,
			},
		},
		{
			"OK",
			"/takeprofit",
			Result{
				TakeProfit,
				nil,
			},
		},

		{
			"NoMatch",
			"/takeprifit",
			Result{
				-1,
				MyErrors.NoMatches,
			},
		},
		{
			"OK",
			"/takeprofit^jD&Sf dad /\n",
			Result{
				TakeProfit,
				nil,
			},
		},
		{
			"NoMatch",
			"asdgadtakeprofit^jD&Sf dad /\n",
			Result{
				-1,
				MyErrors.NoMatches,
			},
		},
	}

	for idx, test := range tests {
		got, err := MessageType(test.In)

		if got != test.Expected.MsgType || err != test.Expected.err {
			t.Errorf("Test %d expected %v got %v and %e", idx, test.Expected, got, err)
		}
	}
}
