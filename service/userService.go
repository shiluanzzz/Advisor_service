package service

import (
	"database/sql"
	"fmt"
	"service-backend/model"
	"service-backend/utils/errmsg"
	"service-backend/utils/tools"
)

var USERTABLE = "user"

// NewRole 新增用户或者顾问
func NewRole(table string, user *model.Login, tx *sql.Tx) (code int, id int64) {

	maps := []map[string]interface{}{tools.Structs2SQLTable(user)}
	code, id = InsertTableItem(table, maps, tx)
	return
}

// GetUser 对查询用户信息的方法再次封装，补充消息
func GetUser(id int64) (code int, res *model.User) {
	code = GetTableRows2StructByWhere(
		USERTABLE,
		map[string]interface{}{"id": id},
		[]string{"*"},
		&res,
	)
	return code, res
}
func GetUserName(UserId int64) (code int, res string) {
	var userNameUint8 interface{}
	if code, userNameUint8 = GetTableItemById(USERTABLE, UserId, "name"); code != errmsg.SUCCESS {
		return
	}
	res = fmt.Sprintf("%s", userNameUint8)
	return errmsg.SUCCESS, res
}
