package Controller

import (
	"Course/Form"
	"Course/global"
	"Course/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func BookCourse(c *gin.Context) {
	cookie, err := c.Cookie("camp-session")
	if err != nil {
		c.JSON(200, gin.H{
			"Code": Form.LoginRequired,
		})
		return
	}
	var user Form.Member
	global.DB.Where("Username = ?", cookie).First(&user)
	if user.UserType != 2 {
		c.JSON(200, gin.H{
			"Code": Form.StudentNotExisted,
		})
		return
	}
	//TODO: 考虑高并发情况
	CourseID := c.PostForm("CourseID")
	var course Form.Course
	global.DB.Where("Course_ID = ?", CourseID).Find(&course)
	if course.CourseID == "" {
		c.JSON(200, gin.H{
			"Code": Form.CourseNotExisted,
		})
		return
	}
	if course.CourseCap <= 0 {
		c.JSON(200, gin.H{
			"Code": Form.CourseNotAvailable,
		})
		return
	}
	//TODO: 自增ID
	ScheduleID := "1"
	u1 := Form.Schedule{ScheduleID, CourseID, user.UserID}
	global.DB.Create(&u1)
	global.LOG.Info(
		"Book Course",
		zap.String("ScheduleID", ScheduleID),
		zap.String("CourseID", CourseID),
		zap.String("UserID", user.UserID),
	)
	c.JSON(200, types.BookCourseResponse{Code: 0})
}
