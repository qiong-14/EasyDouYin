package tools

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// Encoder password encoder
func Encoder(pwd string) string {
	s := hmac.New(sha256.New, []byte(pwd))
	pwd = hex.EncodeToString(s.Sum(nil))
	return pwd
}
func GetEnvByKey(key string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return ""
}
