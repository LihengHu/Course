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
			"Code": Form.WrongPassword,
		})
	}
	if user.Username == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Code": Form.WrongPassword,
		})
		return
	}
	if user.Password != Password {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Code": Form.WrongPassword,
		})
		return
	}
	cookie, err := c.Cookie("camp-seesion")
	if err != nil {
		cookie = "NotSet"
		c.SetCookie("camp-seesion", user.Username, 3600, "/api/v1/", "localhost", false, true)
	}
	fmt.Printf("Cookie value: %s \n", cookie)
	c.JSON(200, Form.LoginResponse{
		Code: 0,
		Data: struct{ UserID string }{UserID: user.UserID},
	},
	)

}

func Loginout(c *gin.Context) {
	cookie, err := c.Cookie("camp-seesion")
	if err != nil {
		cookie = "NotSet"
		return
	}
	// 设置cookie  MaxAge设置为-1，表示删除cookie
	c.SetCookie("camp-seesion", cookie, -1, "/api/v1/", "localhost", false, true)
	c.JSON(200, Form.LogoutResponse{
		Code: 0,
	})
}

func Whoami(c *gin.Context) {
	cookie, err := c.Cookie("camp-seesion")
	if err != nil {
		cookie = "NotSet"
		c.JSON(200, gin.H{
			"Code": Form.LoginRequired,
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

	c.JSON(200, Form.WhoAmIResponse{
		Code: 0,
		Data: struct {
			UserID   string
			Nickname string
			Username string
			UserType Form.UserType
		}{UserID: user.UserID, Nickname: user.Nickname, Username: user.Username, UserType: user.UserType},
	},
	)
}
