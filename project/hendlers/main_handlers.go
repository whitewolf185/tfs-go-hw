package hendlers

import (
	"encoding/json"
	"fmt"
)

func Start() {
	api := MakeAPI()
	var errHandler MyErrors

	ws, _, err := api.WebsocketConnect()
	if err != nil {
		errHandler.WBConnectErr(err)
	}
	defer ws.Close()

	_, data, err := ws.ReadMessage()
	if err != nil {
		errHandler.WBReadMsgErr(err)
	}

	jsonData := make(map[string]interface{})

	_ = json.Unmarshal(data, &jsonData)

	fmt.Println(jsonData)
}
