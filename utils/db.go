package utils

import (
	"database/sql"
	"github.com/didi/gendry/manager"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gomodule/redigo/redis"
	"log"
	"service-backend/utils/setting"
	"time"
)

var (
	DbConn    *sql.DB
	RedisConn *redis.Pool
)

func InitDB() {
	var err error
	DbConn, err = manager.New(
		setting.DB.DbName,
		setting.DB.DbUser,
		setting.DB.DbPassword,
		setting.DB.DbHost).Set(
		manager.SetCharset("utf8"),
		manager.SetTimeout(1*time.Second),
		manager.SetReadTimeout(1*time.Second),
	).Port(setting.DB.DbPort).Open(true)
	if err != nil {
		log.Fatalln("database build error!", err)
	}
	//defer DbConn.Close()
}
func InitRedis() {

}
