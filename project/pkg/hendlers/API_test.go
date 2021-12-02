package hendlers

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition"
	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
	"github.com/whitewolf185/fs-go-hw/project/pkg/hendlers/mock_hendlers"
	"github.com/whitewolf185/fs-go-hw/project/repository/DB/addDB"
)

func SetUpENV() error {
	err := os.Setenv("TG_BOT_TOKEN", "D:/Documents/GO_projects/tfs-go-hw/project/cmd/config/tgBot_token.txt")
	if err != nil {
		return err
	}

	err = os.Setenv("WS_URL", "wss://demo-futures.kraken.com/ws/v1?chart")
	if err != nil {
		return err
	}

	err = os.Setenv("TOKEN_PATH_PUBLIC", "D:/Documents/GO_projects/tfs-go-hw/project/cmd/config/public-token.txt")
	if err != nil {
		return err
	}

	err = os.Setenv("TOKEN_PATH_PRIVATE", "D:/Documents/GO_projects/tfs-go-hw/project/cmd/config/private-token.txt")
	if err != nil {
		return err
	}

	return nil
}

func TestAPI_WebsocketConnect(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := SetUpENV()
	if err != nil {
		t.Fatalf("Unexpected error in ENV set up")
	}

	var (
		orChan chan addition.Orders
		dbChan chan addDB.Query
		tgChan chan addDB.Query
	)

	mockConnect := mock_hendlers.NewMockConnectionService(mockCtrl)
	testAPI := MakeAPI(context.Background(), mockConnect, orChan, dbChan, tgChan)
	wg := sync.WaitGroup{}

	// test 1
	t.Log("Test 1")
	ws := websocket.Conn{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mockConnect.EXPECT().WebConn(testAPI.urlWebSocket).Return(&websocket.Conn{}, &http.Response{StatusCode: 404}, errors.New("j")).Times(1)
	mockConnect.EXPECT().WebConn(testAPI.urlWebSocket).Return(&websocket.Conn{}, &http.Response{StatusCode: 200}, nil).Times(1)
	mockConnect.EXPECT().PingPong(ctx, &wg, &ws).Return().Times(1)
	got := testAPI.WebsocketConnect(&wg)
	if got != nil {
		t.Fatalf("wrong output expected %v got %e", nil, got)
	}

	// test 2
	t.Log("Test 2")
	mockConnect.EXPECT().WebConn(testAPI.urlWebSocket).Return(&websocket.Conn{}, &http.Response{StatusCode: 200}, nil).Times(1)
	mockConnect.EXPECT().PingPong(ctx, &wg, &ws).Return().Times(1)
	got = testAPI.WebsocketConnect(&wg)
	if got != nil {
		t.Fatalf("wrong output expected %v got %e", nil, got)
	}

	// test 3
	t.Log("Test 3")
	mockConnect.EXPECT().WebConn(testAPI.urlWebSocket).Return(&websocket.Conn{}, &http.Response{StatusCode: 404}, errors.New("j")).Times(5)

	got = testAPI.WebsocketConnect(&wg)
	if got.Error() != MyErrors.WSConnectErr(errors.New("j")).Error() {
		t.Fatalf("wrong output expected %e got %e", MyErrors.WSConnectErr(errors.New("j")), got)
	}
}

func SetOption(optionChan chan addition.Options, len int) {
	for i := 0; i < len; i++ {
		option := addition.Options{}
		optionChan <- option
	}
}

func TestAPI_GetCandles(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	err := SetUpENV()
	if err != nil {
		t.Fatalf("Unexpected error in ENV set up")
	}

	var (
		orChan chan addition.Orders
		dbChan chan addDB.Query
		tgChan chan addDB.Query
	)

	optionChan := make(chan addition.Options)

	mockConnect := mock_hendlers.NewMockConnectionService(mockCtrl)
	testAPI := MakeAPI(context.Background(), mockConnect, orChan, dbChan, tgChan)
	wg := sync.WaitGroup{}
	ws := websocket.Conn{}
	options := addition.Options{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	testAPI.Ws = &ws

	// test 1
	t.Log("Get candles test 1")
	mockConnect.EXPECT().PrepareCandles(ctx, &ws, &wg, options).Return(nil, nil).Times(1)
	go SetOption(optionChan, 1)

	_, err = testAPI.GetCandles(&wg, optionChan)
	if err != nil {
		t.Fatalf("Unexpected error on test 2")
	}

	// test 2
	t.Log("Get candles test 2")
	mockConnect.EXPECT().PrepareCandles(ctx, &ws, &wg, options).Return(nil, errors.New("gg")).Times(11)
	go SetOption(optionChan, 11)
	_, err = testAPI.GetCandles(&wg, optionChan)
	if err.Error() != errors.New("gg").Error() {
		t.Fatalf("Unexpected error on test 2")
	}

	// test 3
	t.Log("Get candles test 3")
	close(optionChan)
	_, err = testAPI.GetCandles(&wg, optionChan)
	if err.Error() != MyErrors.ErrOptionChanErr.Error() {
		t.Fatalf("Unexpected error on test 2")
	}
}
