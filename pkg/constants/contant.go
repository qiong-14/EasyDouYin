package constants

const (
	// UserTableName mysql table name
	UserTableName   = "user"
	VideoTableName  = "videos"
	FeedVideosCount = 20

	// JWTSecret jwt params
	JWTSecret       = "ReadGo"
	JWTAudience     = "douyin"
	JWTDuration     = 24 * 60 * 60 // 一天时间，单位为秒
	JWTIssuer       = "ReadGo"
	JWTSubject      = "jwtToken"
	MySQLDefaultDSN = "readygo:123456@tcp(localhost:9910)/douyindb?charset=utf8&parseTime=True&loc=Local"
)
