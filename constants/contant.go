package constants

const (
	UserTableName      = "user"
	VideoTableName     = "videos"
	LikeVideoTableName = "like_video"
	MessageTableName   = "message"
	FeedVideosCount    = 20
	MySQLDefaultDSN    = "readygo:123456@tcp(localhost:9910)/douyindb?charset=utf8&parseTime=True&loc=Local"
)

const (
	RedisAddr   = "localhost:26379"
	RedisPasswd = "123456"
)

const (
	RedisTokenPtn     = "token:user:%d"
	RedisUserInfoPtn  = "info:user:%d"
	RedisVideoInfoPtn = "info:video:%d"
	RedisFavUserPtn   = "fav:user:%d"
	RedisFavVideoPtn  = "fav:video:%d"
	RedisFollowsPtn   = "follows:%d"
	RedisFansPtn      = "fans:%d"
)

const (
	ConfigHasInvalidSuffix = -1
	ConfigFileAccessError  = -2

	ConfigFileUnSequenceError = -3
	ConfigHasNoEffect         = -4

	RedisInitFailed = -5
)
