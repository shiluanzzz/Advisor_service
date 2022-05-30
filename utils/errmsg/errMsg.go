package errmsg

import "fmt"

const (
	//SUCCESS  = iota
	SUCCESS = 200
	ERROR   = 400
	// 用户错误状态码
	ErrorUserphoneUsed  = 1001
	ErrorPasswordWoron  = 1002
	ErrorUserNotExist   = 1003
	ErrorUsernameModify = 1004
	// 输入错误
	ErrorInput           = 1005
	ErrorUpdateValid     = 1006
	ErrorPhoneInput      = 1007
	ErrorAdvisorNotExist = 1008
	// 服务器内部错误,SQL编译等
	ErrorSqlBuild            = 2001
	ErrorMysql               = 2002
	ErrorNotImplement        = 2003
	ErrorGinBind             = 2004
	ErrorSqlScanner          = 2005
	ErrorSqlTransError       = 2006
	ErrorSqlTransCommitError = 2007
	// TOKEN相关错误
	ErrorTokenNotExist     = 3001
	ErrorTokenTimeOut      = 3002
	ErrorTokenWokenWrong   = 3003
	ErrorTokenTypeWrong    = 3004
	ErrorTokenIdNotExist   = 3005
	ErrorTokenRoleNotExist = 3006
	// service
	ErrorServiceNotExist = 4001
	ErrorServiceExist    = 4002

	// order
	ErrorOrderMoneyInsufficient = 5001
)

var errMsg = map[int]string{
	SUCCESS: "成功",
	ERROR:   "错误",
	// user
	ErrorUserphoneUsed:   "手机号已注册！",
	ErrorPasswordWoron:   "密码错误",
	ErrorUserNotExist:    "用户不存在",
	ErrorAdvisorNotExist: "顾问不存在",
	ErrorUsernameModify:  "不允许修改用户名!",
	ErrorInput:           "输入不符合要求!",
	ErrorUpdateValid:     "不允许直接更新Coin或密码字段",
	ErrorPhoneInput:      "手机号字段请传字符串,会被识别成float64",
	//服务器内部错误
	ErrorSqlBuild:            "gendry库SQL编译错误",
	ErrorMysql:               "数据库操作错误",
	ErrorNotImplement:        "接口未开发",
	ErrorGinBind:             "gin框架绑定数据错误,请确认数据格式",
	ErrorSqlScanner:          "gendry库scanner绑定数据格式错误",
	ErrorSqlTransError:       "MySQL事务创建错误",
	ErrorSqlTransCommitError: "MySQL事务提交失败",
	// TOKEN相关错误
	ErrorTokenNotExist:     "TOKEN不存在",
	ErrorTokenTimeOut:      "TOKEN超时",
	ErrorTokenWokenWrong:   "TOKEN错误",
	ErrorTokenTypeWrong:    "TOKEN格式错误",
	ErrorTokenIdNotExist:   "TOKEN中的ID不存在",
	ErrorTokenRoleNotExist: "TOKEN中的角色不存在",
	// service
	ErrorServiceNotExist: "服务不存在",
	ErrorServiceExist:    "服务已存在",
	// order
	ErrorOrderMoneyInsufficient: "金币不足!",
}

func GetErrMsg(code int) string {
	if msg, ok := errMsg[code]; ok {
		return msg
	} else {
		return fmt.Sprintf("状态码%v未定义", code)
	}
}
