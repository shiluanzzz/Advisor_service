package service

import (
	"database/sql"
	"service-backend/model"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/tools"
	"time"
)

const BILLTABLE = "bill"

func NewBill(data *model.Bill, tx *sql.Tx) (code int) {
	data.Time = time.Now().Unix()
	defer logger.CommonServiceLog(&code, data, "msg", "新增了一笔流水")
	maps := []map[string]interface{}{tools.Structs2SQLTable(data)}
	code, _ = InsertTableItem(BILLTABLE, maps, tx)
	return
}

// GetBill 用户或者顾问获取自己的账单
func GetBill(id int64, role string) (code int, res []*model.Bill) {
	where := map[string]interface{}{}
	switch role {
	case USERTABLE:
		where["user_id"] = id
	case ADVISORTABLE:
		where["advisor_id"] = id
	default:
		code = errmsg.ErrorTokenRoleNotMatch
		return
	}
	code = GetTableRows2StructByWhere(BILLTABLE, where, []string{"*"}, &res)
	// 完善一些展示信息
	for _, v := range res {
		v.BillType = v.Type.Name()
		v.ShowTime = time.Unix(v.Time, 0).Format("Jan 02,2006 15:04:05")
		v.ShowAmount = tools.ConvertCoinI2F(v.Amount)
	}
	return
}
