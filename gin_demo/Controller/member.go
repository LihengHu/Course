package Controller

import (
	"gin_demo/Form"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		cookie = "NotSet"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": Form.LoginRequired,
			"msg":  "用户未登陆",
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
	if user.UserType != 1 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"ErrNo": Form.LoginRequired,
			"msg":   "权限不够",
		})
		return
	}
	Nickname := c.PostForm("Nickname")
	Username := c.PostForm("Username")
	Password := c.PostForm("Password")
	UserType := c.PostForm("UserType")
	usertype, _ := strconv.Atoi(UserType)
	if len(Nickname) < 4 || len(Nickname) > 20 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"ErrNo": Form.ParamInvalid,
			"msg":   "参数不合法",
		})
		return
	}
	if len(Username) < 8 || len(Username) > 20 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"ErrNo": Form.ParamInvalid,
			"msg":   "参数不合法",
		})
		return
	}
	if len(Password) < 8 || len(Password) > 20 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"ErrNo": Form.ParamInvalid,
			"msg":   "参数不合法",
		})
		return
	}
	count1 := 0
	count2 := 0
	count3 := 0
	for i := 0; i < len(Password); i++ {
		if Password[i] >= '0' && Password[i] <= '9' {
			count1++
		} else if Password[i] >= 'A' && Password[i] <= 'Z' {
			count2++
		} else if Password[i] >= 'a' && Password[i] <= 'z' {
			count3++
		}
	}
	if count1 == 0 || count2 == 0 || count3 == 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": Form.ParamInvalid,
			"msg":  "参数不合法",
		})
	}
	if usertype != 1 && usertype != 2 && usertype != 3 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": Form.ParamInvalid,
			"msg":  "参数不合法",
		})
		return
	}
	var user1 Form.Member
	db.Where("username = ?", Username).First(&user1)
	if user1.UserID != "" {
		c.JSON(200, gin.H{
			"code": Form.UserHasExisted,
		})
		return
	}
	u1 := Form.Member{"3", Nickname, Username, Password, Form.UserType(usertype), "0"}
	db.Create(&u1)
}

func GetMember(c *gin.Context) {
	UserID := c.PostForm("UserID")
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 自动迁移
	db.AutoMigrate(&Form.Member{})
	var user Form.Member
	db.Where("user_id = ?", UserID).First(&user)
	if user.Deleted == "" {
		c.JSON(http.StatusOK, gin.H{"code": Form.UserNotExisted})
		return
	}
	if user.Deleted == "1" {
		c.JSON(http.StatusOK, gin.H{"code": Form.UserHasDeleted})
		return
	}
	c.JSON(200, gin.H{
		"UserID":   user.UserID,
		"Nickname": user.Nickname,
		"Username": user.Username,
		"UserType": user.UserType,
	})
}
func Delete(c *gin.Context) {
	UserID := c.PostForm("UserID")
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 自动迁移
	db.AutoMigrate(&Form.Member{})
	var user Form.Member
	db.Model(&user).Where("user_id = ?", UserID).Update("deleted", "1")
}

func Update(c *gin.Context) {
	UserID := c.PostForm("UserID")
	Nickname := c.PostForm("Nickname")
	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 自动迁移
	db.AutoMigrate(&Form.Member{})
	var user Form.Member
	db.Where("user_id = ?", UserID).First(&user)
	if user.UserID == "" {
		c.JSON(200, gin.H{
			"code": Form.UserNotExisted,
		})
		return
	}
	if user.Deleted == "1" {
		c.JSON(200, gin.H{
			"code": Form.UserHasDeleted,
		})
		return
	}
	db.Model(&user).Where("user_id = ?", UserID).Update("nickname", Nickname)
}

func List(c *gin.Context) {

}
