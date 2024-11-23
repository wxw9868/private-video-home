package middleware

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
)

func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, util.Msg(false, http.StatusNotFound, "status not found", nil))
		c.Abort()
	}
}

func InitSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte("secret_20240109"))
	return sessions.Sessions("stock_session", store)
}

func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		_, ok := session.Get("user_id").(uint)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.Fail("没有访问权限"))
			return
		}
		c.Next()
	}
}

// GinCors 跨域
func GinCors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = false
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://127.0.0.1", "http://127.0.0.1:8080", "http://192.168.0.9", "http://192.168.0.9:80"}

	return cors.New(config)
}

// Cors 跨域
//func Cors() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		//这是允许访问所有域
//		c.Header("Access-Control-Allow-Origin", "*")
//		//跨域请求是否需要带cookie信息，默认设置为true
//		c.Header("Access-Control-Allow-Credentials", "true")
//		//header的类型
//		c.Header("Access-Control-Allow-Headers", "Action, Module, X-PINGOTHER, Content-Type, Content-Disposition,AccessToken,X-CSRF-Token, Authorization, Token,Content-Length, session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control,Pragma")
//		//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
//		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
//		//跨域关键设置，让浏览器可以解析
//		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Cache-Control, Content-Language, Expires, Last-Modified, Pragma, FooBar")
//		//缓存请求信息，单位为秒
//		c.Header("Access-Control-Max-Age", "172800")
//		//设置返回格式是json
//		c.Set("content-type", "application/json")
//
//		//放行所有OPTIONS方法
//		if c.Request.Method == "OPTIONS" {
//			c.AbortWithStatus(http.StatusNoContent)
//		}
//		c.Next()
//	}
//}
