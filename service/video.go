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
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

type VideoService struct{}

type VideoInfo struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title" gorm:"comment:标题"`
	Actress        string    `json:"actress" gorm:"comment:演员"`
	Size           int64     `json:"size" gorm:"comment:大小"`
	Duration       float64   `json:"duration" gorm:"comment:时长"`
	Poster         string    `json:"poster" gorm:"comment:封面"`
	Width          int       `json:"width" gorm:"comment:宽"`
	Height         int       `json:"height" gorm:"comment:高"`
	CodecName      string    `json:"codecName" gorm:"comment:编解码器"`
	ChannelLayout  string    `json:"channelLayout" gorm:"comment:音频声道"`
	CreationTime   time.Time `json:"creationTime" gorm:"comment:时间"`
	CollectNum     uint      `json:"collectNum" gorm:"comment:收藏"`
	BrowseNum      uint      `json:"browseNum" gorm:"comment:浏览"`
	LikeNum        uint      `json:"likeNum" gorm:"comment:赞"`
	DisLikeNum     uint      `json:"disLikeNum" gorm:"comment:踩"`
	WatchNum       uint      `json:"watchNum" gorm:"comment:观看"`
	ActressIds     string    `json:"actressIds" gorm:"comment:演员ID"`
	ActressNames   string    `json:"actressNames" gorm:"comment:演员名称"`
	ActressAvatars string    `json:"actressAvatars" gorm:"comment:演员头像"`
}
type Video struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Poster     string `json:"poster"`
	Duration   string `json:"duration"`
	BrowseNum  uint   `json:"browseNum"`
	CollectNum uint   `json:"collectNum"`
}

func (vs *VideoService) Find(actressID int, page, pageSize int, action, sort string) (data map[string]interface{}, err error) {
	var ids []uint

	f := func(ids []uint, total int) (map[string]interface{}, error) {
		videos := make([]Video, len(ids))
		for i, id := range ids {
			result := rdb.HGetAll(ctx, utils.Join("video_video_", strconv.Itoa(int(id)))).Val()
			videos[i] = Video{
				ID:         id,
				Title:      result["title"],
				Poster:     result["poster"],
				Duration:   result["duration"],
				BrowseNum:  cast.ToUint(result["browseNum"]),
				CollectNum: cast.ToUint(result["collectNum"]),
			}
		}
		data = map[string]interface{}{
			"list":  videos,
			"total": total,
		}
		return data, nil
	}

	if actressID != 0 {
		var vadb = db.Model(&model.VideoActress{}).Where("actress_id = ?", actressID)
		var total int64
		if err = vadb.Count(&total).Error; err != nil {
			return nil, err
		}
		if err = vadb.Scopes(Paginate(page, pageSize, int(total))).Pluck("video_id", &ids).Error; err != nil {
			return nil, err
		}
		return f(ids, int(total))
	}

	var key string
	switch action {
	case "v.CreatedAt":
		key = "video_video_createdAt"
	case "l.browse":
		key = "video_video_browse"
	case "l.collect":
		key = "video_video_collect"
	default:
		key = "video_video"
	}

	var vdb = db.Table("video_Video as v")
	vdb = vdb.Select("v.*,l.collect, l.browse, l.zan, l.cai, l.watch").Joins("left join video_VideoLog l on l.video_id = v.id")
	if action != "" && sort != "" {
		vdb = vdb.Order(utils.Join(action, " ", sort))
	}
	var total int64
	if err = vdb.Count(&total).Error; err != nil {
		return nil, err
	}
	if err = db.Table("(?)", vdb).Scopes(Paginate(page, pageSize, int(total))).Pluck("id", &ids).Error; err != nil {
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

	rows, err := vdb.Scopes(Paginate(page, pageSize, int(total))).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []Video
	var keys []string
	keys = append(keys, key)
	for rows.Next() {
		var videoInfo VideoInfo
		db.ScanRows(rows, &videoInfo)

		video := Video{
			ID:         videoInfo.ID,
			Title:      videoInfo.Title,
			Poster:     videoInfo.Poster,
			Duration:   utils.ResolveTime(uint32(videoInfo.Duration)),
			BrowseNum:  videoInfo.BrowseNum,
			CollectNum: videoInfo.CollectNum,
		}
		videos = append(videos, video)
		keys = append(keys, utils.Join("video_video_", strconv.Itoa(int(videoInfo.ID))))
	}

	txf := func(tx *redis.Tx) error {
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, key, "len", len(ids), "ids", string(bts))
			for _, video := range videos {
				pipe.HSet(ctx, utils.Join("video_video_", strconv.Itoa(int(video.ID))), "id", video.ID, "title", video.Title, "poster", video.Poster, "duration", video.Duration, "browseNum", video.BrowseNum, "collectNum", video.CollectNum)
			}
			return nil
		})
		return err
	}
	if err = rdb.Watch(ctx, txf, keys...); errors.Is(err, redis.TxFailedErr) {
		return nil, err
	}

	data = map[string]interface{}{
		"list":  videos,
		"total": total,
	}
	return data, nil
}

