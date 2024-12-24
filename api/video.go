package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/model/request"
	"github.com/wxw9868/video/service"
	"github.com/wxw9868/video/utils"
)

// GetVideoListApi godoc
//
//	@Summary		视频列表
//	@Description	视频列表
//	@Tags			video
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Produce		json
//	@Param			page		query		int		false	"页码"	default(1)
//	@Param			size		query		int		false	"条数"	default(10)
//	@Param			column		query		string	false	"排序字段"
//	@Param			order		query		string	false	"排序方式"	Enums(desc, asc)
//	@Param			actress_id	query		int		false	"演员ID"
//	@Success		200			{object}	Success
//	@Failure		400			{object}	Fail
//	@Failure		404			{object}	NotFound
//	@Failure		500			{object}	ServerError
//	@Router			/video/getVideoList [get]
func GetVideoListApi(c *gin.Context) {
	var bind request.SearchVideo
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := videoService.List(bind)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("视频列表", data))
}

// GetVideoInfoApi godoc
//
//	@Summary	视频信息
//	@Tags		video
//	@Accept		application/x-www-form-urlencoded
//	@Produce	json
//	@Param		id	path		int	true	"Video ID"
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/video/getVideoInfo/{id} [get]
func GetVideoInfoApi(c *gin.Context) {
	var bind request.Common
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := videoService.Info(bind.ID, GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("视频信息", data))
}

type CollectVideo struct {
	VideoID uint `form:"video_id" json:"video_id" binding:"required"`
	Num     int  `form:"num" json:"num" binding:"required,oneof=1 -1"`
}

