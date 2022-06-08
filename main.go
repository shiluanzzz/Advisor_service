package main

import (
	"go.uber.org/zap"
	"service-backend/routes"
	"service-backend/utils"
	"service-backend/utils/cronjob"
	"service-backend/utils/logger"
)

func main() {
	utils.InitDB()
	cronjob.InitCronJob()
	routes.InitRouter()

	//code := service.UpdateAdvisorIndicators(30001)
	//fmt.Println(errmsg.GetErrMsg(code))
	defer func() {
		err := utils.DbConn.Close()
		if err != nil {
			logger.Log.Error("关闭数据库连接错误", zap.Error(err))
		}
	}()
}
