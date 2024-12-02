package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
	"gorm.io/gorm"
)

type VideoService struct{}

type VideoInfo struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title" gorm:"column:title;type:varchar(255);comment:标题"`
	Actress        string    `json:"actress" gorm:"column:actress;type:varchar(100);comment:演员"`
	Size           int64     `json:"size" gorm:"column:size;type:bigint;comment:大小"`
	Duration       float64   `json:"duration" gorm:"column:duration;type:float;default:0;comment:时长"`
	Poster         string    `json:"poster" gorm:"column:poster;type:varchar(255);comment:封面"`
	Width          int       `json:"width" gorm:"column:width;type:int;default:0;comment:宽"`
	Height         int       `json:"height" gorm:"column:height;type:int;default:0;comment:高"`
	CodecName      string    `json:"codec_name" gorm:"column:codec_name;type:varchar(90);comment:编解码器"`
	ChannelLayout  string    `json:"channel_layout" gorm:"column:channel_layout;type:varchar(90);comment:音频声道"`
	CreationTime   time.Time `gorm:"column:creation_time;type:date;comment:时间"`
	Collect        uint      `json:"collect" gorm:"column:collect;type:uint;not null;default:0;comment:收藏"`
	Browse         uint      `json:"browse" gorm:"column:browse;type:uint;not null;default:0;comment:浏览"`
	Zan            uint      `json:"zan" gorm:"column:zan;type:uint;not null;default:0;comment:赞"`
	Cai            uint      `json:"cai" gorm:"column:cai;type:uint;not null;default:0;comment:踩"`
	Watch          uint      `json:"watch" gorm:"column:watch;type:uint;not null;default:0;comment:观看"`
	ActressIds     string    `json:"actress_ids" gorm:"column:actress_ids;type:varchar(255);comment:演员ID"`
	ActressNames   string    `json:"actress_names" gorm:"column:actress_names;type:varchar(255);comment:演员名称"`
	ActressAvatars string    `json:"actress_avatars" gorm:"column:actress_avatars;type:varchar(2550);comment:演员头像"`
}
type Video struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Poster   string `json:"poster"`
	Duration string `json:"duration"`
	Browse   uint   `json:"browse"`
	Watch    uint   `json:"watch"`
	//Actress       string  `json:"actress"`
	//Size          float64 `json:"size"`
	//ModTime       string  `json:"mod_time"`
	//Width         int     `json:"width"`
	//Height        int     `json:"height"`
	//CodecName     string  `json:"codec_name"`
	//ChannelLayout string  `json:"channel_layout"`
	//Collect       uint    `json:"collect"`
	//Zan           uint    `json:"zan"`
	//Cai           uint    `json:"cai"`
}

