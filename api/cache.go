package api

//func CacheVideoActress(c *gin.Context) {
//	if len(actressListSort) > 0 {
//		actresss = make([]actress, len(actressListSort))
//		for index, name := range actressListSort {
//			nameSlice := []rune(name)
//			avatarPath := avatarDir + "/" + name + ".png"
//
//			_, err := os.Stat(avatarPath)
//			if os.IsNotExist(err) {
//				utils.GenerateAvatar(string(nameSlice[0]), avatarPath)
//			}
//
//			actresss[index] = actress{
//				ID:      index + 1,
//				Actress: name,
//				Avatar:  avatarPath,
//			}
//		}
//	}
//
//	actressBytes, _ := json.Marshal(actresss)
//
//	c.HTML(http.StatusOK, "actress.html", gin.H{
//		"title":       "演员列表",
//		"actressList": string(actressBytes),
//	})
//}

//func CacheVideoList(c *gin.Context) {
//	if c.Query("actress_id") != "" {
//		indexInt, _ := strconv.Atoi(c.Query("actress_id"))
//		actress := actressListSort[indexInt]
//		index := actressList[actress]
//		actressVideos := make([]video, len(index))
//		for k, v := range index {
//			actressVideos[k] = videos[v]
//		}
//
//		videosBytes, _ := json.Marshal(actressVideos)
//
//		c.HTML(http.StatusOK, "list.html", gin.H{
//			"title":       "视频列表",
//			"data":        string(videosBytes),
//			"actressList": actressListSort,
//			"actressID":   indexInt,
//		})
//		return
//	}
//
//	if videos == nil && list == nil {
//		writeToCache()
//	}
//
//	videosBytes, _ := json.Marshal(videos)
//
//	c.HTML(http.StatusOK, "list.html", gin.H{
//		"title":       "视频列表",
//		"data":        string(videosBytes),
//		"actressList": actressListSort,
//		"actressID":   -1,
//	})
//}

//func CacheVideoPlay(c *gin.Context) {
//	id := c.Query("id")
//	intId, _ := strconv.Atoi(id)
//	name := list[intId]
//
//	videoUrl := videoDir + "/" + name
//
//	c.HTML(http.StatusOK, "play.html", gin.H{
//		"title":     "视频播放",
//		"videoUrl":  videoUrl,
//		"videoName": name,
//	})
//}

//func writeToCache() error {
//	files, err := os.ReadDir(videoDir)
//	if err != nil {
//		return err
//	}
//
//	list = make([]string, len(files))
//	videos = make([]video, len(files))
//
//	for k, file := range files {
//		filename := file.Name()
//		ext := filepath.Ext(filename)
//		if ext == ".mp4" {
//			strs := strings.Split(filename, ".")
//			title := strs[0]
//			arrs := strings.Split(strs[0], "_")
//			actress := arrs[len(arrs)-1]
//			if _, ok := actressList[actress]; !ok {
//				actressListSort = append(actressListSort, actress)
//			}
//			actressList[actress] = append(actressList[actress], k)
//
//			fi, _ := file.Info()
//			size, _ := strconv.ParseFloat(strconv.FormatInt(fi.Size(), 10), 64)
//
//			filePath := videoDir + "/" + filename
//			fil, _ := os.Open(filePath)
//			duration, _ := utils.GetMP4Duration(fil)
//
//			videoInfo, _ := utils.VideoInfo(filePath)
//
//			posterPath := posterDir + "/" + title + ".jpg"
//			_, err = os.Stat(posterPath)
//			if os.IsNotExist(err) {
//				_ = utils.ReadFrameAsJpeg(filePath, posterPath, "00:02:30")
//			}
//
//			//snapshotPath := snapshotDir + "/" + title + ".gif"
//			//_ = CutVideoForGif(filePath, posterPath)
//
//			list[k] = filename
//			videos[k] = video{
//				ID:            k + 1,
//				Title:         title,
//				Actress:       actress,
//				Size:          size / 1024 / 1024,
//				Duration:      utils.ResolveTime(duration),
//				ModTime:       fi.ModTime().Format("2006-01-02 15:04:05"),
//				Poster:        posterPath,
//				Width:         int(videoInfo["width"].(int64)),
//				Height:        int(videoInfo["height"].(int64)),
//				CodecName:     fmt.Sprintf("%s,%s", videoInfo["codec_name0"].(string), videoInfo["codec_name1"].(string)),
//				ChannelLayout: videoInfo["channel_layout"].(string),
//			}
//		}
//	}
//	return nil
//}
