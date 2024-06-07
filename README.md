
# My Video

## 先决条件
### 1. 安装 FFmpeg
根据自己的系统选择下载并安装：
> 下载 [ffmpeg，ffprobe二进制文件](https://ffbinaries.com/downloads)
```sh
go get github.com/u2takey/ffmpeg-go
```
### 2. 安装 SQLite
```sh
# 通过下面的命令查看 sqlite3 是否已安装，一般情况下系统都默认安装了 SQLite 数据库 
sqlite3 -version
```
### 2. 安装和启动 GoFound
> 下载好源码之后，进入到源码目录，执行下列两个命令
* 编译
> 直接下载 [可执行文件](https://github.com/newpanjing/gofound/releases) 可以不用编译，省去这一步。
```sh
go get && go build 
```
* 启动
```sh
./gofound --addr=:8080 --data=./data
```
* docker部署
```sh
docker build -t gofound .
docker run -d --name gofound -p 5678:5678 -v /mnt/data/gofound:/usr/local/go_found/data gofound:latest
```

## 实时重载
### 安装 air
```sh
go install github.com/cosmtrek/air@latest
```
### 使用方法
> 您可以添加 alias air='~/.air' 到您的 .bashrc 或 .zshrc 后缀的文件.
```sh
# 首先，进入你的项目文件夹
cd /path/to/your_project

# 最简单的方法是执行
# 优先在当前路径查找 `.air.toml` 后缀的文件，如果没有找到，则使用默认的
air -c .air.toml

# 您可以运行以下命令初始化，把默认配置添加到当前路径下的.air.toml 文件。
air init

# 在这之后，你只需执行 air 命令，无需添加额外的变量，它就能使用 .air.toml 文件中的配置了。
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

* [bootstrap-5.3.0](https://v5.bootcss.com/)

## 数据库操作命令
```sql
-- 删除数据
DELETE FROM video_UserCommentLog;
-- 重置主键
UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'video_UserCommentLog';
```