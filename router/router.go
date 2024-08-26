package router

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/wxw9868/video/api"
	"github.com/wxw9868/video/middleware"
)

func Engine(addr string) *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)
	// 强制日志颜色化
	gin.ForceConsoleColor()
	gin.DefaultWriter = colorable.NewColorableStdout()

	router := gin.Default()

	// 允许跨域
	router.Use(middleware.GinCors())

	// 性能监控
	pprof.Register(router)

	router.NoRoute(middleware.NoRoute())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/assets", "./assets")

	router.Use(middleware.InitSession())
	user := router.Group("/user")
	user.POST("/doRegister", api.RegisterApi)
	user.POST("/doLogin", api.LoginApi)
	user.POST("/sendMail", api.SendMailApi)
	user.POST("/sendUrl", api.SendUrlApi)
	user.POST("/captcha", api.CaptchaApi)
	user.POST("/forgotPassword", api.ForgotPasswordApi)

	auth := router.Group("", middleware.AuthSession())
	user = auth.Group("/user")
	user.GET("/logout", api.LogoutApi)
	user.GET("/session", api.GetSession)
	user.GET("/info", api.UserInfoApi)
	user.POST("/update", api.UserUpdateApi)
	user.POST("/uploadAvatar", api.UserUploadAvatarApi)
	user.POST("/changePassword", api.ChangePasswordApi)
	user.GET("/collect", api.UserCollectApi)
	user.GET("/browse", api.UserBrowseApi)

	video := auth.Group("/video")
	video.GET("/getSearch", api.VideoSearchApi)
	video.GET("/getList", api.VideoListApi)
	video.GET("/getPlay", api.VideoPlayApi)
	video.GET("/browse", api.VideoBrowseApi)
	video.POST("/collect", api.VideoCollectApi)
	video.GET("/import", api.VideoImport)
	video.GET("/importVideoActress", api.ImportVideoActress)
	video.GET("/getVideoPic", api.GetVideoPic)
	video.GET("/oneVideoPic", api.OneVideoPic)

	comment := auth.Group("/comment")
	comment.GET("/list", api.VideoCommentListApi)
	comment.POST("/comment", api.VideoCommentApi)
	comment.POST("/reply", api.VideoReplyApi)
	comment.POST("/zan", api.CommentZanApi)
	comment.POST("/cai", api.CommentCaiApi)

	danmu := auth.Group("/danmu")
	danmu.GET("/list", api.DanmuListApi)
	danmu.POST("/save", api.DanmuSaveApi)

	search := auth.Group("/search")
	search.POST("/api/index", api.SearchIndex)
	search.POST("/api/index/batch", api.SearchIndexBatch)
	search.POST("/api/index/remove", api.SearchIndexRemove)
	search.POST("/api/query", api.SearchQuery)
	search.GET("/api/status", api.SearchStatus)
	search.GET("/api/db/drop", api.SearchDbDrop)
	search.GET("/api/word/cut", api.SearchWordCut)

	actress := auth.Group("/actress")
	actress.POST("/add", api.ActressAddApi)
	actress.POST("/edit", api.ActressEditApi)
	actress.GET("/delete", api.ActressDeleteApi)
	actress.GET("/list", api.ActressListApi)
	actress.GET("/addInfo", api.AdditionalInformation)

	tag := auth.Group("/tag")
	tag.POST("/add", api.TagAddApi)
	tag.POST("/edit", api.TagEditApi)
	tag.GET("/delete", api.TagDeleteApi)
	tag.GET("/list", api.TagListApi)

	stock := auth.Group("/stock")
	stock.POST("/import_trading_records", api.ImportTradingRecordsApi)
	stock.GET("/liquidation", api.LiquidationApi)
	stock.GET("/trading_records", api.TradingRecordsApi)

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
