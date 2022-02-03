package Controller

import (
	"Course/Form"
	"Course/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//创建课程
func CreateCourse(c *gin.Context) {
	var createCourseRequest Form.CreateCourseRequest
	// 参数有误
	if err := c.Bind(&createCourseRequest); err != nil || len(createCourseRequest.Name) == 0 || createCourseRequest.Cap <= 0 {
		c.JSON(http.StatusUnprocessableEntity, Form.CreateCourseResponse{
			Code: Form.ParamInvalid,
			Data: struct {
				CourseID string
			}{},
		})
		return
	}

	// 获取最新的记录
	var len int64
	global.DB.Table("courses").Count(&len)
	var firstCourse Form.TCourse
	global.DB.Table("courses").Offset(int(len - 1)).Limit(1).Find(&firstCourse)
	oldId, _ := strconv.ParseInt(firstCourse.CourseID, 10, 64)
	fmt.Println(oldId + 1)

	newCourse := Form.TCourse{
		CourseID:  strconv.FormatInt(oldId+1, 10),
		Name:      createCourseRequest.Name,
		TeacherID: "-1",
		Cap:       createCourseRequest.Cap,
	}

	global.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").Table("courses").Create(&newCourse)
	c.JSON(http.StatusOK, Form.CreateCourseResponse{
		Code: Form.OK,
		Data: struct {
			CourseID string
		}{newCourse.CourseID},
	})
}

//获取课程
func GetCourse(c *gin.Context) {
	// 创建变量绑定输入
	var getCourseRequest Form.GetCourseRequest
	// 参数有误
	if err := c.Bind(&getCourseRequest); err != nil || len(getCourseRequest.CourseID) == 0 {
		c.JSON(http.StatusUnprocessableEntity, Form.GetCourseResponse{
			Code: Form.ParamInvalid,
			Data: Form.TCourse{},
		})
		return
	}

	// 根据courseId获取对应课程
	var getCourse Form.TCourse
	global.DB.Table("courses").Where("course_id=?", getCourseRequest.CourseID).Find(&getCourse)

	// 课程不存在
	if getCourse.CourseID == "" {
		c.JSON(http.StatusUnprocessableEntity, Form.GetCourseResponse{
			Code: Form.CourseNotExisted,
			Data: Form.TCourse{},
		})
	}
	// 返回课程
	c.JSON(http.StatusOK, Form.GetCourseResponse{
		Code: Form.OK,
		Data: getCourse,
	})

}

//绑定课程
func BindCourse(c *gin.Context) {
	// 创建变量绑定输入
	var bindCourseRequest Form.BindCourseRequest
	// 参数不合法
	if err := c.Bind(&bindCourseRequest); err != nil || len(bindCourseRequest.CourseID) == 0 || len(bindCourseRequest.TeacherID) == 0 {
		c.JSON(http.StatusUnprocessableEntity, Form.BindCourseResponse{Code: Form.ParamInvalid})
		return
	}

	// 根据courseId获取对应课程
	var getCourse Form.TCourse
	global.DB.Table("courses").Where("course_id=?", bindCourseRequest.CourseID).Find(&getCourse)

	// 课程不存在
	if getCourse.CourseID == "" {
		c.JSON(http.StatusUnprocessableEntity, Form.BindCourseResponse{Code: Form.CourseNotExisted})
		return
	}

	// 课程已绑定
	if getCourse.TeacherID != "-1" {
		c.JSON(http.StatusUnprocessableEntity, Form.BindCourseResponse{Code: Form.CourseHasBound})
		return
	}

	// 绑定
	global.DB.Table("courses").Where("course_id=?", bindCourseRequest.CourseID).Update("teacher_id", bindCourseRequest.TeacherID)
	c.JSON(http.StatusOK, Form.BindCourseResponse{Code: Form.OK})

}

//解绑课程
func UnbindCourse(c *gin.Context) {
	// 创建变量绑定输入
	var unbindCourseRequest Form.UnbindCourseRequest
	// 参数不合法
	if err := c.Bind(&unbindCourseRequest); err != nil || len(unbindCourseRequest.CourseID) == 0 || len(unbindCourseRequest.TeacherID) == 0 {
		c.JSON(http.StatusUnprocessableEntity, Form.UnbindCourseResponse{Code: Form.ParamInvalid})
		return
	}

	// 根据courseId获取对应课程
	var getCourse Form.TCourse
	global.DB.Table("courses").Where("course_id=?", unbindCourseRequest.CourseID).Find(&getCourse)

	// 课程不存在
	if getCourse.CourseID == "" {
		c.JSON(http.StatusUnprocessableEntity, Form.UnbindCourseResponse{Code: Form.CourseNotExisted})
		return
	}

	// 课程未绑定
	if getCourse.TeacherID == "-1" {
		c.JSON(http.StatusUnprocessableEntity, Form.UnbindCourseResponse{Code: Form.CourseNotBind})
		return
	}

	//解绑
	global.DB.Table("courses").Where("course_id=?", unbindCourseRequest.CourseID).Update("teacher_id", "-1")
	c.JSON(http.StatusOK, Form.UnbindCourseResponse{Code: Form.OK})
}

//获取该老师的所有课程
func GetTeacherCourse(c *gin.Context) {
	// 创建变量绑定输入
	var getTeacherCourseRequest Form.GetTeacherCourseRequest
	// 参数不合法
	if err := c.Bind(&getTeacherCourseRequest); err != nil || len(getTeacherCourseRequest.TeacherID) == 0 {
		c.JSON(http.StatusUnprocessableEntity, Form.GetTeacherCourseResponse{
			Code: Form.ParamInvalid,
			Data: struct {
				CourseList []*Form.TCourse
			}{},
		})
		return
	}

	// 根据TeacherId获取对应课程
	var courses []*Form.TCourse
	global.DB.Table("courses").Where("teacher_id=?", getTeacherCourseRequest.TeacherID).Find(&courses)
	if len(courses) == 0 {
		c.JSON(http.StatusOK, Form.GetTeacherCourseResponse{
			Code: Form.CourseNotExisted,
			Data: struct {
				CourseList []*Form.TCourse
			}{},
		})
		return
	}

	c.JSON(http.StatusOK, Form.GetTeacherCourseResponse{
		Code: Form.OK,
		Data: struct {
			CourseList []*Form.TCourse
		}{courses},
	})
}
