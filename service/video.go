package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/spf13/cast"

	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/model/request"
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

type VideoService struct{}

type Video struct {
	ID       uint    `json:"id"`
	Title    string  `gorm:"column:title;type:varchar(255);uniqueIndex;comment:标题" json:"title"`
	Duration float64 `gorm:"column:duration;type:float;default:0;comment:时长" json:"duration"`
	Poster   string  `gorm:"column:poster;type:varchar(255);comment:封面" json:"poster"`
	Collect  uint    `gorm:"column:collection_volume;type:uint;not null;default:0;comment:收藏" json:"collect"`
	Browse   uint    `gorm:"column:page_views;type:uint;not null;default:0;comment:浏览" json:"browse"`
}

func (vs *VideoService) List(req request.SearchVideo) (data map[string]interface{}, err error) {
	var ids []uint
	var key string
	switch req.Column {
	case "v.CreatedAt":
		key = "video_video_createdAt"
	case "l.browse":
		key = "video_video_browse"
	case "l.collect":
		key = "video_video_collect"
	default:
		key = "video_video"
	}

	f := func(ids []uint, total int) (map[string]interface{}, error) {
		videos := make([]Video, len(ids))
		for i, id := range ids {
			result := rdb.HGetAll(ctx, utils.Join("video_video_", strconv.FormatUint(uint64(id), 10))).Val()
			videos[i] = Video{
				ID:       id,
				Title:    result["title"],
				Poster:   result["poster"],
				Duration: cast.ToFloat64(result["duration"]),
				Browse:   cast.ToUint(result["browseNum"]),
				Collect:  cast.ToUint(result["collectNum"]),
			}
		}
		return map[string]interface{}{"list": videos, "total": total}, nil
	}

	if req.ActressID != 0 {
		var vadb = db.Model(&model.VideoActress{}).Where("actress_id = ?", req.ActressID)
		var total int64
		if err = vadb.Count(&total).Error; err != nil {
			return nil, err
		}
		if err = vadb.Scopes(Paginate(req.Page, req.Size, int(total))).Pluck("video_id", &ids).Error; err != nil {
			return nil, err
		}
		return f(ids, int(total))
	}

	vdb := db.Table("video_Video as v").
		Select("v.id, v.title, v.duration, v.poster, l.collection_volume, l.page_views").
		Joins("left join video_VideoLog l on l.video_id = v.id")
	if req.Column != "" && req.Order != "" {
		vdb = vdb.Order(utils.Join(req.Column, " ", req.Order))
	}
	var total int64
	if err = vdb.Count(&total).Error; err != nil {
		return nil, err
	}
	if err = db.Table("(?)", vdb).Scopes(Paginate(req.Page, req.Size, int(total))).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	bts, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}
	result, _ := rdb.HGet(ctx, key, "ids").Result()
	if strings.Compare(string(bts), result) == 0 {
		return f(ids, int(total))
	}

	var videos []Video
	err = vdb.Scopes(Paginate(req.Page, req.Size, int(total))).Scan(&videos).Error
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, video := range videos {
		keys = append(keys, utils.Join("video_video_", strconv.FormatUint(uint64(video.ID), 10)))
	}

	txf := func(tx *redis.Tx) error {
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, key, "len", len(ids), "ids", string(bts))
			for _, video := range videos {
				pipe.HSet(ctx, utils.Join("video_video_", strconv.Itoa(int(video.ID))), "id", video.ID, "title", video.Title, "poster", video.Poster, "duration", video.Duration, "browseNum", video.Browse, "collectNum", video.Collect)
			}
			return nil
		})
		return err
	}
	if err = rdb.Watch(ctx, txf, keys...); errors.Is(err, redis.TxFailedErr) {
		return nil, err
	}

	return map[string]interface{}{"list": videos, "total": total}, nil
}

