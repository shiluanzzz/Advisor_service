package model

type Order struct {
	Id            int64              `structs:"id" json:"orderId"`
	UserId        int64              `structs:"user_id" json:"userId" validate:"number"`
	ServiceId     int64              `structs:"service_id" json:"serviceId" validate:"required,number"`
	ServiceNameId serviceNameCode    `structs:"service_name_id" json:"serviceNameId"`
	AdvisorId     int64              `structs:"advisor_id" json:"advisorId" validate:"required,number"`
	Situation     string             `structs:"situation" json:"situation" validate:"required,max=3000"`
	Question      string             `structs:"question" json:"question" validate:"required,max=200"`
	Coin          int64              `structs:"coin" json:"coin,omitempty" validate:""`
	RushCoin      int64              `structs:"rush_coin" json:"rushCoin,omitempty"`
	CreateTime    int64              `structs:"create_time" json:"createTime"`
	RushTime      int64              `structs:"rush_time" json:"rushTime,omitempty"`
	Reply         string             `structs:"reply" json:"reply,omitempty" validate:""`
	Status        OrderStatus        `structs:"status" json:"status"`
	Rate          int                `structs:"rate" json:"rate,omitempty" validate:""`
	Comment       string             `structs:"comment" json:"comment,omitempty" validate:""`
	CommentTime   int64              `structs:"comment_time" json:"commentTime,omitempty"`
	CommentStatus OrderCommentStatus `structs:"comment_status" json:"commentStatus,omitempty"`
	// 其他信息
	OrderListOtherInfo
}

// OrderListOtherInfo 订单列表展示的补充信息
type OrderListOtherInfo struct {
	ShowTime          string      `json:"showTime,omitempty"`
	UserName          interface{} `json:"userName,omitempty"`
	ServiceName       string      `json:"serviceName,omitempty"`
	ServiceStatusName string      `json:"serviceStatusName,omitempty"`
}

// OrderInitInfo 新建订单的request
type OrderInitInfo struct {
	OrderId   int64  `json:"orderId"`
	UserId    int64  `structs:"user_id" json:"userId" validate:"number"`
	ServiceId int64  `structs:"service_id" json:"serviceId" validate:"required,number"`
	AdvisorId int64  `structs:"advisor_id" json:"advisorId" validate:"required,number"`
	Situation string `structs:"situation" json:"situation" validate:"required,max=3000"`
	Question  string `structs:"question" json:"question" validate:"required,max=200"`
}
type OrderReply struct {
	Id        int64       `structs:"id" json:"orderId" validate:"min=1"`
	AdvisorId int64       `structs:"advisor_id" json:"advisorId"`
	Reply     string      `structs:"reply" json:"reply" validate:"min=1200,max=5000"`
	Coin      int64       `structs:"coin" json:"coin" validate:""`
	RushCoin  int64       `json:"rushCoin"`
	Status    OrderStatus `json:"status"`
}
type OrderRush struct {
	Id        int64 `structs:"id" json:"orderId"`
	UserId    int64 `structs:"user_id" json:"userId"`
	RushTime  int64 `json:"rushTime"`
	UserMoney int64 `json:"userMoney"`
	RushMoney int64 `json:"rushMoney"`
	Status    int64 `json:"status"`
}
type OrderComment struct {
	Id              int64  `structs:"id" json:"orderId"`
	Comment         string `structs:"comment" json:"comment" validate:"max=300"`
	Rate            int    `structs:"rate" json:"rate" validate:"required,min=0,max=5"`
	UserId          int64  `structs:"user_id" json:"userId"`
	CommentTime     int64  `structs:"comment_time" json:"commentTime"`
	UserName        string `json:"userName"`
	OrderCreateTime int64  `structs:"create_time" json:"orderCreateTime,omitempty"`
	CreateShowTime  string `json:"createShowTime,omitempty"`
	CommentShowTime string `json:"commentShowTime,omitempty"`
}
type CommentStruct struct {
	Id      int64  `structs:"id" json:"orderId"`
	Comment string `structs:"comment" json:"comment" validate:"max=300"`
	Rate    int    `structs:"rate" json:"rate" validate:"required,min=0,max=5"`
}

type OrderDetail struct {
	Order `json:"order"`
	User  `json:"user"`
}