func (vs *VideoService) Find(actressID int, page, pageSize int, action, sort string) (data map[string]interface{}, err error) {
	var ids []uint

	f := func(ids []uint, count int) (map[string]interface{}, error) {
		videos := make([]Video, len(ids))
		for i, id := range ids {
			result := rdb.HGetAll(ctx, utils.Join("video_video_", strconv.Itoa(int(id)))).Val()
			browse, _ := strconv.Atoi(result["browse"])
			watch, _ := strconv.Atoi(result["watch"])
			videos[i] = Video{
				ID:       id,
				Title:    result["title"],
				Poster:   result["poster"],
				Duration: result["duration"],
				Browse:   uint(browse),
				Watch:    uint(watch),
			}
		}
		data = map[string]interface{}{
			"list":  videos,
			"count": count,
		}
		return data, nil
	}

	if actressID != 0 {
		var vadb = db.Model(&model.VideoActress{}).Where("actress_id = ?", actressID)
		var count int64
		if err = vadb.Count(&count).Error; err != nil {
			return nil, err
		}
		if err = vadb.Scopes(Paginate(page, pageSize, int(count))).Pluck("video_id", &ids).Error; err != nil {
			return nil, err
		}
		return f(ids, int(count))
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
	var count int64
	if err = vdb.Count(&count).Error; err != nil {
		return nil, err
	}
	if err = db.Table("(?)", vdb).Scopes(Paginate(page, pageSize, int(count))).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	bts, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}
	result, _ := rdb.HGet(ctx, key, "ids").Result()
	if strings.Compare(string(bts), result) == 0 && result != "" {
		return f(ids, int(count))
	}

	rows, err := vdb.Scopes(Paginate(page, pageSize, int(count))).Rows()
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

		//f, _ := strconv.ParseFloat(strconv.FormatInt(videoInfo.Size, 10), 64)
		video := Video{
			ID:       videoInfo.ID,
			Title:    videoInfo.Title,
			Poster:   videoInfo.Poster,
			Duration: utils.ResolveTime(uint32(videoInfo.Duration)),
			Browse:   videoInfo.Browse,
			Watch:    videoInfo.Watch,
			//Actress:       videoInfo.Actress,
			//Size:          f / 1024 / 1024,
			//ModTime:       videoInfo.CreationTime.Format("2006-01-02 15:04:05"),
			//Width:         videoInfo.Width,
			//Height:        videoInfo.Height,
			//CodecName:     videoInfo.CodecName,
			//ChannelLayout: videoInfo.ChannelLayout,
			//Collect:       videoInfo.Collect,
			//Zan:           videoInfo.Zan,
			//Cai:           videoInfo.Cai,
		}
		videos = append(videos, video)
		keys = append(keys, utils.Join("video_video_", strconv.Itoa(int(videoInfo.ID))))
	}

	txf := func(tx *redis.Tx) error {
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.HSet(ctx, key, "len", len(ids), "ids", string(bts))
			for _, a := range videos {
				pipe.HSet(ctx, utils.Join("video_video_", strconv.Itoa(int(a.ID))), "id", a.ID, "title", a.Title, "poster", a.Poster, "duration", a.Duration, "browse", a.Browse, "watch", a.Watch)
			}
			return nil
		})
		return err
	}
	if err = rdb.Watch(ctx, txf, keys...); errors.Is(err, redis.TxFailedErr) {
		return nil, err
	}

	//pages := pagination.NewPaginator(int(count), pageSize)
	//pages.SetCurrentPage(page)
	//data := map[string]interface{}{
	//	"list": videos,
	//	"page": map[string]interface{}{
	//		"totalPage":   pages.TotalPage(),
	//		"prePage":     pages.PrePage(),
	//		"currentPage": pages.CurrentPage(),
	//		"nextPage":    pages.NextPage(),
	//		"pages":       pages.Pages(),
	//	},
	//}

	data = map[string]interface{}{
		"list":  videos,
		"count": count,
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

func (vs *VideoService) Info(id uint) (VideoInfo, error) {
	var videoInfo VideoInfo
	if err := db.Table("video_Video as v").
		Select("v.id,v.title,v.duration,v.poster,v.size,v.width,v.height,v.codec_name,v.channel_layout,v.creation_time,l.collect, l.browse, l.zan, l.cai, l.watch, group_concat(a.id,',') as actress_ids, group_concat(a.actress,',') as actress_names, group_concat(a.avatar,',') as actress_avatars").
		Joins("left join video_VideoLog as l on l.video_id = v.id").
		Joins("left join video_VideoActress as va on va.video_id = v.id").
		Joins("left join video_Actress as a on a.id = va.actress_id").
		Where("v.id = ?", id).
		Group("v.id,v.title,v.duration,v.poster,v.size,v.width,v.height,v.codec_name,v.channel_layout,v.creation_time,l.collect, l.browse, l.zan, l.cai, l.watch").
		Scan(&videoInfo).Error; err != nil {
		return VideoInfo{}, err
	}
	return videoInfo, nil
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

	tx := db.Begin()

	var expr string
	if collect == 1 {
		// 增加1
		expr = "collect + 1"

		if err := tx.Create(&model.UserCollectLog{UserID: userID, VideoID: videoID}).Error; err != nil {
			tx.Rollback()
			return errors.New("创建失败！")
		}
	} else {
		// 减少1
		expr = "collect - 1"

		if err := tx.Where("user_id = ? and video_id = ?", userID, videoID).Delete(&model.UserCollectLog{}).Error; err != nil {
			tx.Rollback()
			return errors.New("删除失败！")
		}
	}
	if err := tx.Model(&model.VideoLog{}).Where("video_id = ?", videoID).Update("collect", gorm.Expr(expr)).Error; err != nil {
		tx.Rollback()
		return errors.New("更新失败！")
	}

	tx.Commit()

	return nil
}

func (vs *VideoService) Browse(videoID uint, userID uint) error {
	var video model.Video
	if errors.Is(db.First(&video, videoID).Error, gorm.ErrRecordNotFound) {
		return errors.New("视频不存在！")
	}

	tx := db.Begin()

	var userBrowseLog model.UserBrowseLog
	if err := tx.Where(model.UserBrowseLog{UserID: userID, VideoID: videoID}).FirstOrInit(&userBrowseLog).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(model.UserBrowseLog{UserID: userID, VideoID: videoID}).Assign(model.UserBrowseLog{Number: userBrowseLog.Number + 1}).FirstOrCreate(&model.UserBrowseLog{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建失败: %s", err)
	}
	var videoLog model.VideoLog
	if err := tx.Where(model.VideoLog{VideoID: videoID}).FirstOrInit(&videoLog).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(model.VideoLog{VideoID: videoID}).Assign(model.VideoLog{Browse: videoLog.Browse + 1}).FirstOrCreate(&model.VideoLog{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("创建失败: %s", err)
	}

	tx.Commit()

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
	LogUserID uint `gorm:"column:log_user_id;type:uint;not null;default:0;comment:用户ID"`
	Zan       uint `gorm:"column:zan;type:uint;not null;default:0;comment:支持（赞）"`
	Cai       uint `gorm:"column:cai;type:uint;not null;default:0;comment:反对（踩）"`
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
	Text   string  `json:"text" gorm:"column:text;type:text;not null;comment:弹幕文本"`
	Time   float64 `json:"time" gorm:"column:time;type:double;not null;comment:弹幕时间, 默认为当前播放器时间"`
	Mode   uint8   `json:"mode" gorm:"column:mode;type:uint;not null;default:0;comment:弹幕模式: 0: 滚动(默认)，1: 顶部，2: 底部"`
	Color  string  `json:"color" gorm:"column:color;type:text;not null;comment:弹幕颜色，默认为白色"`
	Border bool    `json:"border" gorm:"column:border;type:bool;not null;default:false;comment:弹幕是否有描边, 默认为 false"`
	// Style  string `json:"style" gorm:"column:style;type:text;not null;comment:弹幕自定义样式, 默认为空对象"`
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
