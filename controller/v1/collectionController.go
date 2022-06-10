package v1

import (
	"github.com/gin-gonic/gin"
	"service-backend/model"
	"service-backend/service"
	"service-backend/utils/errmsg"
	"service-backend/utils/validator"
)

func NewCollectionController(ctx *gin.Context) {
	var data *model.TableID
	var code int
	var msg string
	var response *model.Collection

	defer commonControllerDefer(ctx, &code, &msg, &data, &response)
	if err := ctx.ShouldBindQuery(&data); err != nil {
		ginBindError(ctx, err, data)
	}
	if msg, code = validator.Validate(data); code != errmsg.SUCCESS {
		return
	}
	// 顾问是否存在
	if code, _ = service.GetTableItemsById(service.ADVISORTABLE, data.Id, []string{"id"}); code != errmsg.SUCCESS {
		code = errmsg.ErrorAdvisorNotExist
		return
	}
	// 是否已经收藏了?
	if code, _ = service.GetTableItemByWhere(service.COLLECTIONTABLE, map[string]interface{}{"" +
		"user_id": ctx.GetInt64("id"),
		"advisor_id": data.Id,
	}, "id"); code == errmsg.SUCCESS {
		// 查得到说明已经存在了
		code = errmsg.ErrorCollectionExist
		return
	}
	code, response = service.NewCollection(
		&model.Collection{
			UserId:    ctx.GetInt64("id"),
			AdvisorId: data.Id,
		})
	return
}

func GetUserCollectionController(ctx *gin.Context) {
	var code int
	var msg string
	var res []*model.Collection
	defer commonControllerDefer(ctx, &code, &msg, nil, res)
	code, res = service.GetUserCollection(ctx.GetInt64("id"))
	return
}
