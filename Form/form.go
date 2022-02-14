package Form

type ErrNo int

const (
	OK                 ErrNo = 0
	ParamInvalid       ErrNo = 1  // 参数不合法
	UserHasExisted     ErrNo = 2  // 该 Username 已存在
	UserHasDeleted     ErrNo = 3  // 用户已删除
	UserNotExisted     ErrNo = 4  // 用户不存在
	WrongPassword      ErrNo = 5  // 密码错误
	LoginRequired      ErrNo = 6  // 用户未登录
	CourseNotAvailable ErrNo = 7  // 课程已满
	CourseHasBound     ErrNo = 8  // 课程已绑定过
	CourseNotBind      ErrNo = 9  // 课程未绑定过
	PermDenied         ErrNo = 10 // 没有操作权限
	StudentNotExisted  ErrNo = 11 // 学生不存在
	CourseNotExisted   ErrNo = 12 // 课程不存在
	StudentHasNoCourse ErrNo = 13 // 学生没有课程
	StudentHasCourse   ErrNo = 14 // 学生有课程

	UnknownError ErrNo = 255 // 未知错误
)

type Member struct {
	UserID   string
	Nickname string   // required，不小于 4 位 不超过 20 位
	Username string   // required，只支持大小写，长度不小于 8 位 不超过 20 位
	Password string   // required，同时包括大小写、数字，长度不少于 8 位 不超过 20 位
	UserType UserType // required, 枚举值
	Deleted  string
}

type UserType int

const (
	Admin   UserType = 1
	Student UserType = 2
	Teacher UserType = 3
)

type Course struct {
	CourseID  string
	Name      string
	TeacherID string
	CourseCap CourseCap
}

type RedisCourse struct {
	Course_ID  string
	Course_Cap int
}

type CourseCap int

type Schedule struct {
	StudentID string
	CourseID  string
}

//成员管理
type TMember struct {
	UserID   string
	Nickname string
	Username string
	UserType UserType
}

//获取成员信息
type GetMemberRequest struct {
	UserID string
}
type GetMemberResponse struct {
	Code ErrNo
	Data TMember
}

//删除
type DeleteMemberRequest struct {
	UserID string
}
type DeleteMemberResponse struct {
	Code ErrNo
}

//批量获取成员
type GetMemberListRequest struct {
	Offset int
	Limit  int
}
type GetMemberListResponse struct {
	Code ErrNo
	Data struct {
		MemberList []TMember
	}
}

//更新
type UpdateMemberRequest struct {
	UserID   string
	Nickname string
}

//创建成员
type CreateMemberRequest struct {
	Nickname string   // required，不小于 4 位 不超过 20 位
	Username string   // required，只支持大小写，长度不小于 8 位 不超过 20 位
	Password string   // required，同时包括大小写、数字，长度不少于 8 位 不超过 20 位
	UserType UserType // required, 枚举值
}
type CreateMemberResponse struct {
	Code ErrNo
	Data struct {
		UserID string // int64 范围
	}
}

//登录
type LoginRequest struct {
	Username string
	Password string
}
type LoginResponse struct {
	Code ErrNo
	Data struct {
		UserID string
	}
}

type LogoutResponse struct {
	Code ErrNo
}
type WhoAmIResponse struct {
	Code ErrNo
	Data TMember
}

/*课程相关的请求模型*/
// 课程
type TCourse struct {
	CourseID  string `gorm:"Column:course_id"`
	TeacherID string `gorm:"Column:teacher_id"`
	Name      string `gorm:"Column:name"`
	CourseCap int    `gorm:"Column:course_cap"`
}

// 创建课程
type CreateCourseRequest struct {
	Name string
	Cap  int
}
type CreateCourseResponse struct {
	Code ErrNo
	Data struct {
		CourseID string
	}
}

// 获取课程
type GetCourseRequest struct {
	CourseID string
}
type GetCourseResponse struct {
	Code ErrNo
	Data TCourse
}

//绑定课程
type BindCourseRequest struct {
	CourseID  string
	TeacherID string
}
type BindCourseResponse struct {
	Code ErrNo
}

//解绑课程
type UnbindCourseRequest struct {
	CourseID  string
	TeacherID string
}
type UnbindCourseResponse struct {
	Code ErrNo
}

// 该教师所有课程
type GetTeacherCourseRequest struct {
	TeacherID string
}
type GetTeacherCourseResponse struct {
	Code ErrNo
	Data struct {
		CourseList []*TCourse
	}
}

// 排课求解器，使老师绑定课程的最优解， 老师有且只能绑定一个课程
// Method： Post
type ScheduleCourseRequest struct {
	TeacherCourseRelationShip map[string][]string // key 为 teacherID , val 为老师期望绑定的课程 courseID 数组
}

type ScheduleCourseResponse struct {
	Code ErrNo
	Data map[string]string // key 为 teacherID , val 为老师最终绑定的课程 courseID
}

// 获取学生课表
type GetStudentCourseRequest struct {
	StudentID string
}

type GetStudentCourseResponse struct {
	Code ErrNo
	Data struct {
		CourseList []TCourse
	}
}

type BookCourseRequest struct {
	StudentID string
	CourseID  string
}
