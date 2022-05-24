package model

type User struct {
	Name     string `structs:"name" json:"name"`
	Password string `structs:"password" json:"password"`
	Phone    string `structs:"phone" json:"phone"`
	Birth    string `structs:"birth" json:"birth"`
	Gender   string `structs:"gender" json:"gender"`
	Bio      string `structs:"bio" json:"bio"`
	About    string `structs:"about" json:"about"`
	Coin     int    `structs:"coin" json:"coin"`
}