type VideoInfo struct {
	ID               uint      `json:"id"`
	Title            string    `gorm:"column:title;type:varchar(255);uniqueIndex;comment:标题" json:"title"`
	Size             int64     `gorm:"column:size;type:bigint;comment:大小" json:"size"`
	Duration         float64   `gorm:"column:duration;type:float;default:0;comment:时长" json:"duration"`
	Poster           string    `gorm:"column:poster;type:varchar(255);comment:封面" json:"poster"`
	Width            int       `gorm:"column:width;type:int;default:0;comment:宽" json:"width"`
	Height           int       `gorm:"column:height;type:int;default:0;comment:高" json:"height"`
	CodecName        string    `gorm:"column:codec_name;type:varchar(90);comment:编解码器" json:"codec_name"`
	ChannelLayout    string    `gorm:"column:channel_layout;type:varchar(90);comment:音频声道" json:"channel_layout"`
	CreationTime     time.Time `gorm:"column:creation_time;type:date;comment:时间" json:"creation_time"`
	CollectionVolume uint      `gorm:"column:collection_volume;type:uint;not null;default:0;comment:收藏量"`
	PageViews        uint      `gorm:"column:page_views;type:uint;not null;default:0;comment:浏览量"`
	LikesCount       uint      `gorm:"column:likes_count;type:uint;not null;default:0;comment:点赞量"`
	DislikesCount    uint      `gorm:"column:dislikes_count;type:uint;not null;default:0;comment:点踩量"`
	ViewsCount       uint      `gorm:"column:views_count;type:uint;not null;default:0;comment:观看次数"`
	ActressIds       string    `gorm:"column:actress_ids;comment:演员ID" json:"actressIds"`
	ActressNames     string    `gorm:"column:actress_names;comment:演员名称" json:"actressNames"`
	ActressAvatars   string    `gorm:"column:actress_avatars;comment:演员头像" json:"actressAvatars"`
}

func (vs *VideoService) Info(id, userId uint) (map[string]interface{}, error) {
	var videoInfo VideoInfo
	if err := db.Table("video_Video as v").
		Select("v.id,v.title,v.duration,v.poster,v.size,v.width,v.height,v.codec_name,v.channel_layout,v.creation_time,l.collection_volume, l.page_views, l.likes_count, l.dislikes_count, l.views_count, group_concat(a.id,',') as actress_ids, group_concat(a.actress,',') as actress_names, group_concat(a.avatar,',') as actress_avatars").
		Joins("left join video_VideoLog as l on l.video_id = v.id").
		Joins("left join video_VideoActress as va on va.video_id = v.id").
		Joins("left join video_Actress as a on a.id = va.actress_id").
		Where("v.id = ?", id).
		Group("v.id,v.title,v.duration,v.poster,v.size,v.width,v.height,v.codec_name,v.channel_layout,v.creation_time,l.collection_volume, l.page_views, l.likes_count, l.dislikes_count, l.views_count").
		Scan(&videoInfo).Error; err != nil {
		return nil, err
	}

	type actress struct {
		ID      string `json:"id"`
		Actress string `json:"actress"`
		Avatar  string `json:"avatar"`
	}
	idSlice := strings.Split(videoInfo.ActressIds, ",")
	actressSlice := strings.Split(videoInfo.ActressNames, ",")
	avatarSlice := strings.Split(videoInfo.ActressAvatars, ",")
	actresses := make([]actress, len(idSlice))
	for i := 0; i < len(idSlice); i++ {
		actresses[i] = actress{ID: idSlice[i], Actress: actressSlice[i], Avatar: avatarSlice[i]}
	}
	var collectID uint = 0
	var isCollect = false
	userCollectLog, err := new(UserService).CollectLog(userId, videoInfo.ID)
	if err == nil {
		collectID = userCollectLog.ID
		isCollect = true
	}
	data := map[string]interface{}{
		"id":            videoInfo.ID,
		"title":         videoInfo.Title,
		"actresses":     actresses,
		"poster":        videoInfo.Poster,
		"link":          "assets/video/" + videoInfo.Title + ".mp4",
		"duration":      utils.ResolveTime(uint32(videoInfo.Duration)),
		"size":          float64(videoInfo.Size) / 1024 / 1024,
		"width":         videoInfo.Width,
		"height":        videoInfo.Height,
		"codecName":     videoInfo.CodecName,
		"channelLayout": videoInfo.ChannelLayout,
		"modTime":       videoInfo.CreationTime.Format("2006-01-02 15:04:05"),
		"collect":       videoInfo.CollectionVolume,
		"browse":        videoInfo.PageViews,
		"like":          videoInfo.LikesCount,
		"dislike":       videoInfo.DislikesCount,
		"watch":         videoInfo.ViewsCount,
		"collectID":     collectID,
		"isCollect":     isCollect,
	}
	return data, nil
}

