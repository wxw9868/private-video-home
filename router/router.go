package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/api"
	"github.com/wxw9868/video/middleware"
)

func Router(r *gin.Engine) {
	r.Use(middleware.InitSession())
	user := r.Group("/user")
	user.POST("/register", api.RegisterApi)
	user.POST("/login", api.LoginApi)
	user.POST("/forgotPassword", api.ForgotPasswordApi)
	user.POST("/sendMail", api.SendMailApi)
	user.POST("/sendUrl", api.SendUrlApi)
	user.POST("/captcha", api.CaptchaApi)

	auth := r.Group("", middleware.AuthSession())
	user = auth.Group("/user")
	user.GET("/logout", api.LogoutApi)
	user.GET("/session", api.SessionApi)
	user.GET("/info", api.UserInfoApi)
	user.POST("/save", api.UserSaveApi)
	user.POST("/avatar", api.UserUploadAvatarApi)
	user.POST("/changePassword", api.ChangePasswordApi)
	user.GET("/collect", api.UserCollectApi)
	user.GET("/browse", api.UserBrowseApi)

	actress := auth.Group("/actress")
	actress.POST("/add", api.ActressAddApi)
	actress.POST("/edit", api.ActressEditApi)
	actress.GET("/delete/:id", api.ActressDeleteApi)
	actress.POST("/list", api.ActressListApi)
	actress.GET("/info/:id", api.ActressInfoApi)
	actress.GET("/oneAddInfo", api.OneAddInfoToActress)
	actress.GET("/allAddInfo", api.AllAddInfoToActress)

	video := auth.Group("/video")
	video.GET("/search", api.VideoSearchApi)
	video.POST("/list", api.VideoListApi)
	video.GET("/play/:id", api.VideoPlayApi)
	video.GET("/browse/:video_id", api.VideoBrowseApi)
	video.POST("/collect", api.VideoCollectApi)
	video.GET("/import", api.VideoImportApi)
	video.GET("/repairImport", api.RepairVideoImportApi)
	video.GET("/writeGoFound", api.VideoWriteGoFound)
	video.GET("/getVideoPic", api.GetVideoPic)
	video.GET("/oneVideoPic", api.OneVideoPic)

	comment := auth.Group("/comment")
	comment.GET("/list/:video_id", api.VideoCommentListApi)
	comment.POST("/comment", api.VideoCommentApi)
	comment.POST("/reply", api.VideoReplyApi)
	comment.POST("/zan", api.CommentZanApi)
	comment.POST("/cai", api.CommentCaiApi)

	danmu := auth.Group("/danmu")
	danmu.GET("/list/:video_id", api.DanmuListApi)
	danmu.POST("/save", api.DanmuSaveApi)

	search := auth.Group("/search")
	search.POST("/api/index", api.SearchIndex)
	search.POST("/api/index/batch", api.SearchIndexBatch)
	search.POST("/api/index/remove", api.SearchIndexRemove)
	search.POST("/api/query", api.SearchQuery)
	search.GET("/api/status", api.SearchStatus)
	search.GET("/api/db/drop", api.SearchDbDrop)
	search.GET("/api/word/cut", api.SearchWordCut)

	tag := auth.Group("/tag")
	tag.POST("/add", api.TagAddApi)
	tag.POST("/edit", api.TagEditApi)
	tag.POST("/list", api.TagListApi)
	tag.GET("/delete/:id", api.TagDeleteApi)

	stock := auth.Group("/stock")
	stock.POST("/importTradingRecords", api.ImportTradingRecordsApi)
	stock.GET("/liquidation", api.LiquidationApi)
	stock.GET("/tradingRecords", api.TradingRecordsApi)

	auth.GET("/resetTable", api.ResetTableApi)

	// // Setup Security Headers
	// r.Use(func(c *gin.Context) {
	// 	if c.Request.Host != addr {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
	// 		return
	// 	}
	// 	c.Header("X-Frame-Options", "DENY")
	// 	c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
	// 	c.Header("X-XSS-Protection", "1; mode=block")
	// 	c.Header("Strict-Transport-Security", "max-age=pu31536000; includeSubDomains; preload")
	// 	c.Header("Referrer-Policy", "strict-origin")
	// 	c.Header("X-Content-Type-Options", "nosniff")
	// 	c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
	// 	c.Next()
	// })
}
