package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"image/color/palette"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/shiningrush/avatarbuilder"
	"github.com/shiningrush/avatarbuilder/calc"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func Format() {
	// "2006-01-02 15:04:05"
}

// GenerateAvatar 生成头衔
func GenerateAvatar(name, path string) error {
	palette := palette.Plan9
	var bgColor color.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	var frontColor color.Color = color.RGBA{0x00, 0x00, 0x00, 0xff}
	bgColor = palette[randint(len(palette))]
	frontColor = palette[randint(len(palette))]
	ab := avatarbuilder.NewAvatarBuilder("D:/video/assets/ttf/SourceHanSansSC-Medium.ttf", &calc.SourceHansSansSCMedium{})
	ab.SetBackgroundColor(bgColor)
	ab.SetFrontgroundColor(frontColor)
	ab.SetFontSize(300)
	ab.SetAvatarSize(460, 460)
	if err := ab.GenerateImageAndSave(name, path); err != nil {
		return err
	}
	return nil
}

var (
	rng   = rand.New(rand.NewSource(time.Now().UnixNano()))
	rngMu = new(sync.Mutex)
)

func randint(n int) int {
	rngMu.Lock()
	defer rngMu.Unlock()
	return rng.Intn(n)
}

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("unable to determine local IP address")
}

// BoxHeader 信息头
type BoxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

// GetMP4Duration 获取视频时长，以秒计
func GetMP4Duration(reader io.ReaderAt) (lengthOfTime uint32, err error) {
	var info = make([]byte, 0x10)
	var boxHeader BoxHeader
	var offset int64 = 0
	// 获取moov结构偏移
	for {
		_, err = reader.ReadAt(info, offset)
		if err != nil {
			return
		}
		boxHeader = getHeaderBoxInfo(info)
		fourccType := getFourccType(boxHeader)
		if fourccType == "moov" {
			break
		}
		// 有一部分mp4 mdat尺寸过大需要特殊处理
		if fourccType == "mdat" {
			if boxHeader.Size == 1 {
				offset += int64(boxHeader.Size64)
				continue
			}
		}
		offset += int64(boxHeader.Size)
	}
	// 获取moov结构开头一部分
	moovStartBytes := make([]byte, 0x100)
	_, err = reader.ReadAt(moovStartBytes, offset)
	if err != nil {
		return
	}
	// 定义timeScale与Duration偏移
	timeScaleOffset := 0x1C
	durationOffest := 0x20
	timeScale := binary.BigEndian.Uint32(moovStartBytes[timeScaleOffset : timeScaleOffset+4])
	duration := binary.BigEndian.Uint32(moovStartBytes[durationOffest : durationOffest+4])
	lengthOfTime = duration / timeScale
	return
}

// getHeaderBoxInfo 获取头信息
func getHeaderBoxInfo(data []byte) (boxHeader BoxHeader) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &boxHeader)
	return
}

// getFourccType 获取信息头类型
func getFourccType(boxHeader BoxHeader) (fourccType string) {
	fourccType = string(boxHeader.FourccType[:])
	return
}

// ResolveTime 将秒转成时分秒格式
func ResolveTime(seconds uint32) string {
	var (
		h, m, s string
	)
	var day = seconds / (24 * 3600)
	hour := (seconds - day*3600*24) / 3600
	minute := (seconds - day*24*3600 - hour*3600) / 60
	second := seconds - day*24*3600 - hour*3600 - minute*60
	h = strconv.Itoa(int(hour))
	if hour < 10 {
		h = "0" + strconv.Itoa(int(hour))
	}
	m = strconv.Itoa(int(minute))
	if minute < 10 {
		m = "0" + strconv.Itoa(int(minute))
	}
	s = strconv.Itoa(int(second))
	if second < 10 {
		s = "0" + strconv.Itoa(int(second))
	}
	return fmt.Sprintf("%s:%s:%s", h, m, s)
}

// CutVideoForGif 将视频剪切为 GIF
// ss 开始时间  例子：00:00:15
// t  持续时间  例子：00:00:06
// r 设定帧速率，默认为25
// s 设定画面的宽与高
func CutVideoForGif(videoPath, snapshotPath, ss string) error {
	err := ffmpeg.Input(videoPath, ffmpeg.KwArgs{"ss": ss}).
		Output(snapshotPath, ffmpeg.KwArgs{"pix_fmt": "rgb24", "t": "00:00:15", "r": "30"}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}
	return nil
}

// ReadFrameAsJpeg 将视频剪切为 JPG
func ReadFrameAsJpeg(inFileName, outFileName, ss string) error {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName, ffmpeg.KwArgs{"ss": ss}).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", 0)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg", "qscale:v": 2}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		return err
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		return err
	}

	if err = imaging.Save(img, outFileName); err != nil {
		return err
	}

	return nil
}

