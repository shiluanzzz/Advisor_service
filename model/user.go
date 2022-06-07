package model

import (
	"service-backend/utils/tools"
	"time"
)

type UserInfo struct {
	Name   *string  `structs:"name"   json:"name"   `
	Phone  *string  `structs:"phone"  json:"phone"  `
	Birth  *string  `structs:"birth"  json:"birth"  `
	Gender *int     `structs:"gender" json:"gender" `
	Bio    *string  `structs:"bio"    json:"bio"    `
	About  *string  `structs:"about"  json:"about"  `
	Coin   *float32 `structs:"coin"   json:"coin"   `
}
type User struct {
	Name   string     `structs:"name"   json:"name"   `
	Phone  string     `structs:"phone"  json:"phone"  `
	Birth  string     `structs:"birth"  json:"birth"  `
	Gender GendryEnum `structs:"gender" json:"gender" `
	Bio    string     `structs:"bio"    json:"bio"    `
	About  string     `structs:"about"  json:"about"  `
	Coin   int64      `structs:"coin"  `
	// 一些灵活的展示信息
	CoinShow   float32 `json:"coinShow"`
	GenderShow string  `json:"genderShow"`
	BirthShow  string  `json:"birthShow"`
}

func (u *User) UpdateShow(birthFormat string) {
	u.GenderShow = u.Gender.StatusName()
	u.CoinShow = tools.ConvertCoinI2F(u.Coin)
	t, _ := time.Parse("02-01-2006", u.Birth)
	u.BirthShow = t.Format(birthFormat)
}

type Login struct {
	Id       int64  `structs:"id" json:"id" form:"id"`
	Phone    string `structs:"phone" json:"phone" form:"phone" validate:"required,number,len=11"`
	Password string `structs:"password" json:"password" form:"password" validate:"required,min=6,max=12"`
	Token    string `json:"token"`
}
type ChangePwd struct {
	//长度限制，新旧密码不能相等
	Id          int64
	NewPassword string `form:"newPassword" json:"newPassword" validate:"required,min=6,max=12"`
	OldPassword string `form:"oldPassword" json:"oldPassword" validate:"required,min=6,max=12,nefield=NewPassword"`
}
