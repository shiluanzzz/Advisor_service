package model

type Service struct {
	AdvisorPhone string  `structs:"advisor_phone" json:"advisor_phone"`
	ServiceName  string  `structs:"service_name" json:"service_name"  validate:"required"`
	ServiceId    int     `structs:"service_id" json:"service_id"`
	Price        float32 `structs:"price" json:"price" validate:"required,number,gte=1,lte=36"`
	Status       int     `structs:"status" json:"status" validate:"required,number,gte=0,lte=1"`
}
type ServiceKind struct {
	Name string `structs:"name" json:"name" validate:"require"`
	Id   int    `structs:"id" json:"id"`
}
