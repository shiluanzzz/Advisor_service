package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"service/middleware"
	"service/utils"
)
import v1 "service/controller/v1"

func InitRouter() {
	r := gin.Default()
	UserRouter := r.Group("user")
	{
		UserRouter.POST("/add", v1.NewUserController)
		UserRouter.GET("/login", v1.UserLogin)
		UserRouter.POST("/changepwd", v1.UpdateUserPwd)
		UserRouter.Use(middleware.JwtToken()).POST("/update", v1.UpdateUserInfoController)
		UserRouter.Use(middleware.JwtToken()).POST("/get", v1.GetUserInfo)
	}
	err := r.Run(utils.HttpPort)
	if err != nil {
		log.Fatalln("gin框架启动失败", err)
	}
}
