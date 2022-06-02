package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"net/http"
	"service/service"
	"service/utils"
	"service/utils/errmsg"
	"service/utils/logger"
	"strings"
	"time"
)

var (
	jwtKey = []byte(utils.JwtKey)
)

// MyClaims 自定义一个cliams
type MyClaims struct {
	Id   int64
	Role string
	jwt.StandardClaims
}

// NewToken 生成token
func NewToken(id int64, role string) (string, int) {
	// 有效期
	expireTime := time.Now().Add(10 * time.Hour)
	// 声明一个Claims
	SetClaims := MyClaims{
		id,
		role,
		jwt.StandardClaims{ExpiresAt: expireTime.Unix(), Issuer: "service"},
	}
	// 新建一个声明
	reqClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	// 生成Token
	token, err := reqClaims.SignedString(jwtKey)
	if err != nil {
		logger.Log.Error("Token生成错误", zap.Error(err))
		return "", errmsg.ERROR
	}
	return token, errmsg.SUCCESS
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, int) {
	//下面这个函数是官方文档中提供的函数，用来校验token
	setToken, err := jwt.ParseWithClaims(
		token,
		&MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return nil, errmsg.ErrorTokenTimeOut
		} else {
			logger.Log.Error("Jwt校验错误", zap.Error(err))
		}
	}
	// 检验
	if key, _ := setToken.Claims.(*MyClaims); setToken.Valid {
		return key, errmsg.SUCCESS
	} else {
		return nil, errmsg.ErrorTokenWokenWrong
	}
}
func RoleValidate(targetRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("role") != targetRole {
			code := errmsg.ErrorTokenRoleNotMatch
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
	}
}

// JwtToken jwt中间件
// 定义一个gin的中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 固定写法
		tokenHandler := c.Request.Header.Get("Authorization")
		var code int
		// not exist
		if tokenHandler == "" {
			code = errmsg.ErrorTokenNotExist
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		// check the token format
		checkToken := strings.SplitN(tokenHandler, " ", 2)
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			code = errmsg.ErrorTokenTypeWrong
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		// check the validity of the token
		key, valid := CheckToken(checkToken[1])
		if valid != errmsg.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code": valid,
				"msg":  errmsg.GetErrMsg(valid),
			})
			c.Abort()
			return
		}
		// if the token timeout
		if time.Now().Unix() > key.ExpiresAt {
			code = errmsg.ErrorTokenTimeOut
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		code = service.CheckIdExist(key.Id, key.Role)
		if code != errmsg.SUCCESS {
			c.JSON(http.StatusOK, gin.H{
				"code": code,
				"msg":  errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		// 校验ID是否存在
		c.Set("id", key.Id)
		c.Set("role", key.Role)
		c.Next()
	}
}
