package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/rubberduckkk/ducker/internal/infra/config"
)

var (
	mysqlDB *gorm.DB
)

func Init(cfg config.MySQL) error {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: cfg.DSN,
	}))
	if err != nil {
		return err
	}
	mysqlDB = db
	return nil
}

func Instance() *gorm.DB {
	return mysqlDB
}
