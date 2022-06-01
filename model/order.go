package model

type Order struct {
	Id         int64   `structs:"id" json:"id"`
	UserId     int64   `structs:"user_id" json:"userId" validate:"required,number"`
	ServiceId  int64   `structs:"service_id" json:"serviceId" validate:"required,number"`
	AdvisorId  int64   `structs:"advisor_id" json:"advisorId" validate:"required,number"`
	Situation  string  `structs:"situation" json:"situation" validate:"required,max=3000"`
	Question   string  `structs:"question" json:"question" validate:"required,max=200"`
	Coin       float32 `structs:"coin" json:"coin" validate:""`
	RushCoin   float32 `structs:"rush_coin" json:"rushCoin"`
	CreateTime int64   `structs:"create_time" json:"createTime"`
	RushTime   int64   `structs:"rush_time" json:"rushTime"`
	Reply      string  `structs:"reply" json:"reply" validate:""`
	Rate       float32 `structs:"rate" json:"rate" validate:""`
	Comment    string  `structs:"comment" json:"comments" validate:""`
	Status     int     `structs:"status"`
}
type OrderReply struct {
	Id        int64   `structs:"id" json:"orderId" validate:"min=1"`
	AdvisorId int64   `structs:"advisor_id" json:"advisorId"`
	Reply     string  `structs:"reply" json:"reply" validate:"min=1200,max=5000"`
	Coin      float32 `structs:"coin" json:"coin" validate:""`
	RushCoin  float32 `json:"rushCoin"`
	Status    int64   `json:"status"`
}
type OrderRush struct {
	Id       int64 `struct:"id" json:"orderId"`
	UserId   int64 `struct:"user_id" json:"userId"`
	RushTime int64 `json:"rushTime"`
	Status   int64 `json:"status"`
}
