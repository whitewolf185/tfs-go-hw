package tg_bot

// Orders структура, используемая для того, чтобы отправлять запросы на покупку или продажу валюту
type Orders struct {
	Endpoint string
	PostData map[string]string
}

type OrdersTypes int

const (
	BuyOrder OrdersTypes = iota
)

type Messages struct {
}
