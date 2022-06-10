package logger

import (
	"fmt"
	"service-backend/model"
	"service-backend/utils/errmsg"
	"service-backend/utils/setting"
	"service-backend/utils/tools"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func init() {
	// 三类日志的write
	infoWrite := getLogWriter(setting.Logger.InfoLog)
	errorWrite := getLogWriter(setting.Logger.ErrorLog)
	warnWrite := getLogWriter(setting.Logger.WarnLog)

	// 三类日志的level
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level <= zap.InfoLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level > zap.WarnLevel
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
	defer func() {
		_ = Log.Sync()
	}()
}
func getLogEncoder() zapcore.Encoder {
	// TODO 两种config有什么区别?
	var encoderConfig zapcore.EncoderConfig
	if setting.Logger.LoggerMode == "development" {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	//更改打印的时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// TODO ?
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	switch setting.Logger.LoggerMode {
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

func GendryBuildError(err error, args ...interface{}) {
	fields := []zapcore.Field{
		zap.String("function", tools.WhoCallMe()),
		zap.Error(err),
	}
	for i := 0; i < len(args)-1; i += 2 {
		fields = append(fields, zap.String(fmt.Sprintf("%v", args[i]), fmt.Sprintf("%v", args[i+1])))
	}
	Log.Error("Gendry build SQL错误", fields...)
}
func GendryScannerError(err error, args ...interface{}) {
	fields := []zapcore.Field{
		zap.String("function", tools.WhoCallMe()),
		zap.Error(err),
	}
	for i := 0; i < len(args)-1; i += 2 {
		fields = append(fields, zap.String(fmt.Sprintf("%v", args[i]), fmt.Sprintf("%v", args[i+1])))
	}
	Log.Error("Gendry scanner 绑定数据错误", fields...)
}
func SqlError(err error, args ...interface{}) {
	fields := []zapcore.Field{
		zap.String("function", tools.WhoCallMe()),
		zap.Error(err),
	}
	for i := 0; i < len(args)-1; i += 2 {
		fields = append(fields, zap.String(fmt.Sprintf("%v", args[i]), fmt.Sprintf("%v", args[i+1])))
	}
	Log.Error(fmt.Sprintf("mysql error"), fields...)
}

// CommonControllerLog 控制层的defer Log
func CommonControllerLog(code *int, msg *string, requests interface{}, response interface{}) {
	commonLog(model.ControllerLog, code, "function", tools.WhoCallMe(), "msg", msg, "requests", requests, "response", response)
}

// CommonServiceLog 服务层的defer Log
func CommonServiceLog(code *int, input interface{}, args ...interface{}) {
	args = append(args, []interface{}{"function", tools.WhoCallMe(), "input", input}...)
	commonLog(model.ServiceLog, code, args...)
}

func commonLog(kind model.LogType, code *int, args ...interface{}) {

	fields := []zapcore.Field{
		zap.String("LogType", kind.StatusName()),
	}
	for i := 0; i < len(args)-1; i += 2 {
		fields = append(fields, zap.String(fmt.Sprintf("%v", args[i]), fmt.Sprintf("%v", args[i+1])))
	}
	switch *code {
	case errmsg.SUCCESS:
		Log.Info("success", fields...)
	case errmsg.ERROR:
		Log.Error("error", fields...)
	case errmsg.ErrorMysql:
		Log.Error("error", fields...)
	default:
		Log.Warn("warn", fields...)
	}
}
