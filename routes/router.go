package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"service/utils"
)
import v1 "service/controller/v1"

func InitRouter() {
	r := gin.Default()
	UserRouter := r.Group("user")
	{
		UserRouter.POST("/add", v1.NewUserController)
	}

	err := r.Run(utils.HttpPort)
	if err != nil {
		log.Fatalln("gin框架启动失败", err)
	}
}
