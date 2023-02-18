package middleware

// not verified, incomplete
// 冗余程度很高, 但省去了抽象的成本, 每个人阅读自己的那一部分即可

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/qiong-14/EasyDouYin/constants"
	"github.com/qiong-14/EasyDouYin/dal"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	clients         [16]*redis.Client
	TokenClient     *redis.Client
	UserInfoClient  *redis.Client
	VideoInfoClient *redis.Client
	FavUserClient   *redis.Client
	FavVideoClient  *redis.Client
	FollowsClient   *redis.Client
	FansClient      *redis.Client
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
	TokenClient = GetInstance(0)
	UserInfoClient = GetInstance(1)
	VideoInfoClient = GetInstance(2)
	FavUserClient = GetInstance(3)
	FavVideoClient = GetInstance(3)
	FollowsClient = GetInstance(4)
	FansClient = GetInstance(4)
	checkRedis()
}

/*biz start*/

// GetUserTokenRedis token cache, see also SetUserTokenRedis
func GetUserTokenRedis(userId int64) (token string, err error) {
	token, err = TokenClient.Get(fmt.Sprintf(constants.RedisTokenPtn, userId)).Result()
	if err == redis.Nil {
		return "", err
	}
	return token, nil
}

func SetUserTokenRedis(userId int64, token string) error {
	if err := TokenClient.Set(fmt.Sprintf(constants.RedisTokenPtn, userId), token, time.Hour*24).Err(); err != nil {
		return err
	}
	return nil
}

// GetUserInfoRedis user info cache, see also SetUserInfoRedis, renewUserInfoExpire
func GetUserInfoRedis(userId int64) (user dal.User, err error) {
	key := fmt.Sprintf(constants.RedisUserInfoPtn, userId)
	userInfoStr, err := UserInfoClient.Get(key).Result()
	if err == redis.Nil {
		return dal.InvalidUser, err
	}
	if err = json.Unmarshal([]byte(userInfoStr), &user); err != nil {
		// renew
		return dal.InvalidUser, err
	}
	return user, renewExpire(UserInfoClient, key)
}

func SetUserInfoRedis(user dal.User) (err error) {
	userInfoStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	// set
	err = UserInfoClient.Set(
		fmt.Sprintf(constants.RedisUserInfoPtn, user.Id), userInfoStr,
		randomExpire(time.Hour*24)).Err()

	if err != nil {
		return err
	}
	return
}

// GetVideoInfoRedis video info cache, see also SetVideoInfoRedis, renewVideoInfoExpire
// similar to GetUserInfoRedis
func GetVideoInfoRedis(videoId int64) (video dal.VideoInfo, err error) {
	key := fmt.Sprintf(constants.RedisVideoInfoPtn, videoId)
	videoInfoStr, err := VideoInfoClient.Get(key).Result()
	if err == redis.Nil {
		return dal.InvalidVideo, err
	}
	if err = json.Unmarshal([]byte(videoInfoStr), &video); err != nil {
		return dal.InvalidVideo, renewExpire(VideoInfoClient, key)
	}
	return
}

func SetVideoInfoRedis(video dal.VideoInfo) (err error) {
	videoInfoStr, err := json.Marshal(video)
	if err != nil {
		return err
	}
	// set
	err = VideoInfoClient.Set(
		fmt.Sprintf(constants.RedisVideoInfoPtn, video.ID),
		videoInfoStr,
		randomExpire(time.Hour*24),
	).Err()
	if err != nil {
		return err
	}

	return
}

func GetUserFavVideosRedis(userId int64) (videoIdList []int64, err error) {
	key := fmt.Sprintf(constants.RedisFavUserPtn, userId)
	// 改成按时间倒叙，优先展示最后更改的视频
	result, err := FavUserClient.ZRevRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	// 处理
	for _, s := range result {
		id, _ := strconv.ParseInt(s, 10, 64)
		videoIdList = append(videoIdList, id)
	}
	return videoIdList, renewExpire(FavUserClient, key)
}

func GetUserFavVideosCountRedis(userId int64) (res int64, err error) {
	key := fmt.Sprintf(constants.RedisFavUserPtn, userId)
	result, err := FavUserClient.ZCard(key).Result()
	if err != nil {
		return 0, err
	}
	err = renewExpire(FavUserClient, key)
	if err != nil {
		return 0, err
	}
	// 处理
	return result, nil
}
func GetVideosFavRedis(videoId int64) (videoIdList []int64, err error) {
	key := fmt.Sprintf(constants.RedisFavVideoPtn, videoId)
	result, err := FavVideoClient.ZRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	// 处理
	for _, s := range result {
		id, _ := strconv.ParseInt(s, 10, 64)
		videoIdList = append(videoIdList, id)
	}
	return videoIdList, renewExpire(FavVideoClient, key)
}

