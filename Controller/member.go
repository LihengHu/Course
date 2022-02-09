package Controller

import (
	"Course/Form"
	"Course/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		cookie = "NotSet"
		c.JSON(200, gin.H{
			"Code": Form.LoginRequired,
		})
		return
	}
	var user Form.Member
	global.DB.Where("Username = ?", cookie).First(&user)
	if user.UserType != Form.Admin {
		c.JSON(200, gin.H{
			"Code": Form.PermDenied,
		})
		return
	}
	//TODO: 生成自增ID
	var size int64
	global.DB.Table("members").Count(&size)
	var firstUser Form.Member
	global.DB.Table("members").Offset(int(size - 1)).Limit(1).Find(&firstUser)
	oldId, _ := strconv.ParseInt(firstUser.UserID, 10, 64)
	UserID := strconv.FormatInt(oldId+1, 10)
	var create Form.CreateMemberRequest
	c.Bind(&create)
	Nickname := create.Nickname
	Username := create.Username
	Password := create.Password
	UserType := create.UserType
	if len(Nickname) < 4 || len(Nickname) > 20 {
		c.JSON(http.StatusOK, gin.H{
			"Code": Form.ParamInvalid,
		})
		return
	}
	if len(Username) < 8 || len(Username) > 20 {
		c.JSON(http.StatusOK, gin.H{
			"Code": Form.ParamInvalid,
		})
		return
	}
	if len(Password) < 8 || len(Password) > 20 {
		c.JSON(http.StatusOK, gin.H{
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
		c.JSON(http.StatusOK, gin.H{
			"Code": Form.ParamInvalid,
		})
	}
	if UserType != Form.Admin && UserType != Form.Student && UserType != Form.Teacher {
		c.JSON(http.StatusOK, gin.H{
			"Code": Form.ParamInvalid,
		})
		return
	}
	var user1 Form.Member
	global.DB.Where("username = ?", Username).First(&user1)
	if user1.UserID != "" {
		c.JSON(200, gin.H{
			"Code": Form.UserHasExisted,
		})
		return
	}
	err = global.RDB.Del(global.CTX, "members").Err()
	if err != nil {
		panic(err)
	}
	u1 := Form.Member{UserID, Nickname, Username, Password, UserType, "0"}
	global.DB.Create(&u1)
	global.LOG.Info(
		"Create Member",
		zap.String("UserID", UserID),
		zap.String("Username", Username),
	)
	c.JSON(200, Form.CreateMemberResponse{Code: 0, Data: struct{ UserID string }{UserID: u1.UserID}})
}

func GetMember(c *gin.Context) {
	UserID := c.PostForm("UserID")
	var user Form.Member
	global.DB.Where("user_id = ?", UserID).First(&user)
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
	var deleteRequest Form.DeleteMemberRequest
	c.Bind(&deleteRequest)
	UserID := deleteRequest.UserID
	err := global.RDB.Del(global.CTX, "members").Err()
	if err != nil {
		panic(err)
	}
	var user Form.Member
	global.DB.Model(&user).Where("user_id = ?", UserID).Update("deleted", "1")
	global.LOG.Info(
		"Delete Member",
		zap.String("UserID", UserID),
	)
	c.JSON(200, gin.H{"Code": 0})
}

func Update(c *gin.Context) {
	var updateRequest Form.UpdateMemberRequest
	c.Bind(&updateRequest)
	UserID := updateRequest.UserID
	Nickname := updateRequest.Nickname
	err := global.RDB.Del(global.CTX, "members").Err()
	if err != nil {
		panic(err)
	}
	var user Form.Member
	global.DB.Where("user_id = ?", UserID).First(&user)
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
	global.DB.Model(&user).Where("user_id = ?", UserID).Update("nickname", Nickname)
	global.LOG.Info(
		"Update Member",
		zap.String("UserID", UserID),
		zap.String("new Nickname", Nickname),
	)
	c.JSON(200, gin.H{"Code": 0})
}

func List(c *gin.Context) {
	userdb := global.DB.Model(&Form.Member{}).Where(&Form.Member{Deleted: "0"})
	var count int64
	userdb.Count(&count) //总行数
	var getMember Form.GetMemberListRequest
	c.Bind(&getMember)
	pageindex := getMember.Offset
	pagesize := getMember.Limit
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
