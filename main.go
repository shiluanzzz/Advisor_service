package main

import (
	"service/routes"
	"service/utils"
)

func main() {
	utils.InitDB()
	routes.InitRouter()
	defer utils.DbConn.Close()
	//_, res := service.GetAdvisorService("17607175592")
}
