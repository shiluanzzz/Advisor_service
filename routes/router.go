package routes

import (
	"github.com/gin-gonic/gin"
	"log"
	"service-backend/middleware"
	"service-backend/service"
	"service-backend/utils"
	"service-backend/utils/logger"
)
import v1 "service-backend/controller/v1"

func InitRouter() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Log())
	UserRouter := r.Group("user")
	{
		UserRouter.POST("/add", v1.NewUserController)
		UserRouter.GET("/login", v1.UserLoginController)
		UserRouter.Use(middleware.JwtToken()).Use(middleware.RoleValidate(service.USERTABLE))
		UserRouter.POST("/pwd", v1.UpdateUserPwd)
		UserRouter.POST("/update", v1.UpdateUserInfoController)
		UserRouter.POST("/get", v1.GetUserInfoController)
	}
	AdvisorRouter := r.Group("advisor")
	{
		AdvisorRouter.POST("/add", v1.NewAdvisorController)
		AdvisorRouter.GET("/login", v1.AdvisorLogin)
		AdvisorRouter.GET("/list/:page", v1.GetAdvisorList)
		AdvisorRouter.GET("/getInfo", v1.GetAdvisorInfo)
		AdvisorRouter.Use(middleware.JwtToken()).Use(middleware.RoleValidate(service.ADVISORTABLE))
		AdvisorRouter.POST("/update", v1.UpdateAdvisorController)
		AdvisorRouter.POST("/pwd", v1.UpdateAdvisorPwd)
		AdvisorRouter.POST("/status", v1.ModifyAdvisorStatus)
	}
	Service := r.Group("service")
	Service.Use(middleware.JwtToken()).Use(middleware.RoleValidate(service.ADVISORTABLE))
	{
		Service.POST("/status", v1.ModifyServiceStatus)
		Service.POST("/price", v1.ModifyServicePrice)
	}
	order := r.Group("order")
	order.Use(middleware.JwtToken()).Use(middleware.RoleValidate(service.ADVISORTABLE))
	{
		// advisor
		order.GET("/list", v1.GetOrderListController)
		order.POST("/reply", v1.OrderReplyController)
		order.GET("/detail", v1.GetOrderDetailController)
	}
	orderUser := r.Group("order")
	orderUser.Use(middleware.JwtToken()).Use(middleware.RoleValidate(service.USERTABLE))
	{
		// user
		orderUser.POST("/add", v1.NewOrderController)
		orderUser.POST("/rush", v1.RushOrderController)
		orderUser.POST("/comment", v1.CommentOrderController)
	}
	collection := r.Group("collection")
	collection.Use(middleware.JwtToken()).Use(middleware.RoleValidate(service.USERTABLE))
	{
		collection.GET("/list", v1.GetUserCollectionController)
		collection.POST("/add", v1.NewCollectionController)
	}
	logger.Log.Info("服务启动")
	err := r.Run(utils.HttpPort)
	if err != nil {
		log.Fatalln("gin框架启动失败", err)
	}
}
