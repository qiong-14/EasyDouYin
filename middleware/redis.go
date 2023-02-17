package middleware

// not verified, incomplete

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/qiong-14/EasyDouYin/constants"
	"github.com/qiong-14/EasyDouYin/dal"
	"math/rand"
	"os"
	"time"
)

var (
	clients [16]*redis.Client
	tokenClient,
	userInfoClient,
	videoInfoClient,
	favUserClient,
	favVideoClient,
	followsClient,
	fansClient *redis.Client
)

func GetInstance(bucket int) *redis.Client {
	if bucket > 16 {
		return nil
	}
	if clients[bucket] != nil {
		return clients[bucket]
	}
	client := redis.NewClient(&redis.Options{
		Addr:     constants.RedisAddr,
		Password: constants.RedisPasswd,
		DB:       bucket,
	})

	clients[bucket] = client
	return client
}

func InitRedis() {
	tokenClient = GetInstance(0)
	userInfoClient = GetInstance(1)
	videoInfoClient = GetInstance(2)
	favUserClient = GetInstance(3)
	favVideoClient = GetInstance(3)
	followsClient = GetInstance(4)
	fansClient = GetInstance(4)
	checkRedis()
}

/*biz start*/

func GetUserTokenRedis(userId int64) (token string, err error) {
	token, err = tokenClient.Get(fmt.Sprintf(constants.RedisTokenPtn, userId)).Result()
	if err == redis.Nil {
		return "", err
	}
	return token, nil
}

func SetUserTokenRedis(userId int64, token string) error {
	if err := tokenClient.Set(fmt.Sprintf(constants.RedisTokenPtn, userId), token, time.Hour*24).Err(); err != nil {
		return err
	}
	return nil
}

func GetUserInfoRedis(userId int64) (user dal.User, err error) {
	userInfoStr, err := userInfoClient.Get(fmt.Sprintf(constants.RedisUserInfoPtn, userId)).Result()
	if err == redis.Nil {
		return dal.InvalidUser, err
	}
	if err = json.Unmarshal([]byte(userInfoStr), &user); err != nil {
		return dal.InvalidUser, nil
	}
	return
}

func SetUserInfoRedis(user dal.User) (err error) {
	userInfoStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	// set
	err = userInfoClient.Set(fmt.Sprintf(constants.RedisUserInfoPtn, user.Id), userInfoStr,
		randomExpire(time.Hour*24),
	).Err()
	if err != nil {
		return err
	}

	return
}

func RenewUserInfoExpire(userId int64) (err error) {
	err = userInfoClient.Expire(fmt.Sprintf(constants.RedisUserInfoPtn, userId), randomExpire(time.Hour*24)).Err()
	return err
}

func GetVideoInfoRedis(videoId int64) (video dal.VideoInfo, err error) {
	videoInfoStr, err := videoInfoClient.Get(fmt.Sprintf(constants.RedisVideoInfoPtn, videoId)).Result()
	if err == redis.Nil {
		return dal.InvalidVideo, err
	}
	if err = json.Unmarshal([]byte(videoInfoStr), &video); err != nil {
		return dal.InvalidVideo, nil
	}
	return
}

func SetVideoInfoRedis(video dal.VideoInfo) (err error) {
	videoInfoStr, err := json.Marshal(video)
	if err != nil {
		return err
	}
	// set
	err = videoInfoClient.Set(fmt.Sprintf(constants.RedisVideoInfoPtn, video.ID), videoInfoStr,
		randomExpire(time.Hour*24),
	).Err()
	if err != nil {
		return err
	}

	return
}

func RenewVideoInfoExpire(userId int64) (err error) {
	err = videoInfoClient.Expire(fmt.Sprintf(constants.RedisVideoInfoPtn, userId), randomExpire(time.Hour*24)).Err()
	return err
}

/*biz end*/

func checkRedis() {
	for i := 0; i < 5; i++ {
		if clients[i] == nil {
			os.Exit(constants.RedisInitFailed)
		}
	}
}

// get for test
func get(bucket int, key string) (string, error) {
	return GetInstance(bucket).Get(key).Result()
}

// set for test
func set(bucket int, key string, val string) (string, error) {
	return GetInstance(bucket).Set(key, val, time.Duration(-1)).Result()
}

// flushDB for test
func flushDB(bucket int) (string, error) {
	return GetInstance(bucket).FlushDB().Result()
}

// randomExpire todo 常数转移
func randomExpire(baseTs time.Duration) (exp time.Duration) {
	return time.Duration(rand.Intn(100)) + baseTs
}
