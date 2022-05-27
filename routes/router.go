package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"service/middleware"
	"service/utils"
	"service/utils/logger"
)
import v1 "service/controller/v1"

func InitRouter() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Log())
	UserRouter := r.Group("user")
	{
		UserRouter.POST("/add", v1.NewUser)
		UserRouter.GET("/login", v1.UserLogin)
		UserRouter.Use(middleware.JwtToken())
		UserRouter.POST("/pwd", v1.UpdateUserPwd)
		UserRouter.POST("/update", v1.UpdateUserInfoController)
		UserRouter.POST("/get", v1.GetUserInfo)
	}
	AdvisorRouter := r.Group("advisor")
	{
		AdvisorRouter.POST("/add", v1.NewAdvisorController)
		AdvisorRouter.GET("/login", v1.AdvisorLogin)
		AdvisorRouter.GET("/list/:page", v1.GetAdvisorList)
		AdvisorRouter.Use(middleware.JwtToken())
		AdvisorRouter.POST("/update", v1.UpdateAdvisorController)
		AdvisorRouter.POST("/pwd", v1.UpdateAdvisorPwd)
		AdvisorRouter.GET("/getInfo", v1.GetAdvisorInfo)
		AdvisorRouter.POST("/status", v1.ModifyAdvisorStatus)
	}
	Service := r.Group("service")
	Service.Use(middleware.JwtToken())
	{
		Service.POST("/add", v1.NewService)
		Service.POST("/status", v1.ModifyServiceStatus)
		Service.POST("/price", v1.ModifyServicePrice)
	}
	logger.Log.Info("服务启动")
	err := r.Run(utils.HttpPort)
	if err != nil {
		log.Fatalln("gin框架启动失败", err)
	}
}
