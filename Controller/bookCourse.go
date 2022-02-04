package Controller

import (
	"Course/Form"
	"Course/global"
	"Course/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

func BookCourse(c *gin.Context) {
	//TODO: 考虑高并发情况
	StudentID := c.PostForm("StudentID")
	CourseID := c.PostForm("CourseID")

	// 检验学生是否存在
	count, err := global.RDB.SCard(global.CTX, "members").Result()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		var redisMembers []string
		global.DB.Table("members").Select("User_ID").Where("User_Type = ?", 2).Find(&redisMembers)
		err := global.RDB.SAdd(global.CTX, "members", redisMembers).Err()
		if err != nil {
			panic(err)
		}
	}
	isStudent, err := global.RDB.SIsMember(global.CTX, "members", StudentID).Result()
	if err != nil {
		panic(err)
	}
	if !isStudent {
		c.JSON(200, gin.H{
			"Code": Form.StudentNotExisted,
		})
		return
	}

	// 检验课程是否存在
	count, err = global.RDB.HLen(global.CTX, "courses").Result()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		var courses []*Form.RedisCourse
		global.DB.Table("courses").Select("Course_ID,Course_Cap").Find(&courses)
		redisCourses := make(map[string]interface{})
		for _, course := range courses {
			redisCourses[course.Course_ID] = course.Course_Cap
		}
		err := global.RDB.HMSet(global.CTX, "courses", redisCourses).Err()
		if err != nil {
			panic(err)
		}
	}
	isCourse, err := global.RDB.HExists(global.CTX, "courses", CourseID).Result()
	if !isCourse {
		c.JSON(200, gin.H{
			"Code": Form.CourseNotExisted,
		})
		return
	}
	if err != nil {
		panic(err)
	}
	courseCapStr, err := global.RDB.HGet(global.CTX, "courses", CourseID).Result()
	if err != nil {
		panic(err)
	}
	courseCap, _ := strconv.Atoi(courseCapStr)
	if courseCap <= 0 {
		c.JSON(200, gin.H{
			"Code": Form.CourseNotAvailable,
		})
		return
	}
	var schedule Form.Schedule
	global.DB.Where("Course_ID = ? and Student_ID = ?", CourseID, StudentID).First(&schedule)
	if schedule.CourseID != "" {
		c.JSON(200, gin.H{
			"Code": Form.StudentHasCourse,
		})
		return
	}
	u1 := Form.Schedule{StudentID, CourseID}
	global.DB.Create(&u1)
	global.LOG.Info(
		"Book Course",
		zap.String("CourseID", CourseID),
		zap.String("UserID", StudentID),
	)
	c.JSON(200, types.BookCourseResponse{Code: 0})
}
