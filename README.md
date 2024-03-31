# my video

## 先决条件
### 1. 安装 FFmpeg
根据自己的系统选择下载并安装：
> [ffmpeg，ffprobe二进制文件](https://ffbinaries.com/downloads)
```sh
go get github.com/u2takey/ffmpeg-go
```
### 2. 安装 SQLite
```sh
# 通过下面的命令查看 sqlite3 是否已安装，一般情况下系统都默认安装了 SQLite 数据库 
sqlite3 -version
```

## 实时重载
### 安装air
```sh
go install github.com/cosmtrek/air@latest
```
### 使用方法
您可以添加 alias air='~/.air' 到您的 .bashrc 或 .zshrc 后缀的文件.

首先，进入你的项目文件夹
```sh
cd /path/to/your_project
```
最简单的方法是执行
```sh
# 优先在当前路径查找 `.air.toml` 后缀的文件，如果没有找到，则使用默认的
air -c .air.toml
```
您可以运行以下命令初始化，把默认配置添加到当前路径下的.air.toml 文件。
```sh
air init
```
在这之后，你只需执行 air 命令，无需添加额外的变量，它就能使用 .air.toml 文件中的配置了。
```sh
air
```
如欲修改配置信息，请参考 [air_example.toml](https://github.com/cosmtrek/air/blob/master/air_example.toml) 文件.

## 部署到安卓上教程
```sh
adb version

SET GOOS=android
SET GOARCH=arm64

$env:GOOS="android"
$env:GOARCH="arm64"

go build -o myvideo main.go

adb push myvideo /data/local/tmp
chmod 755 myvideo
```

## 使用的开源库和工具
* [ffmpeg-go](https://github.com/u2takey/ffmpeg-go)
* [air](https://github.com/cosmtrek/air/blob/master/README-zh_cn.md)

+ [西瓜视频播放器](https://h5player.bytedance.com/)
+ [ckplayer](https://www.ckplayer.com/)

- [Bootstrap V5](https://v5.bootcss.com/)
- [jQuery API 3.5.1 速查表](https://jquery.cuishifeng.cn/index.html)

+ [免费Favicon.ico图标在线生成器](https://www.logosc.cn/logo/favicon)
+ [一个工具箱 - 好用的在线工具都在这里！](http://www.atoolbox.net/)

* bootstrap-5.3.0