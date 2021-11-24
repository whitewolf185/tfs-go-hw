package hendlers

type SoloCandle struct {
	Time   int64  `json:"time"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume int    `json:"volume"`
}

type EventMsg struct {
	Candle    SoloCandle `json:"candle"`
	Result    string     `json:"result,omitempty"`
	ProductId string     `json:"product_id,omitempty"`
}

type SendStatus struct {
	Status string `json:"status"`
}

type ResponseMsg struct {
	Result string     `json:"result"`
	Status SendStatus `json:"sendStatus"`
}
