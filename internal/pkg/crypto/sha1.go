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
