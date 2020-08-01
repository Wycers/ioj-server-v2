package crypto

import (
	"encoding/base64"
)

func Base64Encode(data []byte) string {
	// Go supports both standard and URL-compatible
	// base64. Here's how to encode using the standard
	// encoder. The encoder requires a `[]byte` so we
	// cast our `string` to that type.
	res := base64.StdEncoding.EncodeToString(data)

	return res
}

func Base64Decode(data string) ([]byte, error) {
	// Decoding may return an error, which you can check
	// if you don't already know the input to be
	// well-formed.
	res, err := base64.StdEncoding.DecodeString(data)

	return res, err
}
