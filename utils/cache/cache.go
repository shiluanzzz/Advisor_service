package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"service-backend/utils/errmsg"
	"service-backend/utils/logger"
	"service-backend/utils/setting"
	"time"
)

var redisConn *redis.Pool

func InitRedis() {
	redisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				log.Fatalf("redis启动失败")
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
					defer func() {
						err := c.Close()
						if err != nil {
						}
					}()
					log.Fatalf("redis启动失败")
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
func closeRedisConn(conn *redis.Conn) {
	if err := (*conn).Close(); err != nil {
		logger.Log.Warn("关闭redis缓存连接错误", zap.Error(err))
	}
}
func set(key string, data interface{}, time int) (code int) {
	conn := redisConn.Get()
	defer closeRedisConn(&conn)
	defer func() {
		fields := []zapcore.Field{
			zap.String("key", key),
			zap.Int("expireTime", time),
			zap.String("data", fmt.Sprintf("%v", data)),
		}
		if code == errmsg.SUCCESS {
			logger.Log.Info("缓存写入成功", fields...)
		} else {
			logger.Log.Warn("缓存写入失败", fields...)
		}
	}()
	var value []byte
	var err error
	if value, err = json.Marshal(data); err != nil {
		return errmsg.ErrorCacheMarshal
	}

	if _, err = conn.Do("SET", key, value); err != nil {
		return errmsg.ErrorCacheDoSet
	}

	if _, err = conn.Do("EXPIRE", key, time); err != nil {
		return errmsg.ErrorCacheDoExpire
	}
	return errmsg.SUCCESS
}
func exists(key string) (exists bool) {
	conn := redisConn.Get()
	defer closeRedisConn(&conn)
	var err error
	if exists, err = redis.Bool(conn.Do("EXISTS", key)); err != nil {
		return false
	}
	return exists
}
func get(key string) (code int, res []byte) {
	conn := redisConn.Get()
	var err error
	defer closeRedisConn(&conn)
	defer func() {
		fields := []zapcore.Field{
			zap.String("key", key),
			zap.String("data", fmt.Sprintf("%v", res)),
			zap.Error(err),
		}
		if code == errmsg.SUCCESS {
			logger.Log.Info("缓存读取成功", fields...)
		} else {
			logger.Log.Warn("缓存读取失败", fields...)
		}
	}()

	if res, err = redis.Bytes(conn.Do("GET", key)); err != nil {

		return errmsg.ErrorCacheGetBytes, nil
	}
	return errmsg.SUCCESS, res
}
func deleteData(key string) (code int, do bool) {
	conn := redisConn.Get()
	defer closeRedisConn(&conn)
	var err error
	if do, err = redis.Bool(conn.Do("DEL", key)); err != nil {
		return errmsg.ErrorCacheDeleteKey, false
	}
	return errmsg.SUCCESS, do
}

// redis key值相关的设定

func GetOrderKey(id int64) string {
	return fmt.Sprintf("orderdetail_%d", id)
}
func GetCommentKey(id int64) string {
	return fmt.Sprintf("advisor's_ordercomment_%d", id)
}

// 封装好缓存的接口三个接口

func GetCacheData(key string, response interface{}) (code int) {
	if exists(key) {
		var cacheData []byte
		if code, cacheData = get(key); code != errmsg.SUCCESS {
			return code
		}
		if err := json.Unmarshal(cacheData, &response); err != nil {
			return errmsg.ErrorCacheUnmarshal
		}
		return errmsg.SUCCESS
	}
	return errmsg.ErrorCacheKeyNotExist
}
func SetCacheData(key string, data interface{}) {
	set(key, data, 3600)
}
func DeleteCacheData(key string) {
	deleteData(key)
}
