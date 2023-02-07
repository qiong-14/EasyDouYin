package utils

import "os"

func GetEnvByKey(key string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return ""
}
