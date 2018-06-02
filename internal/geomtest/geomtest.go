package geomtest

import "encoding/hex"

// MustHexDecode decodes hex bytes from s. It panics on any error.
func MustHexDecode(s string) []byte {
	data, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
