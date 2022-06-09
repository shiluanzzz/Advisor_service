package setting

import (
	"log"
	"time"
)
import "gopkg.in/ini.v1"
import "github.com/didi/gendry/scanner"

type serviceConfig struct {
	RushOrderCost           float32
	RushOrder2PendingTime   int64
	PendingOrder2ExpireTime int64
	CoinBase                int64
}

var ServiceCfg = &serviceConfig{}

type logger struct {
	LoggerMode string
	InfoLog    string
	ErrorLog   string
	WarnLog    string
}

var Logger = &logger{}

type server struct {
	HttpPort string
	AppMode  string
	JwtKey   string
}

var Server = &server{}

type db struct {
	Db         string
	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	DbName     string
}

var DB = &db{}

type redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

var RedisSetting = &redis{}

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatalln("read config file error", err)
	}
	Setting()
	MapTo(file, "redis", RedisSetting)
	MapTo(file, "logger", Logger)
	MapTo(file, "database", DB)
	MapTo(file, "server", Server)
	MapTo(file, "service", ServiceCfg)
}

// Setting 一些库的相关设置
func Setting() {
	// 用于scanner的字段反射
	scanner.SetTagName("structs")
	// 因为结构体转map用到了structs这个库，这个库是根据structs这个tag来反射生成map的。
}
func MapTo(cfg *ini.File, section string, v interface{}) {
	if err := cfg.Section(section).MapTo(v); err != nil {
		log.Fatalf("读取配置文件错误,section:%v,err:%v\n", section, err)
	}
}
