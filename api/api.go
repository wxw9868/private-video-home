package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/wxw9868/util"
	"github.com/wxw9868/util/pagination"
	"github.com/wxw9868/video/service"
	"github.com/wxw9868/video/utils"
)

func VideoIndex(c *gin.Context) {
	cards := make([]string, 14)
	for i := 0; i < 14; i++ {
		cardPath := "./assets/image/card/card" + strconv.Itoa(i+1) + ".jpeg"
		cards[i] = cardPath
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "首页",
		"data":  cards,
	})
}

type Search struct {
	Query     string      `json:"query" binding:"required"`
	Page      int         `json:"page"`
	Limit     int         `json:"limit"`
	Order     string      `json:"order"`
	Highlight interface{} `json:"highlight"`
	ScoreExp  string      `json:"scoreExp"`
}

func VideoSearch(c *gin.Context) {
	query := c.Query("query")

	b, err := json.Marshal(&Search{Query: query, Page: 1, Limit: 240, Order: "desc"})
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

	c.HTML(http.StatusOK, "search.html", gin.H{
		"title": "视频列表",
		"data":  string(robots),
	})
}

func VideoList(c *gin.Context) {
	actressID := cast.ToInt(c.Query("actress_id"))
	page := cast.ToInt(c.Query("page"))
	pageSize := cast.ToInt(c.Query("pagesize"))
	c.HTML(http.StatusOK, "video-list.html", gin.H{
		"title":     "视频列表",
		"actressID": actressID,
		"page":      page,
		"pagesize":  pageSize,
	})
}

func VideoActress(c *gin.Context) {
	c.HTML(http.StatusOK, "video-actress.html", gin.H{
		"title": "演员列表",
	})
}

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
	Action    string `uri:"action" form:"action" json:"action"`
	Sort      string `uri:"sort" form:"sort" json:"sort"`
}

