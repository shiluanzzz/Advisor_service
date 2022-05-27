package errmsg

import "fmt"

const (
	//SUCCESS  = iota
	SUCCESS = 200
	ERROR   = 400
	// 用户错误状态码
	ERROR_USERPHONE_USED  = 1001
	ERROR_PASSWORD_WORON  = 1002
	ERROR_USER_NOT_EXIST  = 1003
	ERROR_USERNAME_MODIFY = 1004
	// 输入错误
	ERROR_INPUT             = 1005
	ERROR_UPDATE_VALID      = 1006
	ERROR_PHONE_INPUT       = 1007
	ERROR_ADVISOR_NOT_EXIST = 1008
	// 服务器内部错误,SQL编译等
	ERROR_SQL_BUILD     = 2001
	ERROR_MYSQL         = 2002
	ERROR_NOT_IMPLEMENT = 2003
	ErrorGinBind        = 2004
	// TOKEN相关错误
	ERROR_TOKEN_NOT_EXIST   = 3001
	ERROR_TOKEN_TIME_OUT    = 3002
	ERROR_TOKEN_WOKEN_WRONG = 3003
	ERROR_TOKEN_TYPE_WRONG  = 3004
	// service
	ERROR_SERVICE_NOT_EXIST = 4001
	ERROR_SERVICE_EXIST     = 4002

	// order
	ERROR_ORDER_MONEY_INSUFFICIENT = 5001
)

var errMsg = map[int]string{
	SUCCESS: "成功",
	ERROR:   "错误",
	// user
	ERROR_USERPHONE_USED:    "手机号已注册！",
	ERROR_PASSWORD_WORON:    "密码错误",
	ERROR_USER_NOT_EXIST:    "用户不存在",
	ERROR_ADVISOR_NOT_EXIST: "顾问不存在",
	ERROR_USERNAME_MODIFY:   "不允许修改用户名!",
	ERROR_INPUT:             "输入不符合要求!",
	ERROR_UPDATE_VALID:      "不允许直接更新Coin或密码字段",
	ERROR_PHONE_INPUT:       "手机号字段请传字符串,会被识别成float64",
	//服务器内部错误
	ERROR_SQL_BUILD:     "gendry库SQL编译错误",
	ERROR_MYSQL:         "数据库操作错误",
	ERROR_NOT_IMPLEMENT: "接口未开发",
	ErrorGinBind:        "gin框架绑定数据错误",
	// TOKEN相关错误
	ERROR_TOKEN_NOT_EXIST:   "TOKEN不存在",
	ERROR_TOKEN_TIME_OUT:    "TOKEN超时",
	ERROR_TOKEN_WOKEN_WRONG: "TOKEN错误",
	ERROR_TOKEN_TYPE_WRONG:  "TOKEN格式错误",
	// service
	ERROR_SERVICE_NOT_EXIST: "服务不存在",
	ERROR_SERVICE_EXIST:     "服务已存在",
	// order
	ERROR_ORDER_MONEY_INSUFFICIENT: "金币不足!",
}

func GetErrMsg(code int) string {
	if msg, ok := errMsg[code]; ok {
		return msg
	} else {
		return fmt.Sprintf("状态码%v未定义", code)
	}
}
