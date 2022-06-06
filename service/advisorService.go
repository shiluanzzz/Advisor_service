package service

import (
	"go.uber.org/zap"
	"service/model"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
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
