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
	"regexp"
	"strings"
)

// CheckUserNameForm 检查用户用户名（邮箱格式）是否正确
func CheckUserNameForm(name string) bool {
	reg := regexp.MustCompile(`^([a-zA-Z\\d][\\w-]{2,})@(\\w{2,})\\.([a-z]{2,})(\\.[a-z]{2,})?$`)
	return reg.Match([]byte(name))
}

// CheckPasswordStrength 密码强度评分
// >= 90: 非常安全 >= 80: 安全（Secure）
// >= 70: 非常强 >= 60: 强（Strong）
// >= 50: 一般（Average）
// >= 25: 弱（Weak） >= 0: 非常弱
func CheckPasswordStrength(password string) int {
	var count int
	pwd := []byte(password)
	lower := regexp.MustCompile(`[a-z]`).Match(pwd)
	upper := regexp.MustCompile(`[A-Z]`).Match(pwd)
	digit := regexp.MustCompile(`[0-9]`).Match(pwd)
	special := regexp.MustCompile(`[\\!#\"$%&'\(\)\*+,-.\/:;<=>?@\[\]\^_\{\|\}~]`).Match(pwd)
	if len(pwd) >= 10 {
		count += 25
	} else if len(pwd) >= 8 && len(pwd) < 10 {
		count += 15
	} else if len(pwd) >= 5 && len(pwd) < 8 {
		count += 5
	}
	if lower && upper {
		count += 20
	} else if lower || upper {
		count += 10
	}
	if upper || lower && digit {
		count += 2
	}
	if upper || lower && digit && special {
		count += 3
	}
	if upper && lower && digit && special {
		count += 5
	}
	// 数字个数
	res := regexp.MustCompile(`[0-9]`).FindAllIndex(pwd, -1)
	if len(res) >= 3 {
		count += 20
	} else if len(res) >= 1 && len(res) <= 3 {
		count += 10
	}
	// 特殊字符个数
	res = regexp.MustCompile(`[\\!#\"$%&'\(\)\*+,-.\/:;<=>?@\[\]\^_\{\|\}~]`).FindAllIndex(pwd, -1)
	if len(res) == 1 {
		count += 10
	} else if len(res) > 1 {
		count += 25
	}

	return count
}

// Encoder password encoder
func Encoder(pwd string) string {
	s := hmac.New(sha256.New, []byte(pwd))
	pwd = hex.EncodeToString(s.Sum(nil))
	return pwd
}

// GetEnvByKey get the environment key, if not exist return  //
func GetEnvByKey(key string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return ""
}

func MatchEnvLoaded() bool {
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

	if !MatchEnvLoaded() {
		_, _ = fmt.Fprintf(os.Stderr, "failed yo get env from json file: %s", filename)
		os.Exit(constants.ConfigHasNoEffect)
	}
	//jsonpath.Read()
}
