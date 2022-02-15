package initialize

import (
	"Course/Form"
	"Course/global"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func GormMysql() *gorm.DB {
	dsn := "root:bytedancecamp@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local"
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		return db
	}
}

func Zap() (logger *zap.Logger) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Printf("can't initialize zap logger: %v\n", err)
	}
	return logger
}

func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		&Form.Member{},
		&Form.Member{},
	)
	if err != nil {
		global.LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.LOG.Info("register table success")
}

func Redis() (rdb *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 30,
	})
	err := rdb.SAdd(global.CTX, "schedules", "nothing").Err()
	if err != nil {
		panic(err)
	}
	return rdb
}
