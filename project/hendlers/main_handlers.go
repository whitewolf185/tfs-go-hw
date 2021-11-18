package hendlers

import (
	"encoding/json"
	"fmt"
)

func Start() {
	api := MakeAPI(Connection{})
	api.WebsocketConnect()
	defer api.Close()
	var errHandler MyErrors

	_, data, err := api.Ws.ReadMessage()
	if err != nil {
		errHandler.WSReadMsgErr(err)
	}

	jsonData := make(map[string]interface{})

	_ = json.Unmarshal(data, &jsonData)

	fmt.Println(jsonData)

	canChan, err := api.connServ.GetCandles(api.Ws, []string{"PI_XBTUSD"})
	if err != nil {
		errHandler.GetCandlesErr(err)
	}

	for {
		can := <-canChan

		fmt.Println(can)
	}
}
