package model

type Order struct {
	Id            int64  `structs:"id" json:"orderId"`
	UserId        int64  `structs:"user_id" json:"userId" validate:"number"`
	ServiceId     int64  `structs:"service_id" json:"serviceId" validate:"required,number"`
	AdvisorId     int64  `structs:"advisor_id" json:"advisorId" validate:"required,number"`
	Situation     string `structs:"situation" json:"situation" validate:"required,max=3000"`
	Question      string `structs:"question" json:"question" validate:"required,max=200"`
	Coin          int64  `structs:"coin" json:"coin" validate:""`
	RushCoin      int64  `structs:"rush_coin" json:"rushCoin"`
	CreateTime    int64  `structs:"create_time" json:"createTime"`
	RushTime      int64  `structs:"rush_time" json:"rushTime"`
	Reply         string `structs:"reply" json:"reply" validate:""`
	Status        int    `structs:"status"`
	Rate          int    `structs:"rate" json:"rate" validate:""`
	Comment       string `structs:"comment" json:"comment" validate:""`
	CommentTime   int64  `structs:"comment_time" json:"commentTime"`
	CommentStatus int64  `structs:"comment_status" json:"commentStatus"`
}
type OrderInitInfo struct {
	UserId    int64  `structs:"user_id" json:"userId" validate:"number"`
	ServiceId int64  `structs:"service_id" json:"serviceId" validate:"required,number"`
	AdvisorId int64  `structs:"advisor_id" json:"advisorId" validate:"required,number"`
	Situation string `structs:"situation" json:"situation" validate:"required,max=3000"`
	Question  string `structs:"question" json:"question" validate:"required,max=200"`
}
type OrderReply struct {
	Id        int64  `structs:"id" json:"orderId" validate:"min=1"`
	AdvisorId int64  `structs:"advisor_id" json:"advisorId"`
	Reply     string `structs:"reply" json:"reply" validate:"min=1200,max=5000"`
	Coin      int64  `structs:"coin" json:"coin" validate:""`
	RushCoin  int64  `json:"rushCoin"`
	Status    int64  `json:"status"`
}
type OrderRush struct {
	Id       int64 `struct:"id" json:"orderId"`
	UserId   int64 `struct:"user_id" json:"userId"`
	RushTime int64 `json:"rushTime"`
	Status   int64 `json:"status"`
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
