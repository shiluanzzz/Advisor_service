package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"service/utils"
)

var Log *zap.Logger

func init() {
	// 三类日志的write
	infoWrite := getLogWriter(utils.InfoLog)
	errorWrite := getLogWriter(utils.ErrorLog)
	warnWrite := getLogWriter(utils.WarnLog)

	// 三类日志的level
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level <= zap.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.WarnLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level == zap.WarnLevel
	})
	infoCore := zapcore.NewCore(getLogEncoder(), infoWrite, infoLevel)
	errorCore := zapcore.NewCore(getLogEncoder(), errorWrite, errorLevel)
	warnCore := zapcore.NewCore(getLogEncoder(), warnWrite, warnLevel)
	// 组合三种核心
	core := zapcore.NewTee(infoCore, errorCore, warnCore)
	Log = zap.New(core, zap.AddCaller())
	// logger会有缓存因此退出的时候需要同步
	defer Log.Sync()
}
func getLogEncoder() zapcore.Encoder {
	// TODO 两种config有什么区别?
	encoderConfig := zap.NewProductionEncoderConfig()
	//更改打印的时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// TODO ?
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	switch utils.LoggerMode {
	case "development":
		// 开发环境下直接打印查看
		return zapcore.NewConsoleEncoder(encoderConfig)
	default:
		return zapcore.NewJSONEncoder(encoderConfig)
	}
}
func getLogWriter(filepath string) zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    1,  //大小限制 MB
		MaxAge:     30, // 保留天数
		MaxBackups: 5,  // 保留数量
		Compress:   false,
	})
}

// common log
func GendryError(err error) {
	Log.Error("Gendry错误", zap.Error(err))
}
func SqlInsertError(err error) {
	Log.Error("Mysql Insert error", zap.Error(err))
}
func SqlUpdateError(err error) {
	Log.Error("Mysql Update error", zap.Error(err))
}
func SqlSelectError(err error) {
	Log.Error("Mysql Select error", zap.Error(err))
}
func SqlDeleteError(err error) {
	Log.Error("Mysql Delete error", zap.Error(err))
}