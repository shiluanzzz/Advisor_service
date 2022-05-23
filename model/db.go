package model

import (
	"database/sql"
	"github.com/didi/gendry/manager"
	"log"
	"service/utils"
	"time"
)

var (
	db *sql.DB
)

func initDB() {
	var err error
	db, err = manager.New(
		utils.DbName, utils.DbUser, utils.DbPassword, utils.DbHost).Set(
		manager.SetCharset("utf8"),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second),
	).Port(utils.DbPort).Open(true)
	if err != nil {
		log.Fatalln("database connect error!")
	}
}
