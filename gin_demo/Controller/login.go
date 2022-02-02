package Controller

import (
	"fmt"
	"gin_demo/Form"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

//登录
func Login(c *gin.Context) {
	Username := c.PostForm("Username")
	Password := c.PostForm("Password")
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 自动迁移
	db.AutoMigrate(&Form.Member{})
	var user Form.Member
	db.Where("Username = ?", Username).First(&user)
	if user.Deleted == "1" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"ErrNo": Form.WrongPassword,
			"msg":   "密码错误",
		})
	}
	if user.Username == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"ErrNo": Form.WrongPassword,
			"msg":   "密码错误",
		})
		return
	}
	if user.Password != Password {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"ErrNo": Form.WrongPassword,
			"msg":   "密码错误",
		})
		return
	}
	cookie, err := c.Cookie("camp-seesion")
	if err != nil {
		cookie = "NotSet"
		c.SetCookie("camp-seesion", user.Username, 3600, "/api/v1/", "localhost", false, true)
	}
	fmt.Printf("Cookie value: %s \n", cookie)
	c.JSON(200, gin.H{
		"status":   "posted",
		"Password": user.Password,
		"Username": user.Username,
	})

}

func Loginout(c *gin.Context) {
	cookie, err := c.Cookie("camp-seesion")
	if err != nil {
		cookie = "NotSet"
		return
	}
	// 设置cookie  MaxAge设置为-1，表示删除cookie
	c.SetCookie("camp-seesion", cookie, -1, "/api/v1/", "localhost", false, true)
	c.String(200, "登出成功")
}

func Whoami(c *gin.Context) {
	cookie, err := c.Cookie("camp-seesion")
	if err != nil {
		cookie = "NotSet"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"ErrNo": Form.LoginRequired,
			"msg":   "用户未登陆",
		})
		return
	}
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 自动迁移
	db.AutoMigrate(&Form.Member{})
	var user Form.Member
	db.Where("Username = ?", cookie).First(&user)

	c.JSON(200, gin.H{
		"用户ID": user.UserID,
		"用户昵称": user.Nickname,
		"用户类型": user.UserType,
	})
}
