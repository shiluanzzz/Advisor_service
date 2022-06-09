package model

type Bill struct {
	Id         int64    `structs:"id" json:"id,omitempty"`
	UserId     int64    `structs:"user_id" json:"userId,omitempty"`
	AdvisorId  int64    `structs:"advisor_id" json:"advisorId,omitempty"`
	OrderId    int64    `structs:"order_id" json:"orderId"`
	Amount     int64    `structs:"amount" json:"amount,omitempty"`
	Type       billType `structs:"type" json:"type,omitempty"`
	Time       int64    `structs:"time" json:"time"`
	BillType   string   `json:"billType,omitempty"`
	ShowTime   string   `json:"showTime,omitempty"`
	ShowAmount float32  `json:"showAmount,omitempty"`
}

type billType int

const (
	ORDERCOST     billType = 1
	ORDERBACK     billType = 2
	ORDERRUSHCOST billType = 3
	ORDERRUSABACK billType = 4
	ORDERINCOME   billType = 5
)

func (b billType) Name() string {
	switch b {
	case ORDERCOST:
		return "新建订单支出"
	case ORDERRUSHCOST:
		return "订单加急支出"
	case ORDERINCOME:
		return "顾问收入"
	case ORDERBACK:
		return "订单过期退回"
	case ORDERRUSABACK:
		return "订单加急超时退回"
	}
	return ""
}
