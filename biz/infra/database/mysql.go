package database

import (
	"fmt"
	"github.com/li1553770945/personal-file-service/biz/infra/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
)

func NewDatabase(conf *config.Config) *gorm.DB {
	dbconfig := conf.DatabaseConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s", dbconfig.Username, dbconfig.Password, dbconfig.Address, dbconfig.Port, dbconfig.Database, strconv.FormatBool(dbconfig.UseTLS))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}
	return db
}
