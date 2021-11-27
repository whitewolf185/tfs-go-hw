package addition

import (
	"io/ioutil"
	"main.go/project/MyErrors"
	"os"
)

// Orders структура, используемая для того, чтобы отправлять запросы на покупку или продажу валюту
type Orders struct {
	Endpoint string
	PostData map[string]string
}

func ENVParser(ENV string) string {
	FilePath, ok := os.LookupEnv(ENV)
	if !ok {
		MyErrors.TokensReadErr(ENV)
	}

	token, err := ioutil.ReadFile(FilePath)
	if err != nil {
		MyErrors.ReadFileErr(FilePath, err)
	}

	return string(token)
}
