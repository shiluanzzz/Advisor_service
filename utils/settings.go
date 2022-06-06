package utils

import "log"
import "gopkg.in/ini.v1"
import "github.com/didi/gendry/scanner"

var (
	Db         string
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string

	HttpPort string
	AppMode  string
	JwtKey   string

	LoggerMode string
	InfoLog    string
	ErrorLog   string
	WarnLog    string

	RushOrderCost           float32
	RushOrder2PendingTime   int64
	PendingOrder2ExpireTime int64
	CoinBase                int64
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatalln("read config file error", err)
	}
	Setting()
	LoadDBConfig(file)
	LoadServer(file)
	LoadLogger(file)
	LoadServiceLogger(file)
}

// Setting 一些库的相关设置
func Setting() {
	// 用于scanner的字段反射
	scanner.SetTagName("structs")
}

// LoadDBConfig 读取数据库相关配置文件
func LoadDBConfig(file *ini.File) {
	var err error
	Db = file.Section("database").Key("Db").String()
	DbHost = file.Section("database").Key("DbHost").String()
	DbUser = file.Section("database").Key("DbUser").String()
	DbPassword = file.Section("database").Key("DbPassword").String()
	DbName = file.Section("database").Key("DbName").String()
	DbPort, err = file.Section("database").Key("DbPort").Int()
	if err != nil {
		log.Fatalln("database port error")
	}
}

// LoadServer 读取服务器相关配置文件
func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8000")
	JwtKey = file.Section("server").Key("JwyKey").MustString("fdasasferqw")
}

// LoadLogger 读取日志相关配置文件
func LoadLogger(file *ini.File) {
	LoggerMode = file.Section("logger").Key("LoggerMode").MustString("development")
	InfoLog = file.Section("logger").Key("InfoLog").MustString("./info.log")
	ErrorLog = file.Section("logger").Key("ErrorLog").MustString("./error.log")
	WarnLog = file.Section("logger").Key("WarnLog").MustString("./warn.log")
}

// LoadServiceLogger 读取业务相关的配置文件
func LoadServiceLogger(file *ini.File) {
	RushOrderCost = float32(file.Section("service").Key("RushOrderCost").MustFloat64(0.5))
	RushOrder2PendingTime = file.Section("service").Key("RushOrder2PendingTime").MustInt64(10)
	PendingOrder2ExpireTime = file.Section("service").Key("PendingOrder2ExpireTime").MustInt64(24 * 60)
	CoinBase = file.Section("service").Key("CoinBase").MustInt64(100)
}
