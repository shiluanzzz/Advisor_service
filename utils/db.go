package utils

import (
	"database/sql"
	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

var (
	DbConn *sql.DB
)

func InitDB() {
	var err error
	DbConn, err = manager.New(
		DbName, DbUser, DbPassword, DbHost).Set(
		manager.SetCharset("utf8"),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second),
	).Port(DbPort).Open(true)
	DbConn.SetConnMaxLifetime(100 * time.Second)
	DbConn.SetMaxIdleConns(20)
	DbConn.SetMaxOpenConns(20)

	if err != nil {
		log.Fatalln("database build error!", err)
	}
	if err = DbConn.Ping(); err != nil {
		log.Fatalln("database connect error", err)
	}
	//defer DbConn.Close()
}
