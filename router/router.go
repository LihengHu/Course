package router

import (
	api "Course/Controller"
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
	g.POST("/auth/logout", api.Logout)
	g.GET("/auth/whoami", api.Whoami)

	// 排课
	g.POST("/course/create", api.CreateCourse)
	g.GET("/course/get", api.GetCourse)

	g.POST("/teacher/bind_course", api.BindCourse)
	g.POST("/teacher/unbind_course", api.UnbindCourse)
	g.GET("/teacher/get_course", api.GetTeacherCourse)
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course", api.BookCourse)
	g.GET("/student/course")

}
