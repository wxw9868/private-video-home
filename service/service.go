package service

import (
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/wxw9868/video/model"
	"github.com/wxw9868/video/utils"
)

type VideoService struct {
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

func do() {
	files, err := os.ReadDir(videoDir)
	if err != nil {
		log.Fatal(err)
	}

	list = make([]string, len(files))
	videos = make([]model.Video, len(files))

	for k, file := range files {
		filename := file.Name()
		ext := filepath.Ext(filename)
		if ext == ".mp4" {
			strs := strings.Split(filename, ".")
			title := strs[0]
			arrs := strings.Split(strs[0], "_")
			actress := arrs[len(arrs)-1]
			if _, ok := actressList[actress]; !ok {
				actressListSort = append(actressListSort, actress)
			}
			actressList[actress] = append(actressList[actress], k)

			fi, _ := file.Info()

			filePath := videoDir + "/" + filename
			vi, _ := utils.VideoInfo(filePath)

			posterPath := posterDir + "/" + title + ".jpg"
			_, err = os.Stat(posterPath)
			if os.IsNotExist(err) {
				_ = utils.ReadFrameAsJpeg(filePath, posterPath, "00:01:45")
			}

			//snapshotPath := snapshotDir + "/" + title + ".gif"
			//_ = CutVideoForGif(filePath, posterPath)

			list[k] = filename
			videos[k] = model.Video{
				Title:    title,
				Actress:  actress,
				Size:     *big.NewInt(fi.Size()),
				Duration: 0,
				ModTime:  fi.ModTime(),
				Poster:   posterPath,
			}
		}
	}

	actresss := make([]model.Actress, len(actressListSort))

	if len(actressListSort) > 0 {
		for index, name := range actressListSort {
			nameSlice := []rune(name)
			avatarPath := avatarDir + "/" + name + ".png"

			_, err := os.Stat(avatarPath)
			if os.IsNotExist(err) {
				if err := utils.GenerateAvatar(string(nameSlice[0]), avatarPath); err != nil {
					log.Fatal(err)
				}
			}

			actresss[index] = model.Actress{
				Actress: name,
				Avatar:  avatarPath,
			}
		}
	}

}