// CollectVideoApi godoc
//
//	@Summary		视频收藏
//	@Description	视频收藏
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			data	body		CollectVideo	true	"视频收藏"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/video/collectVideo [post]
func CollectVideoApi(c *gin.Context) {
	var bind CollectVideo
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := videoService.Collect(bind.VideoID, GetUserID(c), bind.Num); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "收藏成功"
	if bind.Num == -1 {
		msg = "取消收藏"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

// RecordPageViewsApi godoc
//
//	@Summary	视频浏览记录
//	@Tags		video
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"视频ID"
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/video/recordPageViews/{id} [get]
func RecordPageViewsApi(c *gin.Context) {
	var bind request.Common
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := videoService.Browse(bind.ID, GetUserID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "浏览记录成功"
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

type PostVideoComment struct {
	VideoComment
	Content string `form:"content" json:"content" binding:"required"`
}

// PostVideoCommentApi godoc
//
//	@Summary		视频评论
//	@Description	视频评论
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Param			data	body		VideoComment	true	"视频评论"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/comment/postVideoComment [post]
func PostVideoCommentApi(c *gin.Context) {
	var bind PostVideoComment
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := videoService.Comment(GetUser(c), bind.VideoID, bind.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("评论成功", data))
}

type ReplyVideoComment struct {
	PostVideoComment
	ParentID uint `form:"parent_id" json:"parent_id" binding:"required"`
}

// ReplyVideoCommentApi godoc
//
//	@Summary		视频评论回复
//	@Description	视频评论回复
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ReplyVideoComment	true	"视频评论回复"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/comment/replyVideoComment [post]
func ReplyVideoCommentApi(c *gin.Context) {
	var bind ReplyVideoComment
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := videoService.Reply(GetUser(c), bind.VideoID, bind.ParentID, bind.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("回复成功", data))
}

type VideoComment struct {
	VideoID uint `uri:"id" form:"video_id" json:"video_id" binding:"required"`
}

// GetVideoCommentListApi godoc
//
//	@Summary	获取视频评论列表
//	@Tags		comment
//	@Accept		json
//	@Produce	json
//	@Param		video_id	path		int	true	"Video ID"
//	@Success	200			{object}	Success
//	@Failure	400			{object}	Fail
//	@Failure	404			{object}	NotFound
//	@Failure	500			{object}	ServerError
//	@Router		/comment/getVideoCommentList/{id} [get]
func GetVideoCommentListApi(c *gin.Context) {
	var bind VideoComment
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	list, err := videoService.CommentList(GetUserID(c), bind.VideoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("评论列表", list))
}

type LikeVideoComment struct {
	CommentID uint `form:"comment_id" json:"comment_id" binding:"required"`
	Like      int8 `form:"like" json:"like" binding:"required,oneof=1 -1"`
}

// LikeVideoCommentApi godoc
//
//	@Summary	赞视频评论
//	@Tags		comment
//	@Accept		json
//	@Produce	json
//	@Param		data	body		LikeVideoComment	true	"赞视频评论"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/comment/likeVideoComment [post]
func LikeVideoCommentApi(c *gin.Context) {
	var bind LikeVideoComment
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := videoService.LikeVideoComment(GetUserID(c), bind.CommentID, bind.Like); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "点赞成功"
	if bind.Like == -1 {
		msg = "取消点赞"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

type DislikeVideoComment struct {
	CommentID uint `form:"comment_id" json:"comment_id" binding:"required"`
	Dislike   int8 `form:"dislike" json:"dislike" binding:"required,oneof=1 -1"`
}

// DislikeVideoCommentApi godoc
//
//	@Summary	踩视频评论
//	@Tags		comment
//	@Accept		json
//	@Produce	json
//	@Param		data	body		DislikeVideoComment	true	"踩视频评论"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/comment/dislikeVideoComment [post]
func DislikeVideoCommentApi(c *gin.Context) {
	var bind DislikeVideoComment
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := videoService.DislikeVideoComment(GetUserID(c), bind.CommentID, bind.Dislike); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "点踩成功"
	if bind.Dislike == -1 {
		msg = "取消点踩"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

// PostVideoDanmuApi godoc
//
//	@Summary		视频弹幕
//	@Description	视频弹幕
//	@Tags			danmu
//	@Accept			json
//	@Produce		json
//	@Param			data	body		request.CreateDanmu	true	"视频弹幕"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/danmu/postVideoDanmu [post]
func PostVideoDanmuApi(c *gin.Context) {
	var bind request.CreateDanmu
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := videoService.Danmu(GetUserID(c), bind); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("发送成功", nil))
}

// GetVideoDanmuListApi godoc
//
//	@Summary	获取视频弹幕列表
//	@Tags		danmu
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"视频ID"
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/danmu/getVideoDanmuList/{id} [get]
func GetVideoDanmuListApi(c *gin.Context) {
	var bind request.Common
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := videoService.DanmuList(bind.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("弹幕列表", data))
}

type ImportVideoData struct {
	Dir       string   `json:"dir" binding:"required"`
	Actresses []string `json:"actresses" binding:"required"`
}

// ImportVideoDataApi godoc
//
//	@Summary	导入视频数据
//	@Tags		video
//	@Accept		json
//	@Produce	json
//	@Param		data	body		ImportVideoData	true	"视频数据"
//	@Success	200		{object}	Success
//	@Failure	400		{object}	Fail
//	@Failure	404		{object}	NotFound
//	@Failure	500		{object}	ServerError
//	@Router		/video/importVideo [post]
func ImportVideoDataApi(c *gin.Context) {
	var bind ImportVideoData
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	if err := videoService.ImportVideoData(bind.Dir, bind.Actresses...); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

// VideoWriteGoFound godoc
//
//	@Summary	视频导入GoFound
//	@Tags		video
//	@Accept		json
//	@Produce	json
//	@Success	200	{object}	Success
//	@Failure	400	{object}	Fail
//	@Failure	404	{object}	NotFound
//	@Failure	500	{object}	ServerError
//	@Router		/video/writeGoFound [get]
func VideoWriteGoFound(c *gin.Context) {
	if err := service.VideoWriteGoFound(""); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

type Search struct {
	Query     string      `json:"query" binding:"required" example:""`
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Order     string      `json:"order" example:""`
	Highlight interface{} `json:"highlight"`
	ScoreExp  string      `json:"scoreExp" example:""`
}

// VideoSearchApi godoc
//
//	@Summary		视频搜索
//	@Description	get string by ID
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	true	"关键词"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/video/search [get]
func SearchVideoApi(c *gin.Context) {
	query := c.Query("query")

	b, err := json.Marshal(&Search{Query: query, Page: 1, Limit: 1000, Order: "desc"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	resp, err := client.POST(utils.Join("/query", "?", "database=video"), "application/json", bytes.NewReader(b))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	defer resp.Body.Close()

	robots, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	_, err = c.Writer.Write(robots)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
}

func ResetTableApi(c *gin.Context) {
	var table = c.Query("table")
	if err := service.ResetTable(table); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}
