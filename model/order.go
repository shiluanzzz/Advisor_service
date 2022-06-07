package model

type Order struct {
	Id            int64              `structs:"id" json:"orderId"`
	UserId        int64              `structs:"user_id" json:"userId" validate:"number"`
	ServiceId     int64              `structs:"service_id" json:"serviceId" validate:"required,number"`
	ServiceNameId serviceStatusCode  `structs:"service_name_id" json:"serviceNameId"`
	AdvisorId     int64              `structs:"advisor_id" json:"advisorId" validate:"required,number"`
	Situation     string             `structs:"situation" json:"situation" validate:"required,max=3000"`
	Question      string             `structs:"question" json:"question" validate:"required,max=200"`
	Coin          int64              `structs:"coin"  validate:""`
	RushCoin      int64              `structs:"rush_coin" `
	CreateTime    int64              `structs:"create_time" json:"createTime"`
	RushTime      int64              `structs:"rush_time" json:"rushTime"`
	Reply         string             `structs:"reply" json:"reply" validate:""`
	Status        OrderStatus        `structs:"status"`
	Rate          int                `structs:"rate" json:"rate" validate:""`
	Comment       string             `structs:"comment" json:"comment" validate:""`
	CommentTime   int64              `structs:"comment_time" json:"commentTime"`
	CommentStatus OrderCommentStatus `structs:"comment_status" json:"commentStatus"`
	// 其他信息
	OrderListOtherInfo
	OrderDetailInfo
}

// 订单列表展示的补充信息
type OrderListOtherInfo struct {
	ShowTime    string      `json:"showTime"`
	UserName    interface{} `json:"userName"`
	ServiceName string      `json:"serviceName"`
}

// 订单详情的补充信息
type OrderDetailInfo struct {
	User
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
	Id        int64 `struct:"id" json:"orderId"`
	UserId    int64 `struct:"user_id" json:"userId"`
	RushTime  int64 `json:"rushTime"`
	UserMoney int64
	RushMoney int64
	Status    int64 `json:"status"`
}
type OrderComment struct {
	CommentStruct
	UserId      int64 `struct:"user_id" json:"userId"`
	CommentTime int64 `structs:"comment_time" json:"commentTime"`
}
type CommentStruct struct {
	Id      int64  `struct:"id" json:"orderId"`
	Comment string `structs:"comment" json:"comment" validate:"max=300"`
	Rate    int    `structs:"rate" json:"rate" validate:"required,min=0,max=5"`
}
