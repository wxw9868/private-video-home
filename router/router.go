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
	router.GET("/preivew", func(c *gin.Context) {
		c.HTML(http.StatusOK, "preview.html", nil)
	})

	auth := router.Group("", AuthSession())
	auth.GET("/logout", api.LogoutApi)
	auth.GET("/", api.VideoIndex)
	auth.GET("/list", api.VideoList)
	auth.GET("/actress", api.VideoActress)
	auth.GET("/play", api.VideoPlay)
	auth.POST("/collect", api.VideoCollectApi)
	auth.GET("/browse", api.VideoBrowseApi)
	auth.POST("/comment", api.VideoCommentApi)
	auth.POST("/reply", api.VideoReplyApi)
	auth.GET("/commentList", api.VideoCommentListApi)

	auth.GET("/rename", api.VideoRename)
	auth.GET("/import", api.VideoImport)

	// 简单的路由组: v1
	v1 := router.Group("/v1")
	{
		auth := v1.Group("")
		auth.Use(AuthSession())
		auth.GET("/cache/list", api.CacheVideoList)
		auth.GET("/cache/actress", api.CacheVideoActress)
		auth.GET("/cache/play", api.CacheVideoPlay)
		auth.GET("/cache/rename", api.VideoRename)
	}

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

		userid, ok := session.Get("userID").(uint)
		fmt.Println("user: ", userid, ok)
		if !ok {
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
