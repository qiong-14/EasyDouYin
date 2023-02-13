package constants

import (
	"fmt"
)

var (
	MYSQLHost       = "127.0.0.1"
	MYSQLPort       = "9910"
	MYSQLUser       = "readygo"
	MYSQLPwd        = "123456"
	MYSQLDb         = "douyindb"
	MySQLDefaultDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		MYSQLUser,
		MYSQLPwd,
		MYSQLHost,
		MYSQLPort,
		MYSQLDb)
)

const (
	// UserTableName mysql table name
	UserTableName   = "user"
	VideoTableName  = "videos"
	FeedVideosCount = 20
)
