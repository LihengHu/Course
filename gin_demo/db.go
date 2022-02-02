package main

import (
	Form "gin_demo/Form"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 自动迁移
	db.AutoMigrate(&Form.Member{})

	u1 := Form.Member{"21721217", "AncientMoon", "hlh", "111", Form.Student}
	u2 := Form.Member{"21721218", "lol", "zwj", "222", Form.Student}
	u3 := Form.Member{"1", "admin", "JudgeAdmin", "JudgePassword2022", Form.Admin}
	// 创建记录
	db.Create(&u1)
	db.Create(&u2)
	db.Create(&u3)

}
