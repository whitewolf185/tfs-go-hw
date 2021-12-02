package add_Conn

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestsConvert struct {
	Name     string
	In       string
	Expected float32
}

func TestConvertToFloat(t *testing.T) {
	tests := []TestsConvert{
		{
			"OK",
			"23467.3",
			23467.3,
		},

		{
			"Not a number",
			"234ad67.3",
			-1,
		},
	}

	for idx, test := range tests {
		got := ConvertToFloat(test.In)

		if got != test.Expected {
			t.Errorf("Test %d expected %f got %f", idx, test.Expected, got)
		}
	}
}

func SetUpENV() error {
	err := os.Setenv("TG_BOT_TOKEN", "D:/Documents/GO_projects/tfs-go-hw/project/tgBot_token.txt")
	if err != nil {
		return err
	}

	err = os.Setenv("WS_URL", "wss://demo-futures.kraken.com/ws/v1?chart")
	if err != nil {
		return err
	}

	err = os.Setenv("TOKEN_PATH_PUBLIC", "D:/Documents/GO_projects/tfs-go-hw/project/public-token_test.txt")
	if err != nil {
		return err
	}

	err = os.Setenv("TOKEN_PATH_PRIVATE", "D:/Documents/GO_projects/tfs-go-hw/project/private-token_test.txt")
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
