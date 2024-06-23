package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/api"
	"github.com/wxw9868/video/middleware"
)

func Engine(addr string) *gin.Engine {
	router := gin.Default()

	// 允许跨域
	// router.Use(cors.Default())
	router.Use(middleware.Cors())

	// 性能监控
	// pprof.Register(router)

	router.NoRoute(middleware.NoRoute())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("template/*")

	router.Use(middleware.InitSession())
	router.GET("/", api.VideoIndex)
	user := router.Group("/user")
	user.GET("/register", api.Register)
	user.GET("/login", api.Login)
	user.POST("/doRegister", api.RegisterApi)
	user.POST("/doLogin", api.LoginApi)

	auth := router.Group("", middleware.AuthSession())
	user = auth.Group("/user")
	user.GET("/account", api.Account)
	user.GET("/session", api.GetSession)
	user.GET("/logout", api.LogoutApi)

	video := auth.Group("/video")
	video.GET("/list", api.VideoList)
	video.GET("/search", api.VideoSearch)
	video.GET("/actress", api.VideoActress)
	video.GET("/play", api.VideoPlay)
	video.GET("/getPlay", api.VideoPlayApi)
	video.GET("/getList", api.VideoListApi)
	video.GET("/getActress", api.VideoActressApi)
	video.GET("/browse", api.VideoBrowseApi)
	video.POST("/collect", api.VideoCollectApi)
	video.GET("/rename", api.VideoRename)
	video.GET("/import", api.VideoImport)
	video.GET("/importVideoActress", api.ImportVideoActress)

	comment := auth.Group("/comment")
	comment.GET("/list", api.VideoCommentListApi)
	comment.POST("/comment", api.VideoCommentApi)
	comment.POST("/reply", api.VideoReplyApi)
	comment.POST("/zan", api.CommentZanApi)
	comment.POST("/cai", api.CommentCaiApi)

	search := router.Group("/search")
	search.POST("/api/index", api.SearchIndex)
	search.POST("/api/index/batch", api.SearchIndexBatch)
	search.POST("/api/index/remove", api.SearchIndexRemove)
	search.POST("/api/query", api.SearchQuery)
	search.GET("/api/status", api.SearchStatus)
	search.GET("/api/db/drop", api.SearchDbDrop)
	search.GET("/api/word/cut", api.SearchWordCut)

	// 简单的路由组: v1
	v1 := router.Group("/v1")
	{
		auth := v1.Group("")
		auth.Use(middleware.AuthSession())
		auth.GET("/cache/list", api.CacheVideoList)
		auth.GET("/cache/actress", api.CacheVideoActress)
		auth.GET("/cache/play", api.CacheVideoPlay)
		auth.GET("/cache/rename", api.VideoRename)
	}

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

	return router
}