func (vs *VideoService) First(id string) (model.Video, error) {
	var video model.Video
	if err := db.Where("id = ?", id).First(&video).Error; err != nil {
		return video, err
	}
	return video, nil
}

func (vs *VideoService) Info(id, userId uint) (map[string]interface{}, error) {
	var videoInfo VideoInfo
	if err := db.Table("video_Video as v").
		Select("v.id,v.title,v.duration,v.poster,v.size,v.width,v.height,v.codec_name,v.channel_layout,v.creation_time,l.collect, l.browse, l.zan, l.cai, l.watch, group_concat(a.id,',') as actress_ids, group_concat(a.actress,',') as actress_names, group_concat(a.avatar,',') as actress_avatars").
		Joins("left join video_VideoLog as l on l.video_id = v.id").
		Joins("left join video_VideoActress as va on va.video_id = v.id").
		Joins("left join video_Actress as a on a.id = va.actress_id").
		Where("v.id = ?", id).
		Group("v.id,v.title,v.duration,v.poster,v.size,v.width,v.height,v.codec_name,v.channel_layout,v.creation_time,l.collect, l.browse, l.zan, l.cai, l.watch").
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
	var isCollect bool = false
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
		"collectNum":    videoInfo.CollectNum,
		"browseNum":     videoInfo.BrowseNum,
		"likeNum":       videoInfo.LikeNum,
		"dislikeNum":    videoInfo.DisLikeNum,
		"watchNum":      videoInfo.WatchNum,
		"collectID":     collectID,
		"isCollect":     isCollect,
	}
	return data, nil
}

