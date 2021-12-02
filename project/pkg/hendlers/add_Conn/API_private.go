package add_Conn

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"net/url"
)

// GenerateAuthent функция, которая генерирует токен под названием authent
func GenerateAuthent(PostData, endpontPath, apiKeyPrivate string) (string, error) {
	// step 1 and 2
	sha := sha256.New()
	src := PostData + endpontPath
	sha.Write([]byte(src))

	// step 3
	apiDecode, err := base64.StdEncoding.DecodeString(apiKeyPrivate)
	if err != nil {
		return "", err
	}

	// step 4
	h := hmac.New(sha512.New, apiDecode)
	h.Write(sha.Sum(nil))

	// step 5
	result := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return result, nil
}

func MakePostData(data map[string]string) string {
	values := url.Values{}

	for argument, value := range data {
		values.Add(argument, value)
	}

	return values.Encode()
}

func MakeQuery(u *url.URL, data map[string]string) {
	q := u.Query()
	for argument, value := range data {
		q.Set(argument, value)
	}
	u.RawQuery = q.Encode()
}
