package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/wxw9868/util"
	"github.com/wxw9868/video/service"
	"github.com/wxw9868/video/utils"
)

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
func VideoSearchApi(c *gin.Context) {
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

type Video struct {
	Paginate
	OrderBy
	ActressID int `uri:"actress_id" form:"actress_id" json:"actress_id"`
}

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
	var bind Video
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := videoService.List(bind.ActressID, bind.Page, bind.Size, bind.Column, bind.Order)
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
	var bind Common
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

type VideoCollect struct {
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
//	@Param			data	body		VideoCollect	true	"视频收藏"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/video/collectVideo [post]
func CollectVideoApi(c *gin.Context) {
	var bind VideoCollect
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
//	@Summary		视频浏览记录
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"视频ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/video/recordPageViews/{id} [get]
func RecordPageViewsApi(c *gin.Context) {
	var bind Common
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

type VideoComment struct {
	VideoID uint   `form:"video_id" json:"video_id" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

// VideoCommentApi godoc
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
//	@Router			/comment/videoComment [post]
func VideoCommentApi(c *gin.Context) {
	var bind VideoComment
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	commentID, err := videoService.Comment(bind.VideoID, bind.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	userAvatar := session.Get("user_avatar").(string)
	userNickname := session.Get("user_nickname").(string)

	data := map[string]interface{}{
		"commentID":    commentID,
		"userAvatar":   userAvatar,
		"userNickname": userNickname,
		"content":      bind.Content,
	}
	c.JSON(http.StatusOK, util.Success("评论成功", data))
}

type VideoReply struct {
	VideoID  uint   `form:"video_id" json:"video_id" binding:"required"`
	ParentID uint   `form:"parent_id" json:"parent_id" binding:"required"`
	Content  string `form:"content" json:"content" binding:"required"`
}

// ReplyVideoCommentApi godoc
//
//	@Summary		视频评论回复
//	@Description	视频评论回复
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Param			data	body		VideoReply	true	"视频评论回复"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/comment/replyVideoComment [post]
func ReplyVideoCommentApi(c *gin.Context) {
	var bind VideoReply
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	commentID, err := videoService.Reply(bind.VideoID, bind.ParentID, bind.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	userAvatar := session.Get("user_avatar").(string)
	userNickname := session.Get("user_nickname").(string)

	data := map[string]interface{}{
		"commentID":    commentID,
		"userAvatar":   userAvatar,
		"userNickname": userNickname,
		"content":      bind.Content,
	}

	c.JSON(http.StatusOK, util.Success("回复成功", data))
}

type Comment struct {
	ID uint `uri:"id" binding:"required"`
}

// GetVideoCommentListApi godoc
//
//	@Summary		视频弹幕列表
//	@Description	视频弹幕列表
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Param			video_id	path		int	true	"Video ID"
//	@Success		200			{object}	Success
//	@Failure		400			{object}	Fail
//	@Failure		404			{object}	NotFound
//	@Failure		500			{object}	ServerError
//	@Router			/comment/getVideoCommentList/{id} [get]
func GetVideoCommentListApi(c *gin.Context) {
	var bind Comment
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)
	list, err := videoService.CommentList(bind.ID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("评论列表", list))
}

type CommentZan struct {
	CommentID uint `form:"comment_id" json:"comment_id" binding:"required"`
	Like      int  `form:"like" json:"like" binding:"required,oneof=1 -1"`
}

// LikeVideoCommentApi godoc
//
//	@Summary		视频评论回复赞
//	@Description	视频评论回复赞
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Param			data	body		CommentZan	true	"视频评论回复赞"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/comment/likeVideoComment [post]
func LikeVideoCommentApi(c *gin.Context) {
	var bind CommentZan
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	if err := videoService.Zan(bind.CommentID, userID, bind.Like); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "点赞成功"
	if bind.Like == -1 {
		msg = "取消点赞"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

type CommentCai struct {
	CommentID uint `form:"comment_id" json:"comment_id" binding:"required"`
	DisLike   int  `form:"dis_like" json:"dis_like" binding:"required,oneof=1 -1"`
}

// DislikeVideoCommentApi godoc
//
//	@Summary		视频评论回复踩
//	@Description	视频评论回复踩
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Param			data	body		CommentCai	true	"视频评论回复踩"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/comment/dislikeVideoComment [post]
func DislikeVideoCommentApi(c *gin.Context) {
	var bind CommentCai
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	if err := videoService.Cai(bind.CommentID, userID, bind.DisLike); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "点踩成功"
	if bind.DisLike == -1 {
		msg = "取消踩"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

// GetVideoBarrageListApi godoc
//
//	@Summary		获取视频弹幕列表
//	@Description	获取视频弹幕列表
//	@Tags			danmu
//	@Accept			json
//	@Produce		json
//	@Param			video_id	path		int	true	"Video ID"
//	@Success		200			{object}	Success
//	@Failure		400			{object}	Fail
//	@Failure		404			{object}	NotFound
//	@Failure		500			{object}	ServerError
//	@Router			/danmu/getVideoBarrageList/{id} [get]
func GetVideoBarrageListApi(c *gin.Context) {
	var bind Common
	if err := c.ShouldBindUri(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	list, err := videoService.DanmuList(bind.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("弹幕列表", list))
}

type DanmuSave struct {
	VideoID uint    `form:"video_id" json:"video_id" binding:"required"`
	Text    string  `form:"text" json:"text" binding:"required"`
	Time    float64 `form:"time" json:"time"`
	Mode    uint8   `form:"mode" json:"mode"`
	Color   string  `form:"color" json:"color" binding:"required"`
	Border  bool    `form:"border" json:"border"`
	Style   string  `form:"style" json:"style"`
}

// SendVideoBarrageApi godoc
//
//	@Summary		视频弹幕
//	@Description	视频弹幕
//	@Tags			danmu
//	@Accept			json
//	@Produce		json
//	@Param			data	body		DanmuSave	true	"视频弹幕"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/danmu/sendVideoBarrage [post]
func SendVideoBarrageApi(c *gin.Context) {
	var bind DanmuSave
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)
	if err := videoService.DanmuSave(bind.VideoID, userID, bind.Text, bind.Time, bind.Mode, bind.Color, bind.Border, bind.Style); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("发送成功", nil))
}

type ImportVideoData struct {
	Dir       string `json:"dir" binding:"required"`
	Actresses string `json:"actresses" binding:"required"`
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
	if err := videoService.ImportVideoData(bind.Dir, bind.Actresses); err != nil {
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

func ResetTableApi(c *gin.Context) {
	var table = c.Query("table")
	if err := service.ResetTable(table); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}
