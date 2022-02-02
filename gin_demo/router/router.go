package router

import (
	api "gin_demo/Controller"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", api.Create)
	g.GET("/member", api.GetMember)
	g.GET("/member/list", api.List)
	g.POST("/member/update", api.Update)
	g.POST("/member/delete", api.Delete)

	// 登录

	g.POST("/auth/login", api.Login)
	g.POST("/auth/logout", api.Loginout)
	g.GET("/auth/whoami", api.Whoami)

	// 排课
	g.POST("/course/create")
	g.GET("/course/get")

	g.POST("/teacher/bind_course")
	g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course")
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

}
