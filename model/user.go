package model

type UserInfo struct {
	Id     int64   `structs:"id"     json:"id"`
	Name   string  `structs:"name"   json:"name"   validate:"min=4,max=20"`
	Phone  string  `structs:"phone"  json:"phone"  validate:"number,len=11"`
	Birth  string  `structs:"birth"  json:"birth"  validate:"datetime=02-01-2006"`
	Gender string  `structs:"gender" json:"gender" validate:"oneof=Female Male 'Not Specified'"`
	Bio    string  `structs:"bio"    json:"bio"    validate:"max=50"`
	About  string  `structs:"about"  json:"about"`
	Coin   float64 `structs:"coin"   json:"coin"   validate:"number"`
}

type Login struct {
	Id       int64  `structs:"id" json:"id" form:"id"`
	Phone    string `structs:"phone" json:"phone" form:"phone" validate:"required,number,len=11"`
	Password string `structs:"password" json:"password" form:"password" validate:"required,min=6,max=12"`
	Token    string `structs:"token"`
}
type ChangePwd struct {
	//长度限制，新旧密码不能相等
	Id          int64
	NewPassword string `form:"newPassword" json:"newPassword" validate:"required,min=6,max=12"`
	OldPassword string `form:"oldPassword" json:"oldPassword" validate:"required,min=6,max=12,nefield=NewPassword"`
}
