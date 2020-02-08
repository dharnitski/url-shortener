package shortener_test

import (
	"math"
	"testing"

	"github.com/dharnitski/url-shortener/shortener"
	"github.com/stretchr/testify/assert"
)

// Numbers converted to base 62
// test cases are designed to cover base 62 scenarios
func TestConvertions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		i int64
		s string
	}{
		{
			1,
			"1",
		},
		{
			9,
			"9",
		},
		{
			9,
			"9",
		},
		{
			10,
			"a",
		},
		{
			61,
			"Z",
		},
		{
			62,
			"10",
		},
		{
			123456789,
			"8m0Kx",
		},
		{
			math.MaxInt64, // 9223372036854775807
			"aZl8N0y58M7",
		},
		{
			math.MaxInt64 - 1,
			"aZl8N0y58M6",
		},
	}
	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			actualS := shortener.IntToShort(test.i)
			assert.Equal(t, test.s, actualS)
			actualI, err := shortener.ShortToInt(actualS)
			assert.NoError(t, err)
			assert.Equal(t, test.i, actualI)
		})
	}
}

func TestEncodingErrors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		s string
	}{
		{
			"___",
		},
		{
			"-",
		},
		{
			"111111111111111111111111111111111",
		},
		{
			"😊",
		},
	}
	for _, test := range tests {
		t.Run(test.s, func(t *testing.T) {
			_, err := shortener.ShortToInt(test.s)
			assert.Error(t, err)
		})
	}
}
