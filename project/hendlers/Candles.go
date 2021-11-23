package hendlers

type SoloCandle struct {
	Time   int64  `json:"time"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume int    `json:"volume"`
}

type Candles struct {
	Candle SoloCandle `json:"candle"`
}