func GetVideosFavsCountRedis(videoId int64) (res int64, err error) {
	key := fmt.Sprintf(constants.RedisFavVideoPtn, videoId)
	result, err := FavVideoClient.ZCard(key).Result()
	if err != nil {
		return 0, err
	}
	err = renewExpire(FavVideoClient, key)
	if err != nil {
		return 0, err
	}
	// 处理
	return result, nil
}

func ActionUserFavVideoRedis(userId, videoId int64) (err error) {
	// note 正向, 对视频点赞
	key := fmt.Sprintf(constants.RedisFavVideoPtn, videoId)
	// note 如果只需要计数值, 可以使用ZINCRBY等
	err = FavVideoClient.ZAdd(key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.Itoa(int(userId)),
	}).Err()
	if err != nil {
		return err
	}
	err = renewExpire(FavVideoClient, key)
	if err != nil {
		return err
	}

	// note 反向, 添加用户点赞过的视频
	key = fmt.Sprintf(constants.RedisFavUserPtn, userId)
	err = FavUserClient.ZAdd(key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.Itoa(int(videoId)),
	}).Err()
	if err != nil {
		return err
	}
	return renewExpire(FavUserClient, key)
}

func ActionUserUnFavVideoRedis(userId, videoId int64) (err error) {
	// 删除用户点赞记录
	key := fmt.Sprintf(constants.RedisFavVideoPtn, videoId)
	err = FavVideoClient.ZRem(key, strconv.Itoa(int(userId))).Err()
	if err != nil {
		return err
	}
	err = renewExpire(FavVideoClient, key)
	if err != nil {
		return err
	}
	// note 反向, 删除
	key = fmt.Sprintf(constants.RedisFavUserPtn, userId)
	err = FavUserClient.ZRem(key, strconv.Itoa(int(videoId))).Err()
	if err != nil {
		return err
	}
	return renewExpire(FavUserClient, key)
}

// ActionUserFollowRedis userId 关注 userId2
func ActionUserFollowRedis(userId, userId2 int64) (err error) {
	keyFollow, keyFan := fmt.Sprintf(constants.RedisFollowsPtn, userId), fmt.Sprintf(constants.RedisFansPtn, userId2)
	err = FollowsClient.ZAdd(keyFollow, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.Itoa(int(userId2)),
	}).Err()

	if err != nil {
		return err
	}

	err = renewExpire(FollowsClient, keyFollow)
	if err != nil {
		return err
	}

	err = FansClient.ZAdd(keyFan, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: strconv.Itoa(int(userId)),
	}).Err()
	if err != nil {
		return err
	}
	return renewExpire(FansClient, keyFan)
}

// ActionUserUnFollowRedis userId 取消关注 userId2
func ActionUserUnFollowRedis(userId, userId2 int64) (err error) {
	keyFollow, keyFan := fmt.Sprintf(constants.RedisFollowsPtn, userId), fmt.Sprintf(constants.RedisFansPtn, userId2)
	err = FollowsClient.ZRem(keyFollow, strconv.Itoa(int(userId2))).Err()

	if err != nil {
		return err
	}

	err = renewExpire(FollowsClient, keyFollow)
	if err != nil {
		return err
	}

	err = FansClient.ZRem(keyFan, strconv.Itoa(int(userId))).Err()
	if err != nil {
		return err
	}
	return renewExpire(FansClient, keyFan)
}

// GetUserFansCountRedis 获取用户粉丝数
func GetUserFansCountRedis(userId int64) (res int64, err error) {
	key := fmt.Sprintf(constants.RedisFansPtn, userId)
	result, err := FansClient.ZCard(key).Result()
	if err != nil {
		return 0, err
	}
	err = renewExpire(FansClient, key)
	if err != nil {
		return 0, err
	}
	// 处理
	return result, nil
}

// GetUserFollowCountRedis 获取用户粉丝数
func GetUserFollowCountRedis(userId int64) (res int64, err error) {
	key := fmt.Sprintf(constants.RedisFansPtn, userId)
	result, err := FollowsClient.ZCard(key).Result()
	if err != nil {
		return 0, err
	}
	err = renewExpire(FollowsClient, key)
	if err != nil {
		return 0, err
	}
	// 处理
	return result, nil
}

// GetUserFollowsRedis 获取用户关注的所有用户的列表, 返回他们的ID
func GetUserFollowsRedis(userId int64) (res []int64, err error) {
	panic("not implemented")
}

// GetUserFansRedis 获取用户的所有粉丝列表, 返回他们的ID
func GetUserFansRedis(userId int64) (res []int64, err error) {
	panic("not implemented")
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

func randomExpire(baseTs time.Duration) (exp time.Duration) {
	return time.Duration(rand.Intn(100)) + baseTs
}

// renewExpire todo 常数转移
func renewExpire(client *redis.Client, key string) (err error) {
	return client.Expire(key, randomExpire(time.Hour*24)).Err()
}
