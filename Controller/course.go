package Controller

import (
	"Course/Form"
	"Course/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

//创建课程
func CreateCourse(c *gin.Context) {
	var createCourseRequest Form.CreateCourseRequest
	// 参数有误
	if err := c.Bind(&createCourseRequest); err != nil || len(createCourseRequest.Name) == 0 || createCourseRequest.Cap <= 0 {
		c.JSON(http.StatusOK, Form.CreateCourseResponse{
			Code: Form.ParamInvalid,
			Data: struct {
				CourseID string
			}{},
		})
		return
	}
	var exist int64
	global.DB.Table("courses").Where("name=?", createCourseRequest.Name).Count(&exist)
	// 该用户名已存在
	if exist != 0 {
		c.JSON(http.StatusOK, Form.CreateCourseResponse{
			Code: Form.UnknownError,
			Data: struct {
				CourseID string
			}{},
		})
		return
	}

	// 获取最新的记录
	var size int64
	global.DB.Table("courses").Count(&size)
	var firstCourse Form.TCourse
	global.DB.Table("courses").Offset(int(size - 1)).Limit(1).Find(&firstCourse)
	oldId, _ := strconv.ParseInt(firstCourse.CourseID, 10, 64)
	fmt.Println(oldId + 1)

	newCourse := Form.TCourse{
		CourseID:  strconv.FormatInt(oldId+1, 10),
		Name:      createCourseRequest.Name,
		TeacherID: "-1",
		CourseCap: createCourseRequest.Cap,
	}

	///*redis*/
	//err := global.RDB.Del(global.CTX, "courses").Err()
	//if err != nil {
	//	panic(err)
	//}
	global.DB.Table("courses").Create(&newCourse)
	global.LOG.Info(
		"Create Course",
		zap.Any("Course", newCourse),
	)

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
		c.JSON(http.StatusOK, Form.GetCourseResponse{
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
		c.JSON(http.StatusOK, Form.GetCourseResponse{
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
		c.JSON(http.StatusOK, Form.BindCourseResponse{Code: Form.ParamInvalid})
		return
	}

	// 根据courseId获取对应课程
	var getCourse Form.TCourse
	global.DB.Table("courses").Where("course_id=?", bindCourseRequest.CourseID).Find(&getCourse)

	// 课程不存在
	if getCourse.CourseID == "" {
		c.JSON(http.StatusOK, Form.BindCourseResponse{Code: Form.CourseNotExisted})
		return
	}

	// 课程已绑定
	if getCourse.TeacherID != "-1" {
		c.JSON(http.StatusOK, Form.BindCourseResponse{Code: Form.CourseHasBound})
		return
	}

	// 绑定
	global.DB.Table("courses").Where("course_id=?", bindCourseRequest.CourseID).Update("teacher_id", bindCourseRequest.TeacherID)
	global.LOG.Info(
		"Bind Course",
		zap.Any("Course", bindCourseRequest),
	)
	c.JSON(http.StatusOK, Form.BindCourseResponse{Code: Form.OK})

}

//解绑课程
func UnbindCourse(c *gin.Context) {
	// 创建变量绑定输入
	var unbindCourseRequest Form.UnbindCourseRequest
	// 参数不合法
	if err := c.Bind(&unbindCourseRequest); err != nil || len(unbindCourseRequest.CourseID) == 0 || len(unbindCourseRequest.TeacherID) == 0 {
		c.JSON(http.StatusOK, Form.UnbindCourseResponse{Code: Form.ParamInvalid})
		return
	}

	// 根据courseId获取对应课程
	var getCourse Form.TCourse
	global.DB.Table("courses").Where("course_id=?", unbindCourseRequest.CourseID).Find(&getCourse)

	// 课程不存在
	if getCourse.CourseID == "" {
		c.JSON(http.StatusOK, Form.UnbindCourseResponse{Code: Form.CourseNotExisted})
		return
	}

	// 课程未绑定、课程已绑定的TeacherID与传入的TeacherID不一致
	if getCourse.TeacherID == "-1" || getCourse.TeacherID != unbindCourseRequest.TeacherID {
		c.JSON(http.StatusOK, Form.UnbindCourseResponse{Code: Form.CourseNotBind})
		return
	}

	//解绑
	global.DB.Table("courses").Where("course_id=?", unbindCourseRequest.CourseID).Update("teacher_id", "-1")
	global.LOG.Info(
		"Unbind Course",
		zap.Any("Course", unbindCourseRequest),
	)
	c.JSON(http.StatusOK, Form.UnbindCourseResponse{Code: Form.OK})
}

//获取该老师的所有课程
func GetTeacherCourse(c *gin.Context) {
	// 创建变量绑定输入
	var getTeacherCourseRequest Form.GetTeacherCourseRequest
	// 参数不合法
	if err := c.Bind(&getTeacherCourseRequest); err != nil || len(getTeacherCourseRequest.TeacherID) == 0 {
		c.JSON(http.StatusOK, Form.GetTeacherCourseResponse{
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
			Code: Form.OK,
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

//排课求解器
func Schedule(c *gin.Context) {
	// 创建变量绑定输入
	var scheduleCourseRequest Form.ScheduleCourseRequest
	// 参数不合法
	if err := c.Bind(&scheduleCourseRequest); err != nil || len(scheduleCourseRequest.TeacherCourseRelationShip) == 0 {
		c.JSON(http.StatusOK, Form.ScheduleCourseResponse{
			Code: Form.ParamInvalid,
			Data: nil,
		})
		return
	}
	// key是courseID，val是TeacherID，记录该课程已分配的教师号
	whoPickCourse := make(map[string]string)

	for k, _ := range scheduleCourseRequest.TeacherCourseRelationShip {
		// key是courseID，val是记录本次是否访问过这个courseID(一轮dfs不能对同一个courseID重新分配多次，否则会死循环)，每次dfs开始的时候重置vis
		vis := make(map[string]bool)
		dfsSchedule(k, scheduleCourseRequest.TeacherCourseRelationShip, whoPickCourse, vis)
	}
	// 交换whoPickCourse的key和value就是结果
	res := make(map[string]string)
	for k, v := range whoPickCourse {
		res[v] = k
	}

	c.JSON(http.StatusOK, Form.ScheduleCourseResponse{
		Code: Form.OK,
		Data: res,
	})
}

func dfsSchedule(TeacherID string, TeacherCourseRelationShip map[string][]string, whoPickCourse map[string]string, vis map[string]bool) (res bool) {
	// 得到这个教师可以绑定的CourseID数组
	arr := TeacherCourseRelationShip[TeacherID]
	// 遍历CourseID数组
	for _, v := range arr {
		// 如果这个CourseID在这轮dfs没被访问过，则可以**尝试**把这个CourseID与当前的TeacherID绑定
		_, ok := vis[v]
		if !ok {
			// 标记这轮已访问过该CourseID
			vis[v] = true
			// 查询这个CourseID此前是否已经被别人绑定了，
			pickedTeacherID, ok := whoPickCourse[v]
			// 如果已被绑定，则尝试让绑定了这个CourseID的TeacherID换一个（dfs）
			//	未绑定、更换成功都可以将这个CourseID与当前TeacherID绑定，然后返回true
			if !ok || dfsSchedule(pickedTeacherID, TeacherCourseRelationShip, whoPickCourse, vis) {
				whoPickCourse[v] = TeacherID
				res = true
				return
			}
		}
	}
	// 遍历到底该TeacherID都未能绑定一个CourseID，返回false
	res = false
	return
}

//获取学生课程
func StudentCourse(c *gin.Context) {
	// 创建变量绑定输入
	var getStudentCourseRequest Form.GetStudentCourseRequest
	// 参数不合法
	if err := c.Bind(&getStudentCourseRequest); err != nil || len(getStudentCourseRequest.StudentID) == 0 {
		c.JSON(http.StatusOK, Form.GetStudentCourseResponse{
			Code: Form.OK,
			Data: struct {
				CourseList []Form.TCourse
			}{},
		})
		return
	}

	var student Form.TMember
	global.DB.Table("members").Where("user_id=? and deleted=?", getStudentCourseRequest.StudentID, 0).Find(&student)
	// 学生不存在
	if student.UserID != getStudentCourseRequest.StudentID || student.UserType != Form.Student {
		c.JSON(http.StatusOK, Form.GetStudentCourseResponse{
			Code: Form.StudentNotExisted,
			Data: struct {
				CourseList []Form.TCourse
			}{},
		})
		return
	}

	var courseIDList []string
	global.DB.Table("schedules").Where("student_id=?", getStudentCourseRequest.StudentID).Select("course_id").Find(&courseIDList)
	// 学生没课程
	if len(courseIDList) == 0 {
		c.JSON(http.StatusOK, Form.GetStudentCourseResponse{
			Code: Form.StudentHasNoCourse,
			Data: struct {
				CourseList []Form.TCourse
			}{},
		})
		return
	}
	var res []Form.TCourse
	for _, courseId := range courseIDList {
		var getCourse Form.TCourse
		global.DB.Table("courses").Where("course_id=?", courseId).Find(&getCourse)
		res = append(res, getCourse)
	}
	c.JSON(http.StatusOK, Form.GetStudentCourseResponse{
		Code: Form.OK,
		Data: struct {
			CourseList []Form.TCourse
		}{res},
	})
}
