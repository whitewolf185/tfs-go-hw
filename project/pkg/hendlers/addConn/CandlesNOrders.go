package addConn

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
	ProductID string     `json:"product_id,omitempty"`
}

type OrderExecution struct {
	Ticket     string  `json:"symbol"`
	LimitPrice float32 `json:"limitPrice"`
	Quantity   int     `json:"quantity"`
	Side       string  `json:"side"`
}

type OrVent struct {
	Type     string         `json:"type"`
	Executed OrderExecution `json:"orderPriorExecution"`
}

type SendStatus struct {
	OrEvent []OrVent `json:"orderEvents"`
}

type ResponseMsg struct {
	Result string     `json:"result"`
	Status SendStatus `json:"sendStatus"`
}