func (vs *VideoService) Create(videos []model.Video) error {
	if err := db.Create(&videos).Error; err != nil {
		return err
	}
	return nil
}

func (vs *VideoService) Collect(videoID, userID uint, num int) error {
	var video model.Video
	if errors.Is(db.First(&video, videoID).Error, gorm.ErrRecordNotFound) {
		return errors.New("视频不存在！")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var expr string
		if num == 1 {
			// 增加1
			expr = "collection_volume + 1"
			if err := tx.Create(&model.UserCollectLog{UserID: userID, VideoID: videoID}).Error; err != nil {
				return errors.New("创建失败！")
			}
		} else {
			// 减少1
			expr = "collection_volume - 1"
			if err := tx.Where("user_id = ? and video_id = ?", userID, videoID).Delete(&model.UserCollectLog{}).Error; err != nil {
				return errors.New("删除失败！")
			}
		}
		if err := tx.Model(&model.VideoLog{}).Where("video_id = ?", videoID).Update("collection_volume", gorm.Expr(expr)).Error; err != nil {
			return errors.New("更新失败！")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return vs.GofoundIndex(videoID)
}

func (vs *VideoService) Browse(videoID uint, userID uint) error {
	var video model.Video
	if errors.Is(db.First(&video, videoID).Error, gorm.ErrRecordNotFound) {
		return errors.New("视频不存在！")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var userBrowseLog model.UserPageViewsLog
		if err := tx.Where(model.UserPageViewsLog{UserID: userID, VideoID: videoID}).FirstOrInit(&userBrowseLog).Error; err != nil {
			return err
		}

		if err := tx.Where(model.UserPageViewsLog{UserID: userID, VideoID: videoID}).Assign(model.UserPageViewsLog{PageViews: userBrowseLog.PageViews + 1}).FirstOrCreate(&model.UserPageViewsLog{}).Error; err != nil {
			return fmt.Errorf("创建失败: %s", err)
		}

		var videoLog model.VideoLog
		if err := tx.Where(model.VideoLog{VideoID: videoID}).FirstOrInit(&videoLog).Error; err != nil {
			return err
		}

		if err := tx.Where(model.VideoLog{VideoID: videoID}).Assign(model.VideoLog{PageViews: videoLog.PageViews + 1}).FirstOrCreate(&model.VideoLog{}).Error; err != nil {
			return fmt.Errorf("创建失败: %s", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return vs.GofoundIndex(videoID)
}

func (vs *VideoService) GofoundIndex(videoID uint) error {
	data, _ := vs.Info(videoID, 0)
	id := data["id"].(uint)
	title := data["title"].(string)

	rdb.HSet(ctx, utils.Join("video_video_", strconv.FormatUint(uint64(id), 10)), "id", id, "title", title, "poster", data["poster"], "duration", data["duration"], "browseNum", data["browse"], "collectNum", data["collect"])
	index := Index{
		Id:   id,
		Text: title,
		Document: VideoData{
			ID:            data["id"].(uint),
			Title:         data["title"].(string),
			Poster:        data["poster"].(string),
			Duration:      data["duration"].(string),
			Size:          data["size"].(float64),
			CreationTime:  data["modTime"].(string),
			Width:         data["width"].(int),
			Height:        data["height"].(int),
			CodecName:     data["codecName"].(string),
			ChannelLayout: data["channelLayout"].(string),
			CollectNum:    data["collect"].(uint),
			BrowseNum:     data["browse"].(uint),
			LikeNum:       data["like"].(uint),
			DisLikeNum:    data["dislike"].(uint),
			WatchNum:      data["watch"].(uint),
		},
	}
	b, err := json.Marshal(&index)
	if err != nil {
		return err
	}
	if err = Post(utils.Join("/index", "?", "database=", "private-video"), bytes.NewReader(b)); err != nil {
		return err
	}
	return nil
}

func (vs *VideoService) Comment(user *model.User, videoID uint, content string) (map[string]interface{}, error) {
	comment := model.VideoComment{
		ParentId:    0,
		VideoId:     videoID,
		UserId:      user.ID,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      "APPROVED",
		IsAnonymous: 1,
		Content:     content,
		IsShow:      1,
	}
	if err := db.Create(&comment).Error; err != nil {
		return nil, err
	}

	data := map[string]interface{}{"id": comment.ID, "avatar": user.Avatar, "nickname": user.Nickname, "content": content, "createdAt": comment.CreatedAt}
	return data, nil
}

func (vs *VideoService) Reply(user *model.User, videoID, parentID uint, content string) (map[string]interface{}, error) {
	comment := model.VideoComment{
		ParentId:    parentID,
		VideoId:     videoID,
		UserId:      user.ID,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      "APPROVED",
		IsAnonymous: 1,
		Content:     content,
		IsShow:      1,
	}

	db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&comment).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.VideoComment{}).Where("id = ?", parentID).Update("reply_num", gorm.Expr("reply_num + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	data := map[string]interface{}{"commentID": comment.ID, "userAvatar": user.Avatar, "userNickname": user.Nickname, "content": content}
	return data, nil
}

type VideoComment struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `gorm:"comment:" json:"createdAt"`
	ParentId      uint      `gorm:"column:parent_id;type:uint;not null;default:0;comment:父级评论的ID" json:"parentId"`
	VideoId       uint      `gorm:"column:video_id;type:uint;not null;default:0;comment:被评论的视频ID" json:"videoId"`
	UserId        uint      `gorm:"column:user_id;type:uint;not null;default:0;comment:评论人的ID" json:"userId"`
	Nickname      string    `gorm:"column:nickname;type:varchar(13);null;comment:评论人的昵称" json:"nickname"`
	Avatar        string    `gorm:"column:avatar;type:varchar(255);null;comment:评论人的头像地址" json:"avatar"`
	Content       string    `gorm:"column:content;type:text;not null;comment:评论内容" json:"content"`
	LikesCount    uint      `gorm:"column:likes_count;type:uint;not null;default:0;comment:点赞量" json:"likesCount"`
	DislikesCount uint      `gorm:"column:dislikes_count;type:uint;not null;default:0;comment:点踩量" json:"dislikesCount"`
	Like          int8      `gorm:"column:like;type:uint;not null;default:0;comment:支持（赞）" json:"like"`
	Dislike       int8      `gorm:"column:dislike;type:uint;not null;default:0;comment:反对（踩）" json:"dislike"`
	LogUserID     uint      `gorm:"column:log_user_id;comment:用户ID" json:"logUserId"`
}

func (vs *VideoService) CommentList(userID, videoID uint) ([]*CommentTree, error) {
	var list []VideoComment
	query := db.Model(&model.UserCommentLog{}).Where("video_id = ? and user_id = ?", videoID, userID)
	err := db.Table("video_VideoComment as c").
		Select("c.id, c.CreatedAt, c.parent_id, c.video_id, c.user_id, c.nickname, c.avatar, c.content, c.likes_count, c.dislikes_count, l.like, l.dislike, l.user_id as log_user_id").
		Joins("left join (?) l on l.comment_id = c.id", query).
		Where("c.video_id = ?", videoID).
		Order("c.id desc").Find(&list).Error
	if err != nil {
		return nil, err
	}
	return tree(list), nil
}

func (vs *VideoService) LikeVideoComment(userID, commentID uint, like int8) error {
	var comment model.VideoComment
	if errors.Is(db.First(&comment, commentID).Error, gorm.ErrRecordNotFound) {
		return errors.New("评论不存在！")
	}

	// 减少1
	expr := "likes_count - 1"
	if like == 1 {
		// 增加1
		expr = "likes_count + 1"
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(model.UserCommentLog{UserID: userID, VideoID: comment.VideoId, CommentID: commentID}).Assign(model.UserCommentLog{Like: &like}).FirstOrCreate(&model.UserCommentLog{}).Error; err != nil {
			return fmt.Errorf("创建失败: %s", err)
		}

		if err := tx.Model(&model.VideoComment{}).Where("id = ? and video_id = ?", commentID, comment.VideoId).Update("likes_count", gorm.Expr(expr)).Error; err != nil {
			return fmt.Errorf("更新失败: %s", err)
		}

		return nil
	})
}

func (vs *VideoService) DislikeVideoComment(userID, commentID uint, dislike int8) error {
	var comment model.VideoComment
	if errors.Is(db.First(&comment, commentID).Error, gorm.ErrRecordNotFound) {
		return errors.New("评论不存在！")
	}

	// 减少1
	expr := "dislikes_count - 1"
	if dislike == 1 {
		// 增加1
		expr = "dislikes_count + 1"
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(model.UserCommentLog{UserID: userID, VideoID: comment.VideoId, CommentID: commentID}).Assign(model.UserCommentLog{Dislike: &dislike}).FirstOrCreate(&model.UserCommentLog{}).Error; err != nil {
			return fmt.Errorf("创建失败: %s", err)
		}

		if err := tx.Model(&model.VideoComment{}).Where("id = ? and video_id = ?", commentID, comment.VideoId).Update("dislikes_count", gorm.Expr(expr)).Error; err != nil {
			return fmt.Errorf("更新失败: %s", err)
		}

		return nil
	})
}

type CommentTree struct {
	VideoComment
	Childrens []CommentTree `json:"childrens"`
}

func tree(list []VideoComment) []*CommentTree {
	var data = make(map[uint]*CommentTree)
	var childrens = make(map[uint][]CommentTree)
	var dataSort []uint
	var childrensSort []uint
	for _, v := range list {
		if v.ParentId == 0 {
			data[v.ID] = &CommentTree{v, nil}
			dataSort = append(dataSort, v.ID)
		} else {
			childrens[v.ParentId] = append(childrens[v.ParentId], CommentTree{v, nil})
			childrensSort = append(childrensSort, v.ParentId)
		}
	}

	trees := recursiveSort(data, childrens, dataSort, childrensSort)

	result := make([]*CommentTree, len(trees))
	for k, v := range dataSort {
		result[k] = trees[v]
	}

	return result
}

func recursiveSort(data map[uint]*CommentTree, childrens map[uint][]CommentTree, dataSort, childrensSort []uint) map[uint]*CommentTree {
	for _, v := range dataSort {
		videoComments, ok := childrens[v]
		if ok {
			data[v].Childrens = videoComments
			delete(childrens, v)
			childrensSort = deleteArray(childrensSort, v)
			if len(childrens) > 0 {
				data := make(map[uint]*CommentTree, len(videoComments))
				dataSort := make([]uint, len(videoComments))
				for k, v := range videoComments {
					videoComment := v
					data[v.ID] = &videoComment
					dataSort[k] = v.ID
				}
				recursiveSort(data, childrens, dataSort, childrensSort)
			}
		}
	}
	return data
}

func deleteArray(d []uint, e uint) []uint {
	r := make([]uint, len(d)-1)
	j := 0
	for i := 0; i < len(d); i++ {
		if d[i] != e {
			r[j] = d[i]
			j++
		}
	}
	return r
}

type VideoDanmu struct {
	Text   string  `json:"text" gorm:"comment:弹幕文本"`
	Time   float64 `json:"time" gorm:"comment:弹幕时间, 默认为当前播放器时间"`
	Mode   uint8   `json:"mode" gorm:"comment:弹幕模式: 0: 滚动(默认) 1: 顶部 2: 底部"`
	Color  string  `json:"color" gorm:"comment:弹幕颜色，默认为白色"`
	Border bool    `json:"border" gorm:"comment:弹幕是否有描边, 默认为 false"`
	Style  string  `json:"style" gorm:"comment:弹幕自定义样式, 默认为空对象"`
}

func (vs *VideoService) DanmuList(videoID uint) ([]VideoDanmu, error) {
	var list []VideoDanmu
	if err := db.Model(&model.VideoDanmu{}).Where("video_id = ?", videoID).Order("time asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (vs *VideoService) Danmu(userID uint, req request.CreateDanmu) error {
	danmu := model.VideoDanmu{
		VideoId: req.VideoID,
		Text:    req.Text,
		Time:    req.Time,
		Mode:    req.Mode,
		Color:   req.Color,
		Border:  req.Border,
		Style:   req.Style,
		UserId:  userID,
	}
	if err := db.Create(&danmu).Error; err != nil {
		return err
	}
	return nil
}

func (vs *VideoService) ImportVideoData(dir string, actresses ...string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var ids []uint
	db.Model(&model.Video{}).Order("id desc").Limit(1).Pluck("id", &ids)

	videoSQL := generateVideoSQL(dir, files)
	actressSQL := generateActressSQL(actresses)

	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Exec(videoSQL).Error; err != nil {
			return err
		}

		if actressSQL != "" {
			if err = tx.Exec(actressSQL).Error; err != nil {
				return err
			}
		}

		var videoActressSQL = "INSERT OR REPLACE INTO video_VideoActress (video_id, actress_id, CreatedAt, UpdatedAt) VALUES "
		for _, v := range actresses {

			var actress model.Actress
			if err = tx.Where("actress = ?", v).First(&actress).Error; err != nil {
				return err
			}

			var videos []model.Video
			if err = tx.Where("title LIKE ?", "%"+v+"%").Find(&videos).Error; err != nil {
				return err
			}

			for _, video := range videos {
				var videoActress model.VideoActress
				err = tx.Model(&model.VideoActress{}).Where("video_id = ? and actress_id = ?", video.ID, actress.ID).First(&videoActress).Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					videoActressSQL += fmt.Sprintf("(%d, %d, '%v', '%v'), ", video.ID, actress.ID, time.Now().Local(), time.Now().Local())
				}
			}
		}

		b := []byte(videoActressSQL)
		if err = tx.Exec(string(b[:len(b)-2])).Error; err != nil {
			return err
		}
		return nil
	})

	if err == nil {
		oldId := ids[0]
		db.Model(&model.Video{}).Order("id desc").Limit(1).Pluck("id", &ids)
		newId := ids[0]
		if err = VideoWriteGoFound(fmt.Sprintf("v.id between %d and %d", oldId, newId)); err != nil {
			return err
		}
	}

	return err
}

