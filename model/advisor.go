package model

type Advisor struct {
	Name           string  `structs:"name" json:"name" validate:"required,min=4,max=20"`
	Password       string  `structs:"password" json:"password" validate:"required,min=6,max=12"`
	Phone          string  `structs:"phone" json:"phone" validate:"required,number,len=11"`
	Coin           int     `structs:"coin" json:"coin" validate:"number"`
	TotalOrderNum  int     `structs:"total_order_num" json:"total_order_num" validate:"number"`
	Status         string  `structs:"status" json:"status" validate:"required,oneof=open close"`
	Rank           float64 `structs:"rank" json:"rank" validate:"number"`
	RankNum        int     `structs:"rank_num" json:"rank_num" validate:"number"`
	WorkExperience int     `structs:"work_experience" json:"work_experience" validate:"number,gte=0"`
	Bio            string  `structs:"bio" json:"bio" validate:"max=50"`
	About          string  `structs:"about" json:"about"`
}
