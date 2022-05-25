package utils

import "log"
import "gopkg.in/ini.v1"

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
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatalln("read config file error", err)
	}
	LoadDBConfig(file)
	LoadServer(file)
	LoadLogger(file)
}

// 读取数据库相关配置文件
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

// 读取服务器相关配置文件
func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":8000")
	JwtKey = file.Section("server").Key("JwyKey").MustString("fdasasferqw")
}

// 读取日志相关配置文件
func LoadLogger(file *ini.File) {
	LoggerMode = file.Section("logger").Key("LoggerMode").MustString("development")
	InfoLog = file.Section("loggeer").Key("InfoLog").MustString("./info.log")
	ErrorLog = file.Section("loggeer").Key("ErrorLog").MustString("./error.log")
	WarnLog = file.Section("loggeer").Key("WarnLog").MustString("./warn.log")
}
