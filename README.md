# My Video

<!-- vscode-markdown-toc -->
* 1. [先决条件](#)
	* 1.1. [1. 安装 FFmpeg](#FFmpeg)
	* 1.2. [2. 安装 SQLite](#SQLite)
	* 1.3. [2. 安装和启动 GoFound](#GoFound)
* 2. [实时重载](#-1)
	* 2.1. [安装 air](#air)
	* 2.2. [使用方法](#-1)
* 3. [部署到安卓上教程](#-1)
* 4. [使用的开源库和工具](#-1)
* 5. [数据库操作命令](#-1)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

##  1. <a name=''></a>先决条件
###  1.1. <a name='FFmpeg'></a>1. 安装 FFmpeg
根据自己的系统选择下载并安装：
> 下载 [ffmpeg，ffprobe二进制文件](https://ffbinaries.com/downloads)
```sh
go get github.com/u2takey/ffmpeg-go
```
###  1.2. <a name='SQLite'></a>2. 安装 SQLite
```sh
# 通过下面的命令查看 sqlite3 是否已安装，一般情况下系统都默认安装了 SQLite 数据库 
sqlite3 -version
```
###  1.3. <a name='GoFound'></a>2. 安装和启动 GoFound
> 下载好源码之后，进入到源码目录，执行下列两个命令
* 编译
> 直接下载 [可执行文件](https://github.com/newpanjing/gofound/releases) 可以不用编译，省去这一步。
```sh
git clone https://github.com/sea-team/gofound.git
cd gofound
go get && go build 
```
* 命令启动
```sh
./gofound --addr=:8080 --data=./data
```
* Docker命令部署
```sh
docker build -t gofound .
docker run -d --name gofound -p 5678:5678 -v /mnt/data/gofound:/usr/local/go_found/data gofound:latest
```
###  1.4. <a name='Redis'></a>1. 安装和启动 Redis
```sh
# 拉取官方的最新版本的镜像
docker pull redis:latest
# 运行 redis 容器
docker run -d -p 6379:6379 --name my-video-redis -v E:/video/database/redis:/data -v E:/video/config/redis.conf:/etc/redis/redis.conf redis redis-server /etc/redis/redis.conf redis-server --appendonly yes
# 通过 redis-cli 连接测试使用 redis 服务
$ docker exec -it my-video-redis /bin/bash
```
## 安装和启动 Video
```sh
docker build -t my-video:v1 .

docker run -d -p 8080:8080 --name my-video \
--mount type=bind,source=E:/video/assets,target=/usr/src/app/assets \
--mount type=bind,source=E:/video/database/sqlite,target=/usr/src/app/database/sqlite my-video:v1

docker stop <容器id>
```
##  2. <a name='-1'></a>实时重载
###  2.1. <a name='air'></a>安装 air
```sh
go install github.com/cosmtrek/air@latest
```
###  2.2. <a name='-1'></a>使用方法
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

##  3. <a name='-1'></a>部署到安卓上教程
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

##  4. <a name='-1'></a>使用的开源库和工具
* [zap](https://github.com/uber-go/zap.git)
* [ffmpeg-go](https://github.com/u2takey/ffmpeg-go)
* [air](https://github.com/cosmtrek/air/blob/master/README-zh_cn.md)

##  5. <a name='-1'></a>数据库操作命令
```sql
-- 删除数据
DELETE FROM video_UserCommentLog;
-- 重置主键
UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'video_UserCommentLog';
```

[go-gin简单实现tls https服务](https://www.cnblogs.com/davis12/p/16918591.html)