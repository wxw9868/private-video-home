package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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
	ActressID int    `uri:"actress_id" form:"actress_id" json:"actress_id"`
	Page      int    `uri:"page" form:"page" json:"page"`
	Size      int    `uri:"size" form:"size" json:"size"`
	Action    string `uri:"action" form:"action" json:"action" example:"v.CreatedAt"`
	Sort      string `uri:"sort" form:"sort" json:"sort" example:"desc"`
}

// GetVideoListApi godoc
//
//	@Summary		视频列表
//	@Description	视频列表
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			data	body		Video	true	"视频列表"
//	@Success		200		{object}	Success
//	@Failure		400		{object}	Fail
//	@Failure		404		{object}	NotFound
//	@Failure		500		{object}	ServerError
//	@Router			/video/getVideoList [post]
func GetVideoListApi(c *gin.Context) {
	var bind Video
	if err := c.BindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	data, err := videoService.Find(bind.ActressID, bind.Page, bind.Size, bind.Action, bind.Sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("视频列表", data))
}

type Common struct {
	ID uint `uri:"id" binding:"required"`
}

// VideoPlayApi godoc
//
//	@Summary		视频信息
//	@Description	get string by ID
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Video ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/video/videoPlay/{id} [get]
func VideoPlayApi(c *gin.Context) {
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
	Collect int  `form:"collect" json:"collect" binding:"required,oneof=1 -1"`
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

	userID := GetUserID(c)

	if err := videoService.Collect(bind.VideoID, bind.Collect, userID); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "收藏成功"
	if bind.Collect == -1 {
		msg = "取消收藏"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

// BrowseVideoApi godoc
//
//	@Summary		视频浏览记录
//	@Description	get string by ID
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			video_id	path		int	true	"Video ID"
//	@Success		200			{object}	Success
//	@Failure		400			{object}	Fail
//	@Failure		404			{object}	NotFound
//	@Failure		500			{object}	ServerError
//	@Router			/video/browseVideo/{id} [get]
func BrowseVideoApi(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	if aid == 0 {
		c.JSON(http.StatusBadRequest, util.Fail("video_id must be greater than 0"))
		return
	}

	userID := GetUserID(c)
	if err = videoService.Browse(uint(aid), userID); err != nil {
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
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	if aid == 0 {
		c.JSON(http.StatusBadRequest, util.Fail("video_id must be greater than 0"))
		return
	}

	userID := GetUserID(c)
	list, err := videoService.CommentList(cast.ToUint(aid), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("评论列表", list))
}

type CommentZan struct {
	CommentID uint `form:"comment_id" json:"comment_id" binding:"required"`
	Zan       int  `form:"zan" json:"zan" binding:"required,oneof=1 -1"`
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

	if err := videoService.Zan(bind.CommentID, userID, bind.Zan); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "点赞成功"
	if bind.Zan == -1 {
		msg = "取消点赞"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

type CommentCai struct {
	CommentID uint `form:"comment_id" json:"comment_id" binding:"required"`
	Cai       int  `form:"cai" json:"cai" binding:"required,oneof=1 -1"`
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

	if err := videoService.Cai(bind.CommentID, userID, bind.Cai); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	msg := "点踩成功"
	if bind.Cai == -1 {
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
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	if aid == 0 {
		c.JSON(http.StatusBadRequest, util.Fail("video_id must be greater than 0"))
		return
	}

	list, err := videoService.DanmuList(cast.ToUint(aid))
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
//	@Summary		导入视频数据
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			data	body		ImportVideoData	true	"视频数据"
//	@Success		200			{object}	Success
//	@Failure		400			{object}	Fail
//	@Failure		404			{object}	NotFound
//	@Failure		500			{object}	ServerError
//	@Router			/video/importVideo [post]
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
//	@Summary		视频导入GoFound
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/video/writeGoFound [get]
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
