package tools

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/qiong-14/EasyDouYin/constants"
	"io/ioutil"
	"os"
	"strings"
)

// Encoder password encoder
func Encoder(pwd string) string {
	s := hmac.New(sha256.New, []byte(pwd))
	pwd = hex.EncodeToString(s.Sum(nil))
	return pwd
}

// GetEnvByKey get the environment key, if not exist return  ""
func GetEnvByKey(key string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return ""
}

func TestEnvLoaded() bool {
	var envKeys = []string{
		"MINIO_ENDPOINT",
		"MINIO_ACCESS_KEY",
		"MINIO_SECRET_KEY",
		"MINIO_BUCKET",
	}
	for _, key := range envKeys {
		if _, exist := os.LookupEnv(key); !exist {
			return false
		}
	}
	return true
}

type Configs struct {
	MinioEndpoint  string `json:"MINIO_ENDPOINT"`
	MinioAccessKey string `json:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `json:"MINIO_SECRET_KEY"`
	MinioBucket    string `json:"MINIO_BUCKET"`
}

func LoadEnvFromJsonCfg(filename string) {
	if !strings.HasSuffix(filename, ".json") {
		_, _ = fmt.Fprintf(os.Stderr, "%s is not a json file", filename)
		os.Exit(constants.ConfigHasInvalidSuffix)
	}
	config := Configs{}
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "access error or file not exist: %s ", filename)
		os.Exit(constants.ConfigFileAccessError)
	}
	if err := json.Unmarshal(bytes, &config); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unsequence error: %s ", filename)
		os.Exit(constants.ConfigFileUnSequenceError)
	}
	_ = os.Setenv("MINIO_ENDPOINT", config.MinioEndpoint)
	_ = os.Setenv("MINIO_ACCESS_KEY", config.MinioAccessKey)
	_ = os.Setenv("MINIO_SECRET_KEY", config.MinioSecretKey)
	_ = os.Setenv("MINIO_BUCKET", config.MinioBucket)

	if !TestEnvLoaded() {
		_, _ = fmt.Fprintf(os.Stderr, "failed yo get env from json file: %s", filename)
		os.Exit(constants.ConfigHasNoEffect)
	}
	//jsonpath.Read()
}
