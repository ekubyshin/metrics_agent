package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
)

const HashHeader = "HashSHA256"

func HashData(d []byte, k string) ([]byte, error) {
	h := hmac.New(sha256.New, []byte(k))
	_, err := h.Write(d)
	return h.Sum(nil), err
}
