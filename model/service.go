package model

type Service struct {
	Id          int64   `structs:"id" json:"id"`
	AdvisorId   int64   `structs:"advisor_id" json:"advisor_phone"`
	ServiceName string  `structs:"service_name" json:"service_name"  validate:"required"`
	ServiceId   int     `structs:"service_id" json:"service_id"`
	Price       float32 `structs:"price" json:"price" validate:"required,number,gte=1,lte=36"`
	Status      int     `structs:"status" json:"status" validate:"required,number,gte=0,lte=1"`
}

//type ServiceKind struct {
//	Name string `structs:"name" json:"name" validate:"require"`
//	Id   int    `structs:"id" json:"id"`
//}
type ServiceState struct {
	Status int `structs:"status" form:"status" validate:"required,number,min=0,max=1"`
}

var ServiceKind = map[int]string{
	1: "24h Delivered Video Reading",
	2: "24h Delivered Audio Reading",
	3: "24h Delivered Text Reading",
	4: "Live Text Chat",
}
