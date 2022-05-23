package model

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Birth    string `json:"birth"`
	Gender   string `json:"gender"`
	Bio      string `json:"bio"`
	About    string `json:"about"`
	Coin     int    `json:"coin"`
}
