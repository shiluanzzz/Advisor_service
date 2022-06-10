package service

import (
	"service-backend/model"
	"service-backend/utils/errmsg"
	"service-backend/utils/tools"
)

const COLLECTIONTABLE = "collection"

func NewCollection(data *model.Collection) (code int, res *model.Collection) {
	maps := []map[string]interface{}{tools.Structs2SQLTable(data)}
	code, data.Id = InsertTableItem(COLLECTIONTABLE, maps)
	return code, data
}

func GetUserCollection(id int64) (code int, res []*model.Collection) {
	code = GetTableRows2StructByWhere(
		COLLECTIONTABLE,
		map[string]interface{}{"user_id": id},
		[]string{"*"},
		&res,
	)
	return errmsg.SUCCESS, res
}
