package main

import (
	"go.uber.org/zap"
	"service/routes"
	"service/utils"
	"service/utils/cronjob"
	"service/utils/logger"
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
}
