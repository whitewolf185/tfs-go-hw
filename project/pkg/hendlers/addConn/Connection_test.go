package addConn

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetUpENV() error {
	err := os.Setenv("TG_BOT_TOKEN", "D:/Documents/GO_projects/tfs-go-hw/project/cmd/config/tgBot_token.txt")
	if err != nil {
		return err
	}

	err = os.Setenv("WS_URL", "wss://demo-futures.kraken.com/ws/v1?chart")
	if err != nil {
		return err
	}

	err = os.Setenv("TOKEN_PATH_PUBLIC", "D:/Documents/GO_projects/tfs-go-hw/project/cmd/config/public-token_test.txt")
	if err != nil {
		return err
	}

	err = os.Setenv("TOKEN_PATH_PRIVATE", "D:/Documents/GO_projects/tfs-go-hw/project/cmd/config/private-token_test.txt")
	if err != nil {
		return err
	}

	return nil
}

func TestTakeAPITokens(t *testing.T) {
	err := SetUpENV()
	if err != nil {
		t.Fatalf("Unexpected error")
	}

	expected := WSTokens{
		"asdfjkl;",
		"123456879",
		"wss://demo-futures.kraken.com/ws/v1?chart",
	}

	got := TakeAPITokens()
	assert.Equal(t, expected, got)
}
