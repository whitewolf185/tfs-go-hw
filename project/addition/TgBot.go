package addition

const (
	TgTokenENV = "TG_BOT_TOKEN"
)

func TakeTgBotToken() string {

	// TgBot token parser
	token := ENVParser(TgTokenENV)

	return token
}