func generateVideoSQL(dir string, files []os.DirEntry) string {
	var videoSQL = "INSERT OR REPLACE INTO video_Video (title, actress, size, duration, poster, width, height, codec_name, channel_layout, creation_time, CreatedAt, UpdatedAt) VALUES "
	for _, file := range files {
		filename := file.Name()
		if filepath.Ext(filename) == ".mp4" {
			title := strings.Split(filename, ".")[0]
			array := strings.Split(title, "_")
			actress := array[len(array)-1]
			posterPath := "assets/image/poster/" + title + ".jpg"
			videoPath := dir + "/" + filename
			videoInfo, _ := utils.GetVideoInfo(videoPath)
			size := videoInfo["size"].(int64)
			duration := videoInfo["duration"].(float64)
			width := videoInfo["width"].(int64)
			height := videoInfo["height"].(int64)
			codec := fmt.Sprintf("%s,%s", videoInfo["codec_name0"].(string), videoInfo["codec_name1"].(string))
			channelLayout := videoInfo["channel_layout"].(string)
			creationTime := videoInfo["creation_time"].(time.Time)
			videoSQL += fmt.Sprintf("('%s', '%s', %d, %f, '%s', %d, %d, '%s', '%s', '%v', '%v', '%v'), ", title, actress, size, duration, posterPath, width, height, codec, channelLayout, creationTime, time.Now().Local(), time.Now().Local())
		}
	}
	b := []byte(videoSQL)
	return string(b[:len(b)-2])
}

func generateActressSQL(actresses []string) string {
	if len(actresses) <= 0 {
		return ""
	}
	var actressSQL = "INSERT OR REPLACE INTO video_Actress (actress, avatar, CreatedAt, UpdatedAt) VALUES "
	startLen := len(actressSQL)
	for _, actress := range actresses {
		var data model.Actress
		if errors.Is(db.Model(&model.Actress{}).Where("actress = ?", actress).First(&data).Error, gorm.ErrRecordNotFound) {
			//_, err = os.Stat("assets/image/avatar/" + actress + ".jpg")
			//if os.IsNotExist(err) {
			//	nameSlice := []rune(actress)
			//	if err := utils.GenerateAvatar(string(nameSlice[0]), avatarPath); err != nil {
			//		return err
			//	}
			//}
			avatarPath := "assets/image/avatar/defaultgirl.png"
			actressSQL += fmt.Sprintf("('%s', '%s', '%v', '%v'), ", actress, avatarPath, time.Now().Local(), time.Now().Local())
		}
	}
	endLen := len(actressSQL)
	if endLen <= startLen {
		return ""
	}

	b := []byte(actressSQL)
	return string(b[:len(b)-2])
}