func VideoListApi(c *gin.Context) {
	var bind Video
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	videos, count, err := vs.Find(bind.ActressID, bind.Page, bind.Size, bind.Action, bind.Sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	if bind.Size <= 0 {
		bind.Size = 16
	}
	pages := pagination.NewPaginator(int(count), bind.Size)
	pages.SetCurrentPage(bind.Page)

	c.JSON(http.StatusOK, util.Success("视频列表", map[string]interface{}{
		"list": videos,
		"page": map[string]interface{}{
			"totalPage":   pages.TotalPage(),
			"prePage":     pages.PrePage(),
			"currentPage": pages.CurrentPage(),
			"nextPage":    pages.NextPage(),
			"pages":       pages.Pages(),
		},
	}))
}

type Actress struct {
	Page   int    `uri:"page" form:"page" json:"page"`
	Size   int    `uri:"size" form:"size" json:"size"`
	Action string `uri:"action" form:"action" json:"action"`
	Sort   string `uri:"sort" form:"sort" json:"sort"`
}

func VideoActressApi(c *gin.Context) {
	var bind Video
	if err := c.Bind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	actresss, err := as.Find(bind.Page, bind.Size, bind.Action, bind.Sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	// var data = make(map[string]struct{})
	// for _, v := range actresss {
	// 	data[v.Actress] = struct{}{}
	// }
	// fmt.Printf("%d\n", len(data))
	// utils.WriteMapToFile("data.json", &data)
	// fmt.Printf("%+v\n", actresss)
	// fmt.Printf("%d\n", len(actresss))

	c.JSON(http.StatusOK, util.Success("视频列表", map[string]interface{}{
		"list": actresss,
	}))
}

type Play struct {
	ID string `form:"id" binding:"required"`
}

func VideoPlay(c *gin.Context) {
	var play Play
	if err := c.Bind(&play); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	vi, err := vs.Info(cast.ToUint(play.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	actressArr := strings.Split(vi.ActressStr, ",")
	idsArr := strings.Split(vi.ActressIds, ",")
	actressHtml := ""
	for i := 0; i < len(actressArr); i++ {
		if i == 0 {
			actressHtml += `<a href="/video/list?actress_id=` + idsArr[i] + `">` + actressArr[i] + `</a>`
		} else {
			actressHtml += `<a class="ms-2" href="/video/list?actress_id=` + idsArr[i] + `">` + actressArr[i] + `</a>`
		}
	}
	var collectID uint = 0
	var isCollect = false
	usc, err := us.CollectLog(GetUserID(c), vi.ID)
	if err == nil {
		collectID = usc.ID
		isCollect = true
	}

	size, _ := strconv.ParseFloat(strconv.FormatInt(vi.Size, 10), 64)
	c.HTML(http.StatusOK, "video-player.html", gin.H{
		"title":         "视频播放",
		"videoID":       vi.ID,
		"videoTitle":    vi.Title,
		"videoActress":  actressHtml,
		"videoUrl":      videoDir + "/" + vi.Title + ".mp4",
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
		"Avatar":        sessions.Default(c).Get("userAvatar").(string),
	})
}

func VideoPlayApi(c *gin.Context) {
	var play Play
	if err := c.Bind(&play); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		return
	}

	vi, err := vs.Info(cast.ToUint(play.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.Fail(err.Error()))
		return
	}

	actressArr := strings.Split(vi.ActressStr, ",")
	idsArr := strings.Split(vi.ActressIds, ",")
	type actress struct {
		ID      string `json:"id"`
		Actress string `json:"actress"`
	}
	actressSlice := make([]actress, len(actressArr))
	for i := 0; i < len(actressArr); i++ {
		actressSlice[i] = actress{ID: idsArr[i], Actress: actressArr[i]}
	}
	var collectID uint = 0
	var isCollect = false
	usc, err := us.CollectLog(GetUserID(c), vi.ID)
	if err == nil {
		collectID = usc.ID
		isCollect = true
	}

	size, _ := strconv.ParseFloat(strconv.FormatInt(vi.Size, 10), 64)
	c.JSON(http.StatusOK, gin.H{
		"videoID":       vi.ID,
		"videoTitle":    vi.Title,
		"videoActress":  actressSlice,
		"videoUrl":      videoDir + "/" + vi.Title + ".mp4",
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
		"Avatar":        sessions.Default(c).Get("userAvatar").(string),
	})
}

type VideoCollect struct {
	VideoID uint `form:"video_id" json:"video_id" binding:"required"`
	Collect int  `form:"collect" json:"collect" binding:"required,oneof=1 -1"`
}

// 收藏
func VideoCollectApi(c *gin.Context) {
	var bind VideoCollect
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	if err := vs.Collect(bind.VideoID, bind.Collect, userID); err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
		return
	}

	msg := "收藏成功"
	if bind.Collect == -1 {
		msg = "取消收藏"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

type VideoBrowse struct {
	VideoID uint `form:"video_id" json:"video_id" binding:"required"`
}

// 浏览
func VideoBrowseApi(c *gin.Context) {
	var bind VideoBrowse
	if err := c.ShouldBind(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	if err := vs.Browse(bind.VideoID, userID); err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
		return
	}

	msg := "浏览记录成功"
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

type VideoComment struct {
	VideoID uint   `form:"video_id" json:"video_id" binding:"required"`
	Content string `form:"content" json:"content" binding:"required"`
}

// 评论
func VideoCommentApi(c *gin.Context) {
	var bind VideoComment
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	commentID, err := vs.Comment(bind.VideoID, bind.Content, userID)
	if err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	userAvatar := session.Get("userAvatar").(string)
	userNickname := session.Get("userNickname").(string)

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

// 回复
func VideoReplyApi(c *gin.Context) {
	var bind VideoReply
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	commentID, err := vs.Reply(bind.VideoID, bind.ParentID, bind.Content, userID)
	if err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
		return
	}

	session := sessions.Default(c)
	userAvatar := session.Get("userAvatar").(string)
	userNickname := session.Get("userNickname").(string)

	data := map[string]interface{}{
		"commentID":    commentID,
		"userAvatar":   userAvatar,
		"userNickname": userNickname,
		"content":      bind.Content,
	}

	c.JSON(http.StatusOK, util.Success("回复成功", data))
}

func VideoCommentListApi(c *gin.Context) {
	id := c.Query("video_id")

	userID := GetUserID(c)

	list, err := vs.CommentList(cast.ToUint(id), userID)
	if err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Success("评论列表", list))
}

type CommentZan struct {
	CommentID uint `form:"comment_id" json:"comment_id" binding:"required"`
	Zan       int  `form:"zan" json:"zan" binding:"required,oneof=1 -1"`
}

// 赞
func CommentZanApi(c *gin.Context) {
	var bind CommentZan
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	if err := vs.Zan(bind.CommentID, userID, bind.Zan); err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
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

// 踩
func CommentCaiApi(c *gin.Context) {
	var bind CommentCai
	if err := c.ShouldBindJSON(&bind); err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}

	userID := GetUserID(c)

	if err := vs.Cai(bind.CommentID, userID, bind.Cai); err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
		return
	}

	msg := "点踩成功"
	if bind.Cai == -1 {
		msg = "取消踩"
	}
	c.JSON(http.StatusOK, util.Success(msg, nil))
}

func VideoRename(c *gin.Context) {
	var videoDir = c.Query("dir")
	var nameMap = map[string]string{
		"(1)":  "",
		"(2)":  "",
		"(3)":  "",
		"(4)":  "",
		"(5)":  "",
		"(6)":  "",
		"(7)":  "",
		"(8)":  "",
		"(9)":  "",
		"(10)": "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新":      "082118-735 女熱大陸 File.064  #真菜果",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (2)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (3)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (4)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (5)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (6)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (7)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (8)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (9)":  "",
		"无码频道-tg关注 @AVWUMAYUANPIAN  每天更新 (10)": "",
	}

	files, err := os.ReadDir(videoDir)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.Fail(err.Error()))
		return
	}
	for _, file := range files {
		filename := file.Name()
		oldPath := videoDir + "/" + filename

		oldName := strings.Split(filename, ".")[0]
		newName, ok := nameMap[oldName]
		if ok {
			filename = strings.Replace(filename, oldName, newName, -1)
		} else {
			filename = strings.Replace(filename, "无码频道_tg关注_@AVWUMAYUANPIAN_每天更新_", "", -1)
			filename = strings.Replace(filename, "_tg关注_@AVWUMAYUANPIAN", "", -1)
		}

		newPath := videoDir + "/" + filename
		os.Rename(oldPath, newPath)
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

func VideoImport(c *gin.Context) {
	var videoDir = c.Query("dir")
	if err := service.VideoImport(videoDir); err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}

func ImportVideoActress(c *gin.Context) {
	if err := service.VideoActress(); err != nil {
		c.JSON(http.StatusOK, util.Fail(err.Error()))
		return
	}

	c.JSON(http.StatusOK, util.Success("SUCCESS", nil))
}
