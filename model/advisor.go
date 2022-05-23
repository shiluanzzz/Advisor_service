package model

type Advisor struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	Phone         string `json:"phone"`
	Coin          int    `json:"coin"`
	Status        string `json:"status"`
	TotalOrderNum int    `json:"total_order_num"`
	Rank          string `json:"rank"`
	RankNum       string `json:"rank_num"`
}
