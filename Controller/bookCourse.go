package Controller

import (
	"Course/Form"
	"Course/global"
	"Course/types"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func BookCourse(c *gin.Context) {
	StudentID := c.PostForm("StudentID")
	CourseID := c.PostForm("CourseID")
	//清除脏数据时取消注释
	//err := global.RDB.Del(global.CTX, "schedules").Err()
	//if err != nil {
	//	panic(err)
	//}
	//err = global.RDB.Del(global.CTX, "courses").Err()
	//if err != nil {
	//	panic(err)
	//}

	// 刷新members缓存并检验学生是否存在
	count, err := global.RDB.SCard(global.CTX, "members").Result()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		global.MutexMembers.Lock()
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
		global.MutexMembers.Unlock()
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

	// 刷新courses缓存并检验课程是否存在
	count, err = global.RDB.HLen(global.CTX, "courses").Result()
	if err != nil {
		panic(err)
	}
	if count == 0 {
		global.MutexCourses.Lock()
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
		global.MutexCourses.Unlock()
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

	// 刷新选课缓存
	countSchedule, err := global.RDB.SCard(global.CTX, "schedules").Result()
	if err != nil {
		panic(err)
	}
	if countSchedule == 0 {
		global.MutexSchedules.Lock()
		countSchedule, err := global.RDB.SCard(global.CTX, "schedules").Result()
		if err != nil {
			panic(err)
		}
		if countSchedule == 0 {
			var schedules []*Form.Schedule
			global.DB.Table("schedules").Select("Course_ID,Student_ID").Find(&schedules)
			if len(schedules) > 0 {
				var redisSchedules []string
				for _, schedule := range schedules {
					redisSchedules = append(redisSchedules, schedule.CourseID+"_"+schedule.StudentID)
				}
				err := global.RDB.SAdd(global.CTX, "schedules", redisSchedules).Err()
				if err != nil {
					panic(err)
				}
			} else {
				err := global.RDB.SAdd(global.CTX, "schedules", "nothing").Err()
				if err != nil {
					panic(err)
				}
			}
		}
		global.MutexSchedules.Unlock()
	}

	// lua脚本原子更新redis
	scheduleID := CourseID + "_" + StudentID
	const LuaScript = `
		local scheduleKey = KEYS[1]
		local courseKey = KEYS[2]
        local scheduleID = KEYS[3]
        local courseID = KEYS[4]
		local cap = tonumber(redis.call('HGet', courseKey, courseID))
		if cap <= 0 then
			return 1
		end
		local isBooked = redis.call('SIsMember', scheduleKey, scheduleID)
		if isBooked ~= 0 then
			return 2
		end
		redis.call('SAdd', scheduleKey, scheduleID)
		redis.call('HIncrBy', courseKey, courseID, -1)
		return 0
	`
	lua := redis.NewScript(LuaScript)
	val, err := lua.Run(global.CTX, global.RDB, []string{"schedules", "courses", scheduleID, CourseID}).Int()
	if err != nil {
		panic(err)
	}
	if val == 1 {
		c.JSON(200, gin.H{
			"Code": Form.CourseNotAvailable,
		})
		return
	} else if val == 2 {
		c.JSON(200, gin.H{
			"Code": Form.StudentHasCourse,
		})
		return
	}

	// 抢课成功，开启事务写入数据库
	tx := global.DB.Begin()
	err = tx.Table("courses").Where("Course_ID = ?", CourseID).Update("Course_Cap", gorm.Expr("Course_Cap - 1")).Error
	if err != nil {
		tx.Rollback()
		c.JSON(200, gin.H{
			"Code": Form.UnknownError,
		})
		return
	}
	u1 := Form.Schedule{StudentID, CourseID}
	err = tx.Create(&u1).Error
	if err != nil {
		tx.Rollback()
		c.JSON(200, gin.H{
			"Code": Form.UnknownError,
		})
		return
	}
	tx.Commit()
	global.LOG.Info(
		"Book Course",
		zap.String("CourseID", CourseID),
		zap.String("UserID", StudentID),
	)
	c.JSON(200, types.BookCourseResponse{Code: 0})
}
