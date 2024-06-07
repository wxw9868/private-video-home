package router

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/api"
)

func Engine(addr string) *gin.Engine {
	router := gin.Default()

	// Setup Security Headers
	// router.Use(func(c *gin.Context) {
	// 	if c.Request.Host != addr {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
	// 		return
	// 	}
	// 	c.Header("X-Frame-Options", "DENY")
	// 	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	// 	c.Header("X-XSS-Protection", "1; mode=block")
	// 	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	// 	c.Header("Referrer-Policy", "strict-origin")
	// 	c.Header("X-Content-Type-Options", "nosniff")
	// 	c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
	// 	c.Next()
	// })

	// 允许跨域
	// router.Use(cors.Default())

	//
	router.NoRoute(NoRoute())

	//
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
	router.GET("/register", api.Register)
	router.GET("/login", api.Login)
	router.POST("/doRegister", api.RegisterApi)
	router.POST("/doLogin", api.LoginApi)

	search := router.Group("/search")
	search.POST("/api/index", api.SearchIndex)
	search.POST("/api/index/batch", api.SearchIndexBatch)
	search.POST("/api/index/remove", api.SearchIndexRemove)
	search.POST("/api/query", api.SearchQuery)
	search.GET("/api/status", api.SearchStatus)
	search.GET("/api/db/drop", api.SearchDbDrop)
	search.GET("/api/word/cut", api.SearchWordCut)

	auth := router.Group("", AuthSession())
	auth.GET("/session", api.GetSession)
	auth.GET("/account", api.Account)

	auth.GET("/logout", api.LogoutApi)

	auth.GET("/", api.VideoIndex)

	auth.GET("/list", api.VideoList)
	auth.GET("/actress", api.VideoActress)
	auth.GET("/play", api.VideoPlay)

	auth.GET("/browse", api.VideoBrowseApi)
	auth.POST("/collect", api.VideoCollectApi)
	auth.GET("/commentList", api.VideoCommentListApi)
	auth.POST("/comment", api.VideoCommentApi)
	auth.POST("/reply", api.VideoReplyApi)
	auth.POST("/zan", api.CommentZanApi)
	auth.POST("/cai", api.CommentCaiApi)

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
			fmt.Println("user out: ", userid, ok)
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
