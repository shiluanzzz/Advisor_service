package model

type Service struct {
	Id            int64             `structs:"id" json:"id"`
	AdvisorId     int64             `structs:"advisor_id" json:"advisorPhone"`
	ServiceName   string            `structs:"service_name" json:"serviceName"  validate:"required"`
	ServiceNameId serviceStatusCode `structs:"service_name_id" json:"serviceNameId"`
	Price         int64             `structs:"price" json:"price" validate:"required,number,gte=1,lte=36"`
	Status        int               `structs:"status" json:"status" validate:"required,number,gte=0,lte=1"`
}

type serviceStatusCode int

const (
	VideoReading serviceStatusCode = 1
	AudioReading serviceStatusCode = 2
	TextReading  serviceStatusCode = 3
	TextChat     serviceStatusCode = 4
)

func (s serviceStatusCode) StatusName() string {
	switch s {
	case VideoReading:
		return "24h Delivered Video Reading"
	case AudioReading:
		return "24h Delivered Audio Reading"
	case TextReading:
		return "24h Delivered Text Reading"
	case TextChat:
		return "Live Text Chat"
	}
	return ""
}

var ServiceKind = map[serviceStatusCode]string{
	VideoReading: "24h Delivered Video Reading",
	AudioReading: "24h Delivered Audio Reading",
	TextReading:  "24h Delivered Text Reading",
	TextChat:     "Live Text Chat",
}

// ServicePrice 修改服务价格
type ServicePrice struct {
	AdvisorId     int64             `json:"advisorId"`
	ServiceNameId serviceStatusCode `form:"serviceNameId" json:"serviceNameId" validate:"required,number,lte=4"`
	Price         float32           `form:"price" json:"price" validate:"required,number,gte=1,lte=36"`
}
type ServiceState struct {
	AdvisorId     int64             `json:"advisorId"`
	ServiceNameId serviceStatusCode `form:"serviceNameId" json:"serviceNameId"`
	Status        int               `form:"status" structs:"status" form:"status" validate:"number,min=0,max=1"`
}
