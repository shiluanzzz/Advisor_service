package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"service-backend/utils/logger"
	"time"
)

func Log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		ctx.Next()
		// 记录一些信息
		stopTime := time.Since(startTime).Milliseconds()
		Cost := fmt.Sprintf("%d ms", stopTime)
		hostName, err := os.Hostname()
		statusCode := ctx.Writer.Status()
		if err != nil {
			hostName = "Unknown"
		}
		if len(ctx.Errors) > 0 {
			for _, e := range ctx.Errors.Errors() {
				logger.Log.Error(e)
			}
		} else {
			fields := []zapcore.Field{
				zap.Int("status", ctx.Writer.Status()),
				zap.String("method", ctx.Request.Method),
				zap.String("hostname", hostName),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", ctx.ClientIP()),
				zap.String("user-agent", ctx.Request.UserAgent()),
				zap.String("cost", Cost),
			}
			// 根据不同的响应输出不同等级的日志.
			if statusCode >= 500 {
				logger.Log.Error(path, fields...)
			} else if statusCode >= 400 {
				logger.Log.Warn(path, fields...)
			} else {
				logger.Log.Info(path, fields...)
			}
		}
	}
}
