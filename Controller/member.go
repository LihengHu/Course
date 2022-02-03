package Controller

import (
	"gin_demo/Form"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
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
	if user.UserType != 1 {
		c.JSON(200, gin.H{
			"Code": Form.LoginRequired,
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
			"Code": Form.ParamInvalid,
		})
		return
	}
	if len(Username) < 8 || len(Username) > 20 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Code": Form.ParamInvalid,
		})
		return
	}
	if len(Password) < 8 || len(Password) > 20 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Code": Form.ParamInvalid,
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
			"Code": Form.ParamInvalid,
		})
	}
	if usertype != 1 && usertype != 2 && usertype != 3 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"Code": Form.ParamInvalid,
		})
		return
	}
	var user1 Form.Member
	db.Where("username = ?", Username).First(&user1)
	if user1.UserID != "" {
		c.JSON(200, gin.H{
			"Code": Form.UserHasExisted,
		})
		return
	}
	u1 := Form.Member{"3", Nickname, Username, Password, Form.UserType(usertype), "0"}
	db.Create(&u1)
	c.JSON(200, Form.CreateMemberResponse{Code: 0, Data: struct{ UserID string }{UserID: u1.UserID}})
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
		c.JSON(http.StatusOK, gin.H{"Code": Form.UserNotExisted})
		return
	}
	if user.Deleted == "1" {
		c.JSON(http.StatusOK, gin.H{"Code": Form.UserHasDeleted})
		return
	}
	c.JSON(200, Form.GetMemberResponse{
		Code: Form.OK,
		Data: struct {
			UserID   string
			Nickname string
			Username string
			UserType Form.UserType
		}{UserID: user.UserID, Nickname: user.Nickname, Username: user.Username, UserType: user.UserType},
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
	c.JSON(200, gin.H{"Code": 0})
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
			"Code": Form.UserNotExisted,
		})
		return
	}
	if user.Deleted == "1" {
		c.JSON(200, gin.H{
			"Code": Form.UserHasDeleted,
		})
		return
	}
	db.Model(&user).Where("user_id = ?", UserID).Update("nickname", Nickname)
	c.JSON(200, gin.H{"Code": 0})
}

func List(c *gin.Context) {

	db, err := gorm.Open("mysql", "root:root@(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 自动迁移
	db.AutoMigrate(&Form.Member{})
	userdb := db.Model(&Form.Member{}).Where(&Form.Member{Deleted: "0"})
	var count int32
	userdb.Count(&count) //总行数
	pageindex, _ := strconv.Atoi(c.PostForm("Offset"))
	pagesize, _ := strconv.Atoi(c.PostForm("Limit"))
	UserList := []Form.Member{}
	userdb.Offset((pageindex - 1) * pagesize).Limit(pagesize).Find(&UserList) //查询pageindex页的数据
	var length int = len(UserList)
	TMemberList := make([]Form.TMember, length)
	for i := 0; i < len(UserList); i++ {
		TMemberList[i].UserID = UserList[i].UserID
		TMemberList[i].Username = UserList[i].Username
		TMemberList[i].UserType = UserList[i].UserType
		TMemberList[i].Nickname = UserList[i].Nickname
	}
	c.JSON(200, Form.GetMemberListResponse{
		Code: 0,
		Data: struct{ MemberList []Form.TMember }{MemberList: TMemberList}})
}
