package model

type UserInfo struct {
	Name   *string `structs:"name"   json:"name"   validator:"min=4,max=20"`
	Phone  *string `structs:"phone"  json:"phone"  validator:"number,len=11"`
	Birth  *string `structs:"birth"  json:"birth"  validator:"datetime=02-01-2006"`
	Gender *int    `structs:"gender" json:"gender" validator:"min=0,max=2"`
	Bio    *string `structs:"bio"    json:"bio"    validator:"max=50"`
	About  *string `structs:"about"  json:"about"`
	Coin   *int    `structs:"coin"   json:"coin"   validator:"number"`
}

type Login struct {
	Id       int64  `structs:"id" json:"id" form:"id"`
	Phone    string `structs:"phone" json:"phone" form:"phone" validate:"required,number,len=11"`
	Password string `structs:"password" json:"password" form:"password" validate:"required,min=6,max=12"`
	Token    string `structs:"token" json:"token"`
}
type ChangePwd struct {
	//长度限制，新旧密码不能相等
	Id          int64
	NewPassword string `form:"newPassword" json:"newPassword" validate:"required,min=6,max=12"`
	OldPassword string `form:"oldPassword" json:"oldPassword" validate:"required,min=6,max=12,nefield=NewPassword"`
}
