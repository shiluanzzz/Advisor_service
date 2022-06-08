package model

type Collection struct {
	Id        int64 `structs:"id" json:"id"`
	UserId    int64 `structs:"user_id" json:"userId"`
	AdvisorId int64 `structs:"advisor_id" json:"advisorId"`
}

type TableID struct {
	Id int64 `form:"id" json:"id" validate:"required,min=0"`
}
