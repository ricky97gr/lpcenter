package database

import (
	"context"
	"fmt"
	"time"

	conf "github.com/ricky97gr/lpcenter/server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlClient *gorm.DB

func GetDB() (*gorm.DB, error) {
	if mysqlClient == nil {
		config, _ := conf.GetConfig()
		return InitMySQL(config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := mysqlClient.DB()
	if err != nil {
		config, _ := conf.GetConfig()
		mysqlClient, err = InitMySQL(config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
		if err != nil {
			return nil, err
		}
	}

	err = db.PingContext(ctx)
	if err != nil {
		config, _ := conf.GetConfig()
		mysqlClient, err = InitMySQL(config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
		if err != nil {
			return nil, err
		}
	}
	return mysqlClient, nil
}

func InitMySQL(host, user, password, dbName string, port string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
