package addConn

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// Root folder of this project
	Root = filepath.ToSlash(filepath.Join(filepath.Dir(b), "../../.."))
)

func SetUpENV() error {
	fmt.Println(Root)
	err := os.Setenv("TG_BOT_TOKEN", Root+"/cmd/config/tgBot_token_test.txt")
	if err != nil {
		return err
	}

	err = os.Setenv("WS_URL", "wss://demo-futures.kraken.com/ws/v1?chart")
	if err != nil {
		return err
	}

	err = os.Setenv("TOKEN_PATH_PUBLIC", Root+"/cmd/config/public-token_test.txt")
	if err != nil {
		return err
	}

	err = os.Setenv("TOKEN_PATH_PRIVATE", Root+"/cmd/config/private-token_test.txt")
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
