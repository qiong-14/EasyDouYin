package dal

import (
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
	"log"
)

var DB *gorm.DB
var rDB *redis.Client

// DBInit init DB
func Init() {
	mysqlInit()
	redisInit()
}

func mysqlInit() {
	var err error
	log.Println("mysql数据库初始化 DSN:", constants.MySQLDefaultDSN)
	DB, err = gorm.Open(mysql.Open(constants.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}
	if err = DB.Use(gormopentracing.New()); err != nil {
		panic(err)
	}
}
func redisInit() {
	log.Println("redis数据库初始化 Addr:", constants.REDISAddr)

	rDB = redis.NewClient(&redis.Options{
		Addr:     constants.REDISAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

}
