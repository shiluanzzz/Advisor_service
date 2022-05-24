package model

type User struct {
	Name     string `structs:"name" json:"name" validate:"required,min=4,max=20"`
	Password string `structs:"password" json:"password" validate:"required,min=6,max=12"`
	Phone    string `structs:"phone" json:"phone" validate:"required,number,len=11"`
	Birth    string `structs:"birth" json:"birth" validate:"datetime=02-01-2006"`
	Gender   string `structs:"gender" json:"gender" validate:"required,oneof=Female Male 'Not Specified'"`
	Bio      string `structs:"bio" json:"bio" validate:"max=50"`
	About    string `structs:"about" json:"about"`
	Coin     int    `structs:"coin" json:"coin" validate:"number"`
}
