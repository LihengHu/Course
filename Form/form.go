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

type GetMemberListResponse struct {
	Code ErrNo
	Data struct {
		MemberList []TMember
	}
}
type TMember struct {
	UserID   string
	Nickname string
	Username string
	UserType UserType
}

type GetMemberResponse struct {
	Code ErrNo
	Data TMember
}

type DeleteMemberResponse struct {
	Code ErrNo
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

type CreateMemberResponse struct {
	Code ErrNo
	Data struct {
		UserID string // int64 范围
	}
}

/*课程相关的请求模型*/
// 课程
type TCourse struct {
	CourseID  string `gorm:"Column:course_id"`
	TeacherID string `gorm:"Column:teacher_id"`
	Name      string `gorm:"Column:name"`
	Cap       int    `gorm:"Column:cap"`
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
