package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/video/api"
	"github.com/wxw9868/video/middleware"
)

var relativePath string

func Router(r *gin.Engine) {
	r.Use(middleware.InitSession())
	public := r.Group(relativePath)
	private := r.Group(relativePath, middleware.Authentication())

	public.POST("user/register", api.RegisterApi)
	public.POST("user/login", api.LoginApi)
	public.POST("user/forgotPassword", api.ForgotPasswordApi)
	public.POST("verify/sendMail", api.SendMailApi)
	public.POST("verify/sendUrl", api.SendUrlApi)
	public.POST("verify/captcha", api.CaptchaApi)

	user := private.Group("user")
	{
		user.GET("logout", api.LogoutApi)
		user.GET("getSession", api.GetSessionApi)
		user.GET("getUserInfo", api.GetUserInfoApi)
		user.POST("updateUserInfo", api.UpdateUserInfoApi)
		user.POST("changeUserAvatar", api.ChangeUserAvatarApi)
		user.POST("changePassword", api.ChangePasswordApi)
		user.GET("getUserFavoriteList", api.GetUserFavoriteListApi)
		user.GET("getUserBrowseList", api.GetUserBrowseListApi)
		user.GET("getUserLoginLogListApi", api.GetUserLoginLogListApi)
	}

	actress := private.Group("actress")
	{
		actress.POST("createActress", api.CreateActressApi)
		actress.GET("getActressList", api.GetActressListApi)
		actress.GET("getActressInfo/:id", api.GetActressInfoApi)
		actress.POST("updateActress", api.UpdateActressApi)
		actress.DELETE("deleteActress/:id", api.DeleteActressApi)
	}

	video := private.Group("video")
	{
		video.GET("getVideoList", api.GetVideoListApi)
		video.GET("getVideoInfo/:id", api.GetVideoInfoApi)
		video.GET("recordPageViews/:id", api.RecordPageViewsApi)
		video.POST("collectVideo", api.CollectVideoApi)
		video.POST("importVideoData", api.ImportVideoDataApi)
		video.GET("writeGoFound", api.VideoWriteGoFound)
		video.GET("searchVideo", api.VideoSearchApi)

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

	search := private.Group("search")
	{
		search.POST("api/index", api.SearchIndex)
		search.POST("api/index/batch", api.SearchIndexBatch)
		search.POST("api/index/remove", api.SearchIndexRemove)
		search.POST("api/query", api.SearchQuery)
		search.GET("api/status", api.SearchStatus)
		search.GET("api/db/drop", api.SearchDbDrop)
		search.GET("api/word/cut", api.SearchWordCut)
	}

	tag := private.Group("tag")
	tag.POST("createTag", api.CreateTagApi)
	tag.POST("updateTag", api.UpdateTagApi)
	tag.POST("getTagList", api.GetTagListApi)
	tag.GET("deleteTag/:id", api.DeleteTagApi)

	stock := private.Group("stock")
	stock.POST("importTradingRecords", api.ImportTradingRecordsApi)
	stock.GET("liquidation", api.LiquidationApi)
	stock.GET("tradingRecords", api.TradingRecordsApi)

	system := private.Group("system")
	system.GET("resetTable", api.ResetTableApi)

	util := private.Group("util")
	util.GET("oneAddInfo", api.OneAddInfoToActress)
	util.GET("allAddInfo", api.AllAddInfoToActress)
	util.GET("oneVideoPic", api.OneVideoPic)
	util.GET("getVideoPic", api.GetVideoPic)
}
