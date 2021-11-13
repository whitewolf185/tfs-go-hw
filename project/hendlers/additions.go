package hendlers

import (
	"io/ioutil"
	"os"
)

func takeAPITokens(privateENV, publicENV, urlENV string) (string, string, string) {
	var (
		errHandler MyErrors
		ok         bool
		FilePath   string
		url        string
	)

	// private APIkey parsing
	FilePath, ok = os.LookupEnv(privateENV)
	if !ok {
		errHandler.APITokensReadErr("private")
	}

	apiKeyPrivate, err := ioutil.ReadFile(FilePath)
	if err != nil {
		errHandler.ReadFileErr(FilePath, err)
	}

	// public APIkey parsing
	FilePath, ok = os.LookupEnv(publicENV)
	if !ok {
		errHandler.APITokensReadErr("public")
	}

	apiKeyPublic, err := ioutil.ReadFile(FilePath)
	if err != nil {
		errHandler.ReadFileErr(FilePath, err)
	}

	// URL WB parsing
	url, ok = os.LookupEnv(urlENV)
	if !ok {
		errHandler.APITokensReadErr("WB url")
	}

	return string(apiKeyPrivate), string(apiKeyPublic), url
}
