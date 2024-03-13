package router

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/api"
)

func Engine() *gin.Engine {
	router := gin.Default()

	// 允许跨域
	//router.Use(cors.Default())
	router.NoRoute(NoRoute())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 性能监控
	//pprof.Register(router)

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("template/*")

	router.Use(InitSession())
	router.GET("/login", api.LoginApi)
	router.POST("/doLogin", api.DoLoginApi)

	auth := router.Group("")
	auth.Use(AuthSession())
	auth.GET("/logout", api.LogoutApi)
	auth.GET("/", api.VideoIndex)
	auth.GET("/list", api.VideoList)
	auth.GET("/actress", api.VideoActress)
	auth.GET("/play", api.VideoPlay)
	auth.GET("/rename", api.VideoRename)

	return router
}

func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
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

		email, ok1 := session.Get("email").(string)
		password, ok2 := session.Get("password").(string)
		if !ok1 || !ok2 {
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		fmt.Println("ssid: ", session.Get("email").(string), session.Get("password").(string))

		session.Set("email", email)
		session.Set("password", password)
		session.Save()

		c.Next()
	}
}
