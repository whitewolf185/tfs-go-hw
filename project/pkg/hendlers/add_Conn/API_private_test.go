package add_Conn

import (
	"net/url"
	"testing"
)

type Gen struct {
	PostData      string
	EndpointPath  string
	ApiKeyPrivate string
}

type TestsGenerate struct {
	Name     string
	In       Gen
	Expected string
}

func TestGenerateAuthent(t *testing.T) {
	tests := []TestsGenerate{
		{
			"OK",
			Gen{
				"orderType=mkt&symbol=PI_ETHUSD&side=buy&size=5",
				"/api/v3/sendorder",
				"zEoBlEZQcb7TeUxbiNYsNK19t0xEqHQ0MyUghPEHunIs5x3LP9MqR0631TAaMBY7UWSi6wmLQRytPuYriqWDjmz4",
			},
			"5cHHv5Bb1leGfA+SUeBxO4oBikJME4oRSN9rhaJ8Nq/K5IPprTi8P73LotbzhcROXvB5Feh+KIMKhEauYjU0vg==",
		},
		{
			"Empty test",
			Gen{
				"",
				"",
				"",
			},
			"g89jPWwG72Um1PmUTIkrVxR4UUTVrPr5vyQqbpU5Fu8cgeaFfddTw9Y87W3VvSFqxVXwOrB81+62gQjRj3HXCg==",
		},
	}

	for idx, test := range tests {
		got, err := GenerateAuthent(test.In.PostData, test.In.EndpointPath, test.In.ApiKeyPrivate)
		if err != nil {
			t.Fatalf("Unexpected Error in test №%d %s\n %s", idx, test.Name, err)
		}

		if got != test.Expected {
			t.Errorf("Not match in test %d %s", idx, test.Name)
		}
	}
}

type TestsPostData struct {
	Name     string
	In       map[string]string
	Expected string
}

func TestMakePostData(t *testing.T) {
	tests := []TestsPostData{
		{
			"OK",
			map[string]string{
				"orderType": "mkt",
				"symbol":    "PI_XBTUSD",
				"size":      "1",
				"side":      "buy",
			},
			"orderType=mkt&side=buy&size=1&symbol=PI_XBTUSD",
		},
	}

	for idx, test := range tests {
		got := MakePostData(test.In)

		if got != test.Expected {
			t.Errorf("Test %d expected %s got %s", idx, test.Expected, got)
		}
	}
}

type QE struct {
	u    url.URL
	data map[string]string
}

type TestsMakeQuery struct {
	Name     string
	In       QE
	Expected string
}

func TestMakeQuery(t *testing.T) {
	tests := []TestsMakeQuery{
		{
			"OK",
			QE{
				url.URL{
					Scheme: "https",
					Host:   "demo-futures.kraken.com",
					Path:   "derivatives",
				},
				map[string]string{
					"orderType": "mkt",
					"symbol":    "PI_XBTUSD",
					"size":      "1",
					"side":      "buy",
				},
			},
			"https://demo-futures.kraken.com/derivatives?orderType=mkt&side=buy&size=1&symbol=PI_XBTUSD",
		},
	}

	for idx, test := range tests {
		MakeQuery(&test.In.u, test.In.data)

		if test.In.u.String() != test.Expected {
			t.Errorf("Test %d expected %s got %s", idx, test.Expected, test.In.u.String())
		}
	}
}
