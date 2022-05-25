package errmsg

import "fmt"

const (
	SUCCESS = 200
	ERROR   = 400
	// 用户错误状态码
	ERROR_USERPHONE_USED  = 1001
	ERROR_PASSWORD_WORON  = 1002
	ERROR_USER_NOT_EXIST  = 1003
	ERROR_USERNAME_MODIFY = 1004
	// 输入错误
	ERROR_INPUT = 1005
	// 服务器内部错误,SQL编译等
	ERROR_SQL_BUILD = 2001
	ERROR_MYSQL     = 2002
	// TOKEN相关错误
	ERROR_TOKEN_NOT_EXIST   = 3001
	ERROR_TOKEN_TIME_OUT    = 3002
	ERROR_TOKEN_WOKEN_WRONG = 3003
	ERROR_TOKEN_TYPE_WRONG  = 3004
)

var errMsg = map[int]string{
	SUCCESS: "成功",
	ERROR:   "错误",
	// user
	ERROR_USERPHONE_USED:  "手机号已注册！",
	ERROR_PASSWORD_WORON:  "密码错误",
	ERROR_USER_NOT_EXIST:  "用户不存在",
	ERROR_USERNAME_MODIFY: "不允许修改用户名!",
	ERROR_INPUT:           "输入不符合要求!",

	//服务器内部错误
	ERROR_SQL_BUILD: "服务器内部错误",
	ERROR_MYSQL:     "数据库操作错误",

	// TOKEN相关错误
	ERROR_TOKEN_NOT_EXIST:   "TOKEN不存在",
	ERROR_TOKEN_TIME_OUT:    "TOKEN超时",
	ERROR_TOKEN_WOKEN_WRONG: "TOKEN错误",
	ERROR_TOKEN_TYPE_WRONG:  "TOKEN格式错误",
}

func GetErrMsg(code int) string {
	if msg, ok := errMsg[code]; ok {
		return msg
	} else {
		return fmt.Sprintf("状态码%v未定义", code)
	}
}
