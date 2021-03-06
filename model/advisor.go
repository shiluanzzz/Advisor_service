package model

// AdvisorUpdateInfo 用来动态更新advisor的信息
type AdvisorUpdateInfo struct {
	Name           *string `structs:"name" json:"name" `
	Phone          *string `structs:"phone" json:"phone"`
	WorkExperience *int    `structs:"work_experience" json:"workExperience"`
	Bio            *string `structs:"bio" json:"bio" `
	About          *string `structs:"about" json:"about"`
}

type Advisor struct {
	Id             int64             `json:"id" structs:"id"`
	Phone          string            `json:"phone" structs:"phone"`
	Name           string            `json:"name" structs:"name"`
	Coin           int               `json:"coin,omitempty" structs:"coin"`
	Status         ServiceStatusCode `json:"status" structs:"status"`
	WorkExperience int               `json:"workExperience" structs:"work_experience"`
	Bio            string            `json:"bio" structs:"bio"`
	About          string            `json:"about" structs:"about"`

	//AdvisorIndicators
	TotalOrderNum   int     `json:"totalOrderNum" structs:"total_order_num"`
	TotalCommentNum int     `json:"totalCommentNum" structs:"total_comment_num"`
	Rank            float32 `json:"rank" structs:"rank"`
	OnTime          float32 `json:"onTime" structs:"on_time"`
}

type AdvisorIndicators struct {
	TotalOrderNum   int     `json:"totalOrderNum,omitempty" structs:"total_order_num"`
	TotalCommentNum int     `json:"totalCommentNum,omitempty" structs:"total_comment_num"`
	Rank            float32 `json:"rank,omitempty" structs:"rank"`
	OnTime          float32 `json:"onTime,omitempty" structs:"on_time"`
}

type AdvisorDetailInfo struct {
	ServiceData []*Service      `json:"service"`
	Comments    []*OrderComment `json:"comments"`
	Info        *Advisor        `json:"info"`
}
