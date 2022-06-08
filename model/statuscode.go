package model

// OrderStatus 订单状态
type OrderStatus int

const (
	Pending   OrderStatus = 0
	Rush      OrderStatus = 1
	Expired   OrderStatus = 2
	Completed OrderStatus = 3
)

func (o OrderStatus) StatusName() string {
	switch o {
	case Pending:
		return "Pending"
	case Rush:
		return "Rush"
	case Expired:
		return "Expired"
	case Completed:
		return "Completed"
	default:
		return ""
	}
}
func (o OrderStatus) CanReply() bool {
	return o == Rush || o == Pending
}
func (o OrderStatus) CanRush() bool {
	return o == Pending
}

// OrderCommentStatus 订单回复状态
type OrderCommentStatus int

// 订单回复状态枚举
const (
	NotComment OrderCommentStatus = 0
	Commented  OrderCommentStatus = 1
)

// GendryEnum 用户性别枚举
type GendryEnum int

const (
	Unknown GendryEnum = 0
	Male    GendryEnum = 1
	Female  GendryEnum = 2
)

func (g GendryEnum) StatusName() string {
	switch g {
	case Male:
		return "Male"
	case Female:
		return "Female"
	case Unknown:
		return "Not Specified"
	default:
		return ""
	}
}

// ServiceStatusCode 顾问服务状态枚举
type ServiceStatusCode int

const (
	AdvisorServiceNotOpen ServiceStatusCode = 0
	AdvisorServiceOpen    ServiceStatusCode = 1
)

// LogType 日志分层
type LogType int

const (
	ControllerLog LogType = 1
	ServiceLog    LogType = 2
)

func (l LogType) StatusName() string {
	switch l {
	case ControllerLog:
		return "controller"
	case ServiceLog:
		return "service"
	}
	return "unknown"
}
