package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(src string) string {
	h := sha256.New()
	h.Write([]byte(src))

	return hex.EncodeToString(h.Sum([]byte("")))
}

func Sha256Bytes(src []byte) string {
	h := sha256.New()
	h.Write(src)

	return hex.EncodeToString(h.Sum([]byte("")))
}
