package dal

import (
	"github.com/qiong-14/EasyDouYin/pkg/constants"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
	"log"
)

var DB *gorm.DB

// Init init DB
func Init() {
	var err error
	log.Println("数据库初始化 DSN:", constants.MySQLDefaultDSN)
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
