package service

import (
	"fmt"
	"go.uber.org/zap"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
	"time"
)

var ADVISORTABLE = "advisor"

func GetAdvisorList(page int) (int, []map[string]interface{}) {
	uPage := uint(page)
	where := map[string]interface{}{
		"status": 1,
		"_limit": []uint{(uPage - 1) * 10, uPage * 10},
	}
	selects := []string{
		"id", "phone", "name", "bio",
	}
	return GetManyTableItemsByWhere(ADVISORTABLE, where, selects)
}

func NewAdvisorAndService(data *model.Login) (code int, id int64) {
	id = -1
	begin, err := utils.DbConn.Begin()
	defer func() {
		if code != errmsg.SUCCESS {
			err = begin.Rollback()
			logger.Log.Error("事务回滚失败!", zap.Error(err))
		}
	}()
	if err != nil {
		return errmsg.ErrorSqlTransError, -1
	}
	code, id = NewUser(ADVISORTABLE, data, begin)
	if code != errmsg.SUCCESS {
		// 创建顾问失败
		return errmsg.ErrorSqlTransError, -1
	}
	code = NewService(id, begin)
	if code != errmsg.SUCCESS {
		// 顾问的服务项创建失败
		return errmsg.ErrorSqlTransError, -1
	}
	err = begin.Commit()
	if err != nil {
		return errmsg.ErrorSqlTransCommitError, -1
	}
	return errmsg.SUCCESS, id
}

// GetAdvisorScore 获取顾问的评分
func GetAdvisorScore(id int64) (code int, score float32) {
	score = 0.0
	where := map[string]interface{}{
		"advisor_id":     id,
		"status":         model.Completed,
		"comment_status": model.Commented,
	}
	selects := []string{"rate"}
	code, data := GetManyTableItemsByWhere(ORDERTABLE, where, selects)
	if code != errmsg.SUCCESS {
		return
	}
	if len(data) != 0 {
		for _, v := range data {
			score += float32(v["rate"].(int64))
		}
		score /= float32(len(data))
	}
	return errmsg.SUCCESS, score
}

//
func GetAdvisorCommentData(id int64) (code int, res []map[string]interface{}) {
	where := map[string]interface{}{
		"advisor_id":     id,
		"status":         model.Completed,
		"comment_status": model.Commented,
	}
	selects := []string{"rate", "create_time", "comment_time", "service_id", "user_id", "comment"}
	if code, res = GetManyTableItemsByWhere(ORDERTABLE, where, selects); code != errmsg.SUCCESS {
		return
	}
	//扩充数据 user_name,service_name 时间格式转化
	for _, v := range res {
		var userNameUint8 interface{}
		if code, userNameUint8 = GetTableItem(USERTABLE, v["user_id"].(int64), "name"); code != errmsg.SUCCESS {
			return
		}
		v["user_name"] = fmt.Sprintf("%s", userNameUint8)
		v["service_name"] = model.ServiceKind[int(v["service_id"].(int64))]
		v["create_show_time"] = time.Unix(v["create_time"].(int64), 0).Format("Jan 02,2006 15:04:05")
		v["comment_show_time"] = time.Unix(v["comment_time"].(int64), 0).Format("Jan 02,2006 15:04:05")
	}
	return
}