func VideoInfo(inFileName string) (map[string]interface{}, error) {
	a, err := ffmpeg.Probe(inFileName)
	if err != nil {
		return nil, err
	}

	isCreationTime := gjson.Get(a, "format.tags.creation_time").IsBool()
	var creationTime time.Time
	if isCreationTime {
		creationTime = gjson.Get(a, "format.tags.creation_time").Time().Add(0)
	} else {
		fi, _ := os.Stat(inFileName)
		creationTime = fi.ModTime()
	}
	duration := gjson.Get(a, "format.duration").Float()
	size := gjson.Get(a, "format.size").Int()
	width := gjson.Get(a, "streams.0.width").Int()
	height := gjson.Get(a, "streams.0.height").Int()
	codecName0 := gjson.Get(a, "streams.0.codec_name").String()
	codecName1 := gjson.Get(a, "streams.1.codec_name").String()
	channelLayout := gjson.Get(a, "streams.1.channel_layout").String()

	// fmt.Println(a)
	// fmt.Printf("大小：%d\n", size)
	// fmt.Printf("尺寸：%d x %d\n", width, height)
	// fmt.Printf("编解码器：%s %s\n", strings.ToUpper(codecName1), strings.ToUpper(codecName0))
	// fmt.Printf("时长：%f\n", duration)
	// fmt.Printf("音频声道：%s\n", channelLayout)
	// fmt.Printf("创建时间：%s\n", creationTime.Format("2006:01:02 15:04:06"))

	data := map[string]interface{}{
		"duration":       duration,
		"size":           size,
		"creation_time":  creationTime,
		"width":          width,
		"height":         height,
		"codec_name0":    codecName0,
		"codec_name1":    codecName1,
		"channel_layout": channelLayout,
	}

	return data, nil
}

// Join 字符串拼接
func Join(s ...string) string {
	var b bytes.Buffer
	for i := 0; i < len(s); i++ {
		b.WriteString(s[i])
	}
	return b.String()
}

// 读取文件数据到 map
func ReadFileToMap(name string, v any) error {
	bytes, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytes, &v); err != nil {
		return err
	}
	return nil
}

// 读取 map 数据到文件
func WriteMapToFile(name string, v any) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = os.WriteFile(name, bytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

func VideoFileRename(nameMap map[string]string, nameSlice, actressSlice []string) error {
	filename := ""
	for index, oldName := range nameMap {
		filename = oldName
		for _, v := range nameSlice {
			filename = strings.Replace(filename, v, "", -1)
		}
		if filename[6:7] == "-" {
			filename = strings.Replace(filename, filename[6:7], "_", -1)
		}
		if !strings.Contains(filename, filename[0:10]+"_") && strings.Contains(filename, filename[0:10]) {
			filename = strings.Replace(filename, filename[0:10], filename[0:10]+"_", -1)
		}
		for _, v := range actressSlice {
			if strings.Contains(filename, v) {
				if v == "Vol." {
					filename = strings.Replace(filename, v, "Vol_", -1)
				} else {
					filename = strings.Replace(filename, v, "_"+v, -1)
				}
			}
		}
		fmt.Println(filename)
		nameMap[index] = filename
	}
	return nil
}

func VideoRename(videoDir string, nameMap map[string]string, nameSlice, actressSlice []string) error {
	files, err := os.ReadDir(videoDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no file")
	}

	for _, file := range files {
		filename := file.Name()
		oldPath := videoDir + "/" + filename

		oldName := strings.Split(filename, ".")[0]
		newName, ok := nameMap[oldName]
		if ok {
			filename = strings.Replace(filename, oldName, newName, -1)
		} else {
			for _, v := range nameSlice {
				filename = strings.Replace(filename, v, "", -1)
			}
			if filename[6:7] == "-" {
				filename = strings.Replace(filename, filename[6:7], "_", -1)
			}
			if !strings.Contains(filename, filename[0:10]+"_") && strings.Contains(filename, filename[0:10]) {
				filename = strings.Replace(filename, filename[0:10], filename[0:10]+"_", -1)
			}
			for _, v := range actressSlice {
				if strings.Contains(filename, v) {
					if v == "Vol." {
						filename = strings.Replace(filename, v, "Vol_", -1)
					} else {
						if !strings.Contains(filename, "_"+v) {
							filename = strings.Replace(filename, v, "_"+v, -1)
						}
					}
				}
			}
			fmt.Println(filename)
		}

		newPath := videoDir + "/" + filename
		if err = os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

func GeneteSQL() string {
	var data = make(map[string]struct{})
	ReadFileToMap("data.json", &data)
	fmt.Println(len(data))
	return ""
	var avatarDir = "./assets/image/avatar"
	var actressSql = "INSERT OR REPLACE INTO video_Actress (actress, avatar, CreatedAt, UpdatedAt) VALUES "
	var root = "/Users/v_weixiongwei/go/src/video"
	root = "D:/GoLang/video"
	for actress, _ := range data {
		avatarPath := avatarDir + "/" + actress + ".png"
		rootPath := root + "/assets/image/avatar" + "/" + actress + ".png"
		_, err := os.Stat(rootPath)
		if os.IsNotExist(err) {
			nameSlice := []rune(actress)
			if err := GenerateAvatar(string(nameSlice[0]), rootPath); err != nil {
				log.Fatal(err)
			}
		}
		actressSql += fmt.Sprintf("('%s', '%s', '%v', '%v'), ", actress, avatarPath, time.Now().Local(), time.Now().Local())
	}
	actressSqlBytes := []byte(actressSql)
	actressSql = string(actressSqlBytes[:len(actressSqlBytes)-2])

	return actressSql
}
