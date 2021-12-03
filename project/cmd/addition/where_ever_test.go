package addition

import (
	"errors"
	"testing"

	"github.com/whitewolf185/fs-go-hw/project/cmd/addition/MyErrors"
)

type Converted struct {
	num float32
	err error
}

type TestsConvert struct {
	Name     string
	In       string
	Expected Converted
}

func TestConvertToFloat(t *testing.T) {
	tests := []TestsConvert{
		{
			"OK",
			"23467.3",
			Converted{23467.3, nil},
		},

		{
			"Not a number",
			"234ad67.3",
			Converted{float32(-1), MyErrors.ConvertErr(errors.New("strconv.ParseFloat: parsing \"234ad67.3\": invalid syntax"))},
		},
	}

	for idx, test := range tests {
		var got Converted
		got.num, got.err = ConvertToFloat(test.In)

		if got.num != test.Expected.num {
			t.Errorf("Test %d expected %f %s got %f %s", idx, test.Expected.num, test.Expected.err.Error(),
				got.num, got.err.Error())
		}
	}
}
