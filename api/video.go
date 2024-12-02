package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

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

// VideoListApi godoc
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
//	@Router			/video/list [post]
func VideoListApi(c *gin.Context) {
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
//	@Router			/video/play/{id} [get]
func VideoPlayApi(c *gin.Context) {
	id := c.Param("id")
	aid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	if aid == 0 {
		c.JSON(http.StatusBadRequest, util.Fail("id must be greater than 0"))
		return
	}

	vi, err := videoService.Info(cast.ToUint(aid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	idsArr := strings.Split(vi.ActressIds, ",")
	actressArr := strings.Split(vi.ActressNames, ",")
	avatarArr := strings.Split(vi.ActressAvatars, ",")
	type actress struct {
		ID      string `json:"id"`
		Actress string `json:"actress"`
		Avatar  string `json:"avatar"`
	}
	actressSlice := make([]actress, len(actressArr))
	for i := 0; i < len(actressArr); i++ {
		actressSlice[i] = actress{ID: idsArr[i], Actress: actressArr[i], Avatar: avatarArr[i]}
	}
	var collectID uint = 0
	var isCollect = false
	usc, err := userService.CollectLog(GetUserID(c), vi.ID)
	if err == nil {
		collectID = usc.ID
		isCollect = true
	}

	size, _ := strconv.ParseFloat(strconv.FormatInt(vi.Size, 10), 64)
	data := gin.H{
		"videoID":       vi.ID,
		"videoTitle":    vi.Title,
		"videoActress":  actressSlice,
		"videoUrl":      "assets/video/" + vi.Title + ".mp4",
		"Size":          size / 1024 / 1024,
		"Duration":      utils.ResolveTime(uint32(vi.Duration)),
		"ModTime":       vi.CreationTime.Format("2006-01-02 15:04:05"),
		"Poster":        vi.Poster,
		"Width":         vi.Width,
		"Height":        vi.Height,
		"CodecName":     vi.CodecName,
		"ChannelLayout": vi.ChannelLayout,
		"Collect":       vi.Collect,
		"Browse":        vi.Browse,
		"Zan":           vi.Zan,
		"Cai":           vi.Cai,
		"Watch":         vi.Watch,
		"CollectID":     collectID,
		"IsCollect":     isCollect,
		"Avatar":        sessions.Default(c).Get("user_avatar").(string),
	}
	c.JSON(http.StatusOK, util.Success("视频信息", data))
}

type VideoCollect struct {
	VideoID uint `form:"video_id" json:"video_id" binding:"required"`
	Collect int  `form:"collect" json:"collect" binding:"required,oneof=1 -1"`
}

// VideoCollectApi godoc
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
//	@Router			/video/collect [post]
func VideoCollectApi(c *gin.Context) {
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

// VideoBrowseApi godoc
//
//	@Summary		视频浏览记录
//	@Description	get string by ID
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			video_id		path		int	true	"Video ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/video/browse/{video_id} [get]
func VideoBrowseApi(c *gin.Context) {
	id := c.Param("video_id")
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
	if err := videoService.Browse(uint(aid), userID); err != nil {
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
//	@Router			/comment/comment [post]
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

// VideoReplyApi godoc
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
//	@Router			/comment/reply [post]
func VideoReplyApi(c *gin.Context) {
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

// VideoCommentListApi godoc
//
//	@Summary		视频弹幕列表
//	@Description	视频弹幕列表
//	@Tags			comment
//	@Accept			json
//	@Produce		json
//	@Param			video_id		path		int	true	"Video ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/comment/list/{video_id} [get]
func VideoCommentListApi(c *gin.Context) {
	id := c.Param("video_id")
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

// CommentZanApi godoc
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
//	@Router			/comment/zan [post]
func CommentZanApi(c *gin.Context) {
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

// CommentCaiApi godoc
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
//	@Router			/comment/cai [post]
func CommentCaiApi(c *gin.Context) {
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

// DanmuListApi godoc
//
//	@Summary		视频弹幕列表
//	@Description	get string by ID
//	@Tags			danmu
//	@Accept			json
//	@Produce		json
//	@Param			video_id		path		int	true	"Video ID"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/danmu/list/{video_id} [get]
func DanmuListApi(c *gin.Context) {
	id := c.Param("video_id")
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

// DanmuSaveApi godoc
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
//	@Router			/danmu/save [post]
func DanmuSaveApi(c *gin.Context) {
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

// VideoImportApi godoc
//
//	@Summary		视频导入
//	@Description	get string by ID
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			dir			query		string		true	"dir"
//	@Param			actresss	query		string		true	"actresss"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/video/import [get]
func VideoImportApi(c *gin.Context) {
	var videoDir = c.Query("dir")
	var actresss = c.Query("actresss")
	if err := service.VideoImport(videoDir, actresss); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

// RepairVideoImportApi godoc
//
//	@Summary		修复视频导入
//	@Description	get string by ID
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Param			actresss	query		string		true	"actresss"
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/video/repairImport [get]
func RepairVideoImportApi(c *gin.Context) {
	var actresss = c.Query("actresss")
	if err := service.RepairVideoImport(strings.Split(actresss, ",")); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

// VideoWriteGoFound godoc
//
//	@Summary		视频导入GoFound
//	@Description	视频导入GoFound
//	@Tags			video
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Success
//	@Failure		400	{object}	Fail
//	@Failure		404	{object}	NotFound
//	@Failure		500	{object}	ServerError
//	@Router			/video/writeGoFound [get]
func VideoWriteGoFound(c *gin.Context) {
	if err := service.VideoWriteGoFound(); err != nil {
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

type VideoPic struct {
	Page int    `uri:"page" form:"page" json:"page"`
	Size int    `uri:"size" form:"size" json:"size"`
	Site string `uri:"site" form:"site" json:"site"`
}

func GetVideoPic(c *gin.Context) {
	var bind VideoPic
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := service.AllVideoPic(bind.Page, bind.Size, bind.Site); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

type VideoPic1 struct {
	Actress string `uri:"actress" form:"actress" json:"actress"`
	Site    string `uri:"site" form:"site" json:"site"`
}

func OneVideoPic(c *gin.Context) {
	var bind VideoPic1
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	if err := service.OneVideoPic(bind.Actress, bind.Site); err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}
