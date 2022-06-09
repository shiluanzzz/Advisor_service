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

	defer func() {
		err := utils.DbConn.Close()
		if err != nil {
			logger.Log.Error("关闭数据库连接错误", zap.Error(err))
		}
	}()
	select {}
}
