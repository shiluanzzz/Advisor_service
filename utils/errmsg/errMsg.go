package errmsg

import "fmt"

const (
	SUCCESS = iota
	ERROR
	// 用户错误状态码
	ERROR_USERNAME_USED
	ERROR_PASSWORD_WORON
	ERROR_USER_NOT_EXIST
	// 服务器内部错误,SQL编译等
	ERROR_SQL_BUILD
	ERROR_MYSQL
)

var errMsg = map[int]string{
	SUCCESS: "成功",
	ERROR:   "错误",
	// user
	ERROR_USERNAME_USED:  "用户名已存在！",
	ERROR_PASSWORD_WORON: "密码错误",
	ERROR_USER_NOT_EXIST: "用户不存在",
	//服务器内部错误
	ERROR_SQL_BUILD: "服务器内部错误",
	ERROR_MYSQL:     "数据库操作错误",
}

func GetErrMsg(code int) string {
	if msg, ok := errMsg[code]; ok {
		return msg
	} else {
		return fmt.Sprintf("状态码%v未定义", code)
	}
}
