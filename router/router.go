package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/api"
	"github.com/wxw9868/video/middleware"
)

func Router(r *gin.Engine) {
	r.Use(middleware.InitSession())
	user := r.Group("user")
	user.POST("register", api.RegisterApi)
	user.POST("login", api.LoginApi)
	user.POST("forgotPassword", api.ForgotPasswordApi)
	user.POST("sendMail", api.SendMailApi)
	user.POST("sendUrl", api.SendUrlApi)
	user.POST("captcha", api.CaptchaApi)

	auth := r.Group("", middleware.AuthSession())
	user = auth.Group("user")
	user.GET("logout", api.LogoutApi)
	user.GET("session", api.SessionApi)
	user.GET("getUserInfo", api.GetUserInfoApi)
	user.POST("updateUserInfo", api.UpdateUserInfoApi)
	user.POST("changeUserAvatar", api.ChangeUserAvatarApi)
	user.POST("changePassword", api.ChangePasswordApi)
	user.GET("getUserFavoriteList", api.GetUserFavoriteListApi)
	user.GET("getUserBrowseList", api.GetUserBrowseListApi)

	actress := auth.Group("actress")
	actress.POST("addActress", api.AddActressApi)
	actress.POST("updateActress", api.UpdateActressApi)
	actress.GET("deleteActress/:id", api.DeleteActressApi)
	actress.POST("getActressList", api.GetActressListApi)
	actress.GET("getActressInfo/:id", api.GetActressInfoApi)

	video := auth.Group("video")
	{
		video.GET("searchVideo", api.VideoSearchApi)
		video.POST("getVideoList", api.GetVideoListApi)
		video.GET("videoPlay/:id", api.VideoPlayApi)
		video.GET("browseVideo/:id", api.BrowseVideoApi)
		video.POST("collectVideo", api.CollectVideoApi)
		video.POST("importVideoData", api.ImportVideoDataApi)
		video.GET("writeGoFound", api.VideoWriteGoFound)

		comment := video.Group("comment")
		{
			comment.GET("getVideoCommentList/:id", api.GetVideoCommentListApi)
			comment.POST("videoComment", api.VideoCommentApi)
			comment.POST("replyVideoComment", api.ReplyVideoCommentApi)
			comment.POST("likeVideoComment", api.LikeVideoCommentApi)
			comment.POST("dislikeVideoComment", api.DislikeVideoCommentApi)
		}
		danmu := video.Group("danmu")
		{
			danmu.POST("sendVideoBarrage", api.SendVideoBarrageApi)
			danmu.GET("getVideoBarrageList/:id", api.GetVideoBarrageListApi)
		}
	}

	search := auth.Group("search")
	search.POST("api/index", api.SearchIndex)
	search.POST("api/index/batch", api.SearchIndexBatch)
	search.POST("api/index/remove", api.SearchIndexRemove)
	search.POST("api/query", api.SearchQuery)
	search.GET("api/status", api.SearchStatus)
	search.GET("api/db/drop", api.SearchDbDrop)
	search.GET("api/word/cut", api.SearchWordCut)

	tag := auth.Group("tag")
	tag.POST("createTag", api.CreateTagApi)
	tag.POST("updateTag", api.UpdateTagApi)
	tag.POST("getTagList", api.GetTagListApi)
	tag.GET("deleteTag/:id", api.DeleteTagApi)

	stock := auth.Group("stock")
	stock.POST("importTradingRecords", api.ImportTradingRecordsApi)
	stock.GET("liquidation", api.LiquidationApi)
	stock.GET("tradingRecords", api.TradingRecordsApi)

	system := auth.Group("system")
	system.GET("resetTable", api.ResetTableApi)

	util := auth.Group("util")
	util.GET("oneAddInfo", api.OneAddInfoToActress)
	util.GET("allAddInfo", api.AllAddInfoToActress)
	util.GET("oneVideoPic", api.OneVideoPic)
	util.GET("getVideoPic", api.GetVideoPic)
}
