package model

import "service/utils/errmsg"

type Service struct {
	Id            int64  `structs:"id" json:"id"`
	AdvisorId     int64  `structs:"advisor_id" json:"advisorPhone"`
	ServiceName   string `structs:"service_name" json:"serviceName"  validate:"required"`
	ServiceNameId int    `structs:"service_name_id" json:"serviceNameId"`
	Price         int64  `structs:"price" json:"price" validate:"required,number,gte=1,lte=36"`
	Status        int    `structs:"status" json:"status" validate:"required,number,gte=0,lte=1"`
}

type ServiceState struct {
	Id     int64 `json:"advisorId"`
	Status int   `structs:"status" form:"status" validate:"required,number,min=0,max=1"`
}

var ServiceKind = map[int]string{
	1: "24h Delivered Video Reading",
	2: "24h Delivered Audio Reading",
	3: "24h Delivered Text Reading",
	4: "Live Text Chat",
}

func GetServiceNameById(id int) (int, string) {
	if name := ServiceKind[id]; name != "" {
		return errmsg.SUCCESS, name
	} else {
		return errmsg.ErrorServiceNotExist, name
	}
}

// ServicePrice 修改服务价格
type ServicePrice struct {
	AdvisorId int64   `json:"advisorId"`
	ServiceID int     `form:"serviceNameId" json:"serviceNameId" validate:"required,number,lte=4"`
	Price     float32 `form:"price" json:"price" validate:"required,number,gte=1,lte=36"`
}
