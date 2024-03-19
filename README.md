# video

## 先决条件
### 安装ffmpeg
根据自己的系统选择下载并安装：[ffmpeg二进制文件下载地址](https://ffbinaries.com/downloads)
```shell
go get github.com/u2takey/ffmpeg-go
```

## 使用的开源库和工具
* [ckplayer.com](https://www.ckplayer.com/)
* [Bootstrap V5](https://v5.bootcss.com/)
* [ffmpeg-go](https://github.com/u2takey/ffmpeg-go)

- [jQuery API 3.5.1 速查表](https://jquery.cuishifeng.cn/index.html)
- [免费Favicon.ico图标在线生成器](https://www.logosc.cn/logo/favicon)
- [一个工具箱 - 好用的在线工具都在这里！](http://www.atoolbox.net/)

```shell
adb version

SET GOOS=android
SET GOARCH=arm64

$env:GOOS="android"
$env:GOARCH="arm64"

go build -o myvideo main.go

adb push myvideo /data/local/tmp
chmod 755 myvideo
```