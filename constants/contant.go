package constants

const (
	UserTableName         = "user"
	VideoTableName        = "videos"
	LikeVideoTableName    = "like_video"
	MessageTableName      = "message"
	UserRelationName   = "user_relation"
	UserInfoName       = "user_info"
	CommentVideoTableName = "comment_video"
	FeedVideosCount       = 20
	MySQLDefaultDSN       = "readygo:123456@tcp(localhost:9910)/douyindb?charset=utf8&parseTime=True&loc=Local"
)
const (
	CheckUserRegisterInfo = false
	CheckUserName         = true
	// CheckUserPassword 检查用户密码强度大于该值
	// 90 以上: 非常安全
	// 80 ~ 90: 安全
	// 70 ~ 80: 非常强
	// 60 ~ 70: 强
	// 50 ~ 60: 一般
	// 25 ~ 50: 弱
	// 0 ~ 25: 非常弱
	CheckUserPassword = 90
)

const (
	RedisAddr   = "localhost:26379"
	RedisPasswd = "123456"
)

const (
	RedisUserIdPtn    = "Id:token:%s"
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
