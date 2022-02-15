package Controller

import (
	"Course/Form"
	"Course/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

//登录
func Login(c *gin.Context) {
	var loginRequest Form.LoginRequest
	if err := c.Bind(&loginRequest); err != nil {
		panic(err)
	}
	Username := loginRequest.Username
	Password := loginRequest.Password

	var user Form.Member
	global.DB.Where("Username = ?", Username).First(&user)
	if user.Username == "" {
		c.JSON(http.StatusOK, gin.H{
			"Code": Form.WrongPassword,
		})
		return
	}
	if user.Deleted == "1" {
		c.JSON(http.StatusOK, gin.H{
			"Code": Form.WrongPassword,
		})
		return
	}
	if user.Password != Password {
		c.JSON(http.StatusOK, gin.H{
			"Code": Form.WrongPassword,
		})
		return
	}
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		cookie = "NotSet"
		c.SetCookie("camp-session", user.Username, 3600, "/api/v1/", "180.184.74.105", false, true)
	}
	global.LOG.Info(
		"Set Cookie",
		zap.String("Cookie", cookie),
	)

	c.JSON(200, Form.LoginResponse{
		Code: 0,
		Data: struct{ UserID string }{UserID: user.UserID},
	},
	)
}

func Logout(c *gin.Context) {
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		cookie = "NotSet"
		c.JSON(200, Form.LogoutResponse{
			Code: Form.LoginRequired,
		})
		return
	}
	// 设置cookie  MaxAge设置为-1，表示删除cookie
	c.SetCookie("camp-session", cookie, -1, "/api/v1/", "180.184.74.105", false, true)
	global.LOG.Info(
		"Remove Cookie",
		zap.String("Cookie", cookie),
	)
	c.JSON(200, Form.LogoutResponse{
		Code: 0,
	})
}

func Whoami(c *gin.Context) {
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