func (vs *VideoService) List() ([]model.Video, error) {
	var videos []model.Video
	if err := db.Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func (vs *VideoService) Create(videos []model.Video) error {
	if err := db.Create(&videos).Error; err != nil {
		return err
	}
	return nil
}

func (vs *VideoService) Collect(videoID uint, collect int, userID uint) error {
	var video model.Video
	if errors.Is(db.First(&video, videoID).Error, gorm.ErrRecordNotFound) {
		return errors.New("视频不存在！")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var expr string
		if collect == 1 {
			// 增加1
			expr = "collect + 1"
			if err := tx.Create(&model.UserCollectLog{UserID: userID, VideoID: videoID}).Error; err != nil {
				return errors.New("创建失败！")
			}
		} else {
			// 减少1
			expr = "collect - 1"
			if err := tx.Where("user_id = ? and video_id = ?", userID, videoID).Delete(&model.UserCollectLog{}).Error; err != nil {
				return errors.New("删除失败！")
			}
		}
		if err := tx.Model(&model.VideoLog{}).Where("video_id = ?", videoID).Update("collect", gorm.Expr(expr)).Error; err != nil {
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
		var userBrowseLog model.UserBrowseLog
		if err := tx.Where(model.UserBrowseLog{UserID: userID, VideoID: videoID}).FirstOrInit(&userBrowseLog).Error; err != nil {
			return err
		}

		if err := tx.Where(model.UserBrowseLog{UserID: userID, VideoID: videoID}).Assign(model.UserBrowseLog{Number: userBrowseLog.Number + 1}).FirstOrCreate(&model.UserBrowseLog{}).Error; err != nil {
			return fmt.Errorf("创建失败: %s", err)
		}

		var videoLog model.VideoLog
		if err := tx.Where(model.VideoLog{VideoID: videoID}).FirstOrInit(&videoLog).Error; err != nil {
			return err
		}

		if err := tx.Where(model.VideoLog{VideoID: videoID}).Assign(model.VideoLog{Browse: videoLog.Browse + 1}).FirstOrCreate(&model.VideoLog{}).Error; err != nil {
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

	rdb.HSet(ctx, utils.Join("video_video_", strconv.FormatUint(uint64(id), 10)), "id", id, "title", title, "poster", data["poster"], "duration", data["duration"], "browseNum", data["browseNum"], "browseNum", data["browseNum"])
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
			CollectNum:    data["collectNum"].(uint),
			BrowseNum:     data["browseNum"].(uint),
			LikeNum:       data["likeNum"].(uint),
			DisLikeNum:    data["dislikeNum"].(uint),
			WatchNum:      data["watchNum"].(uint),
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

func (vs *VideoService) Comment(videoID uint, content string, userID uint) (uint, error) {
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	comment := model.VideoComment{
		ParentId:    0,
		VideoId:     videoID,
		UserId:      userID,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      "APPROVED",
		IsAnonymous: 1,
		Content:     content,
		IsShow:      1,
	}

	result := db.Create(&comment)
	if result.Error != nil {
		return 0, result.Error
	}

	return comment.ID, nil
}

func (vs *VideoService) Reply(videoID uint, parentID uint, content string, userID uint) (uint, error) {
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return 0, err
	}

	comment := model.VideoComment{
		ParentId:    parentID,
		VideoId:     videoID,
		UserId:      userID,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Status:      "APPROVED",
		IsAnonymous: 1,
		Content:     content,
		IsShow:      1,
	}

	tx := db.Begin()

	if err := tx.Create(&comment).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Model(&model.VideoComment{}).Where("id = ?", parentID).Update("reply_num", gorm.Expr("reply_num + 1")).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return comment.ID, nil
}

type VideoComment struct {
	model.VideoComment
	LogUserID uint `gorm:"comment:用户ID"`
	Zan       uint `gorm:"comment:支持（赞）"`
	Cai       uint `gorm:"comment:反对（踩）"`
}

func (vs *VideoService) CommentList(videoID, userID uint) ([]*CommentTree, error) {
	var list []VideoComment
	query := db.Model(&model.UserCommentLog{}).Where("video_id = ? and user_id = ?", videoID, userID)
	if err := db.Table("video_VideoComment as c").
		Where("c.video_id = ?", videoID).
		Select("c.*", "l.user_id as log_user_id", "l.support as zan", "l.oppose as cai").
		Joins("left join (?) l on l.comment_id = c.id", query).Order("c.id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return tree(list), nil
}

func (vs *VideoService) Zan(commentID, userID uint, zan int) error {
	var comment model.VideoComment
	if errors.Is(db.First(&comment, commentID).Error, gorm.ErrRecordNotFound) {
		return errors.New("评论不存在！")
	}

	tx := db.Begin()

	var expr string
	var support uint
	if zan == 1 {
		// 增加1
		support = 1
		expr = "support + 1"
	} else {
		// 减少1
		support = 0
		expr = "support - 1"
	}

	if err := tx.Where(model.UserCommentLog{UserID: userID, VideoID: comment.VideoId, CommentID: commentID}).Assign(model.UserCommentLog{Support: &support}).FirstOrCreate(&model.UserCommentLog{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建失败: %s", err)
	}
	if err := tx.Model(&model.VideoComment{}).Where("id = ? and video_id = ?", commentID, comment.VideoId).Update("support", gorm.Expr(expr)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新失败: %s", err)
	}

	tx.Commit()

	return nil
}

func (vs *VideoService) Cai(commentID, userID uint, cai int) error {
	var comment model.VideoComment
	if errors.Is(db.First(&comment, commentID).Error, gorm.ErrRecordNotFound) {
		return errors.New("评论不存在！")
	}

	tx := db.Begin()

	var expr string
	var oppose uint
	if cai == 1 {
		// 增加1
		oppose = 1
		expr = "oppose + 1"
	} else {
		// 减少1
		oppose = 0
		expr = "oppose - 1"
	}

	if err := tx.Where(model.UserCommentLog{UserID: userID, VideoID: comment.VideoId, CommentID: commentID}).Assign(model.UserCommentLog{Oppose: &oppose}).FirstOrCreate(&model.UserCommentLog{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建失败: %s", err)
	}
	if err := tx.Model(&model.VideoComment{}).Where("id = ? and video_id = ?", commentID, comment.VideoId).Update("oppose", gorm.Expr(expr)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("更新失败: %s", err)
	}

	tx.Commit()

	return nil
}

type CommentTree struct {
	VideoComment
	// model.VideoComment
	Childrens []CommentTree
}

func tree(list []VideoComment) []*CommentTree {
	// func tree(list []model.VideoComment) []*CommentTree {
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
	// return recursive(data, childrens)
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

// func recursive(data map[uint]*CommentTree, childrens map[uint][]CommentTree) map[uint]*CommentTree {
// 	for _, v := range data {
// 		videoComments, ok := childrens[v.ID]
// 		if ok {
// 			v.Childrens = videoComments
// 			delete(childrens, v.ID)
// 			if len(childrens) > 0 {
// 				data := make(map[uint]*CommentTree, len(videoComments))
// 				for _, v := range videoComments {
// 					videoComment := v
// 					data[v.ID] = &videoComment
// 				}
// 				recursive(data, childrens)
// 			}
// 		}
// 	}
// 	return data
// }

type VideoDanmu struct {
	Text   string  `json:"text" gorm:"comment:弹幕文本"`
	Time   float64 `json:"time" gorm:"comment:弹幕时间, 默认为当前播放器时间"`
	Mode   uint8   `json:"mode" gorm:"comment:弹幕模式: 0: 滚动(默认) 1: 顶部 2: 底部"`
	Color  string  `json:"color" gorm:"comment:弹幕颜色，默认为白色"`
	Border bool    `json:"border" gorm:"comment:弹幕是否有描边, 默认为 false"`
	// Style  string `json:"style" gorm:"comment:弹幕自定义样式, 默认为空对象"`
}

func (vs *VideoService) DanmuList(videoID uint) ([]VideoDanmu, error) {
	var list []VideoDanmu
	if err := db.Model(&model.VideoDanmu{}).Where("video_id = ?", videoID).Order("time asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (vs *VideoService) DanmuSave(videoID, userID uint, text string, time float64, mode uint8, color string, border bool, style string) error {
	var danmu = model.VideoDanmu{
		VideoId: videoID,
		UserId:  userID,
		Text:    text,
		Time:    time,
		Mode:    mode,
		Color:   color,
		Border:  border,
		Style:   style,
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
