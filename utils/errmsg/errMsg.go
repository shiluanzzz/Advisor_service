package errmsg

import "fmt"

const (
	//SUCCESS  = iota
	SUCCESS = 200
	ERROR   = 400
	PANIC   = 0
)

// 注册、输入相关错误
const (
	ErrorUserPhoneUsed = iota + 1001
	ErrorPasswordWrong
	ErrorUserNotExist
	ErrorInput
	ErrorUpdateValid
	ErrorAdvisorNotExist
	ErrorMysqlNoRows
)

// 服务器内部错误,SQL编译等
const (
	ErrorSqlBuild = iota + 2001
	ErrorMysql
	ErrorNotImplement
	ErrorGinBind
	ErrorSqlScanner
	ErrorSqlTransError
	ErrorSqlTransCommitError
	ErrorSqlSingleSelectError
)

// TOKEN相关错误
const (
	ErrorTokenNotExist = iota + 3001
	ErrorTokenTimeOut
	ErrorTokenWokenWrong
	ErrorTokenTypeWrong
	ErrorTokenIdNotExist
	ErrorTokenRoleNotExist
	ErrorTokenRoleNotMatch
)

// service
const (
	ErrorServiceNotExist = iota + 4001
	ErrorServiceExist
)

// 业务相关
const (
	ErrorOrderMoneyInsufficient = iota + 5001
	ErrorIdNotMatchWithToken
	ErrorServiceIdNotMatchWithAdvisorID
	ErrorServiceNotOpen
	ErrorServiceName
	ErrorPriceNotMatch
	ErrorServiceStatusNotExist
	ErrorAffectsNotOne
	ErrorRowNotExpect
	ErrorNoResult
	ErrorOrderHasCompleted
	ErrorOrderCantRush
	ErrorOrderIdNotMatchWithAdvisorID
	ErrorOrderIdNotMatchWithUserID
	ErrorOrderCantComment
	ErrorCollectionExist
)

// Cron相关
const (
	ErrorCronAddJob = iota + 6001
	ErrorCronJobId
	ErrorJobStatusNotExpect
	ErrorJobStatusConvert
)

// 缓存相关
const (
	ErrorCacheMarshal = iota + 7001
	ErrorCacheDoSet
	ErrorCacheDoExpire
	ErrorCacheGetBytes
	ErrorCacheDeleteKey
	ErrorCacheUnmarshal
	ErrorCacheKeyNotExist
)

var errMsg = map[int]string{
	SUCCESS: "成功",
	ERROR:   "错误",
	PANIC:   "gin框架错误，请检查数据传参是否正确",
	// user
	ErrorUserPhoneUsed:   "手机号已注册！",
	ErrorPasswordWrong:   "密码错误",
	ErrorUserNotExist:    "用户不存在",
	ErrorAdvisorNotExist: "顾问不存在",
	ErrorInput:           "输入不符合要求!",
	ErrorUpdateValid:     "不允许直接更新Coin或密码字段",
	ErrorMysqlNoRows:     "没有查询到结果",
	//服务器内部错误
	ErrorSqlBuild:             "gendry库SQL编译错误",
	ErrorMysql:                "数据库操作错误",
	ErrorNotImplement:         "接口未开发",
	ErrorGinBind:              "gin框架绑定数据错误,请确认数据格式",
	ErrorSqlScanner:           "gendry库scanner绑定数据格式错误",
	ErrorSqlTransError:        "MySQL事务创建错误",
	ErrorSqlTransCommitError:  "MySQL事务提交失败",
	ErrorSqlSingleSelectError: "MySQL单项数据查询失败",
	// TOKEN相关错误
	ErrorTokenNotExist:     "TOKEN不存在",
	ErrorTokenTimeOut:      "TOKEN超时",
	ErrorTokenWokenWrong:   "TOKEN错误",
	ErrorTokenTypeWrong:    "TOKEN格式错误",
	ErrorTokenIdNotExist:   "TOKEN中的ID不存在",
	ErrorTokenRoleNotExist: "TOKEN中的角色不存在",
	ErrorTokenRoleNotMatch: "TOKEN与业务预期不匹配！",
	// service
	ErrorServiceNotExist: "服务不存在",
	ErrorServiceExist:    "服务已存在",
	// 业务相关
	ErrorOrderMoneyInsufficient:         "金币不足!",
	ErrorIdNotMatchWithToken:            "用户ID与Token中的不匹配",
	ErrorServiceIdNotMatchWithAdvisorID: "服务ID与顾问ID不匹配",
	ErrorServiceNotOpen:                 "顾问的这项服务是关闭状态",
	ErrorServiceName:                    "服务名称与服务器不匹配",
	ErrorPriceNotMatch:                  "服务价格不匹配，顾问可能修改了价格",
	ErrorServiceStatusNotExist:          "服务状态不存在",
	ErrorAffectsNotOne:                  "数据修改设计到多行",
	ErrorRowNotExpect:                   "查询结果与预期不符合",
	ErrorNoResult:                       "无此查询结果,请检查输入",
	ErrorOrderHasCompleted:              "订单已经完成啦！",
	ErrorOrderCantRush:                  "订单不可加急",
	ErrorOrderIdNotMatchWithAdvisorID:   "订单ID与顾问ID不匹配",
	ErrorOrderIdNotMatchWithUserID:      "订单ID与用户ID不匹配",
	ErrorOrderCantComment:               "该订单不可评论",
	ErrorCollectionExist:                "已经收藏了该顾问",
	//定时任务相关
	ErrorCronAddJob:         "cron创建定时任务失败",
	ErrorCronJobId:          "cron定时任务的Id不符合预期",
	ErrorJobStatusNotExpect: "任务与预期状态不符合",
	ErrorJobStatusConvert:   "订单状态转化不符合预期",
	//缓存相关提示信息
	ErrorCacheMarshal:     "缓存数据编码错误",
	ErrorCacheDoSet:       "redis do操作错误",
	ErrorCacheDoExpire:    "redis expire操作错误",
	ErrorCacheGetBytes:    "redis get操作错误",
	ErrorCacheDeleteKey:   "redis 删除缓存数据错误",
	ErrorCacheUnmarshal:   "缓存数据解码错误",
	ErrorCacheKeyNotExist: "缓存key不存在",
}

func GetErrMsg(code int) string {
	if msg, ok := errMsg[code]; ok {
		return msg
	} else {
		return fmt.Sprintf("状态码%v未定义", code)
	}
}
