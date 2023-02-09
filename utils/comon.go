package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Encoder password encoder
func Encoder(pwd string) string {
	s := hmac.New(sha256.New, []byte(pwd))
	pwd = hex.EncodeToString(s.Sum(nil))
	return pwd
}
