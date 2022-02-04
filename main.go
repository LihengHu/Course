package main

import (
	"Course/global"
	"Course/initialize"
	Types "Course/router"
	"context"
	"github.com/gin-gonic/gin"
)

func main() {
	//TODO: 后期数据库要用Viper从yaml文件中读取，避免写死
	router := gin.Default()
	Types.RegisterRouter(router)
	global.LOG = initialize.Zap()
	global.DB = initialize.GormMysql()
	global.RDB = initialize.Redis()
	global.CTX = context.Background()
	if global.LOG != nil {
		defer global.LOG.Sync()
	}
	if global.DB != nil {
		db, _ := global.DB.DB()
		initialize.RegisterTables(global.DB)
		defer db.Close()
	}
	router.Run(":80")
}
