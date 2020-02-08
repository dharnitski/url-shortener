package shortener

import (
	"math/big"
	"strconv"

	"github.com/kenshaw/baseconv"
)

// IntToShort converts int64 to base 62 string representation
// base 62 string representation is used later as path in short url
// resulting string is encoded using numbers and
// lower-case letters 'a' to 'z' for digit values 10 to 35, and
// the upper-case letters 'A' to 'Z' for digit values 36 to 61
// negative numbers and 0 are not supported
func IntToShort(i int64) string {
	return big.NewInt(i).Text(62)
}

// ShortToInt converts base 62 decoded staring into int64
// return an error if convertion is not possible
func ShortToInt(s string) (int64, error) {
	// base 10 representation as a string
	encoded, err := baseconv.Decode62ToDec(s)
	if err != nil {
		return 0, err
	}
	// cionvert string to int64
	// function cannot be used directly because it dos not support base bigger than 32
	return strconv.ParseInt(encoded, 10, 64)
}
