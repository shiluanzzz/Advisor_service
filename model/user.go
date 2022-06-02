package model

type UserInfo struct {
	Name   *string `structs:"name"   json:"name"   `
	Phone  *string `structs:"phone"  json:"phone"  `
	Birth  *string `structs:"birth"  json:"birth"  `
	Gender *int    `structs:"gender" json:"gender" `
	Bio    *string `structs:"bio"    json:"bio"    `
	About  *string `structs:"about"  json:"about"  `
	Coin   *int    `structs:"coin"   json:"coin"   `
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
