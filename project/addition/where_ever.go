package addition

import (
	"io/ioutil"
	"main.go/project/MyErrors"
	"os"
)

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
