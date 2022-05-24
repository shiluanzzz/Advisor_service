package main

import (
	"service/routes"
	"service/utils"
)

func main() {
	utils.InitDB()
	routes.InitRouter()
}
