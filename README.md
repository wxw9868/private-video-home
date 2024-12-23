# My Video

## 目录
- [My Video](#my-video)
  - [目录](#目录)
  - [1. 先决条件](#1-先决条件)
    - [1.1. 安装FFmpeg](#11-安装ffmpeg)
    - [1.2. 安装SQLite](#12-安装sqlite)
    - [1.3. 安装和启动GoFound](#13-安装和启动gofound)
    - [1.4. 安装和启动Redis](#14-安装和启动redis)
  - [2. 启动服务](#2-启动服务)
  - [3. 第三方库](#3-第三方库)
  - [4. 数据库操作命令](#4-数据库操作命令)
    - [待开发功能](#待开发功能)
  
## 1. 先决条件
### 1.1. 安装FFmpeg
根据自己的系统选择下载并安装：[下载ffmpeg，ffprobe二进制文件](https://ffbinaries.com/downloads)
```sh
go get github.com/u2takey/ffmpeg-go
```
### 1.2. 安装SQLite
```sh
# 通过下面的命令查看 sqlite3 是否已安装，一般情况下系统都默认安装了 SQLite 数据库 
sqlite3 -version
```

### 1.3. 安装和启动GoFound
```sh
git clone https://github.com/sea-team/gofound.git

cd gofound
```
>下载好源码之后，进入到源码目录，执行下列两个命令
* 编译
>直接下载 [可执行文件](https://github.com/newpanjing/gofound/releases) 可以不用编译，省去这一步。
```sh
go get && go build 
```
* 启动
```sh
./gofound --addr=:5678 --data=./data --auth=admin:123456
```
* docker部署
```sh
docker build -t gofound .

docker run -d --name gofound -p 5678:5678 -v D:/database/gofound/data:/usr/local/go_found/data gofound:latest
```
* 其他命令 参考 [配置文档](https://github.com/sea-team/gofound/blob/main/docs/config.md)

### 1.4. 安装和启动Redis
```sh
# 拉取官方的最新版本的镜像
docker pull redis:latest

# 运行 redis 容器
docker run -d -p 6379:6379 --name my-video-redis -v D:/database/redis/data:/data -v D:/database/redis/redis.conf:/etc/redis/redis.conf redis redis-server /etc/redis/redis.conf

# 通过 redis-cli 连接测试使用 redis 服务
$ docker exec -it my-video-redis /bin/bash
```
* 通过配置文件启动
```sh
redis-server /Users/v_weixiongwei/go/src/video/config/redis.conf
```

## 2. 启动服务
* air启动 [air使用教程](air.md)
```sh
air
```
* 命令启动
```sh
go run main.go
```
* docker部署
```sh
docker build -t my-video:v1 .

docker run -d -p 8080:8080 --name my-video \
--mount type=bind,source=E:/video/assets,target=/usr/src/app/assets \
--mount type=bind,source=E:/video/database/sqlite,target=/usr/src/app/database/sqlite my-video:v1
```

## 3. 第三方库
* [zap](https://github.com/uber-go/zap.git)
* [ffmpeg-go](https://github.com/u2takey/ffmpeg-go)
* [air](https://github.com/cosmtrek/air/blob/master/README-zh_cn.md)
* [gin-swagger](https://github.com/swaggo/gin-swagger)
* [go-gin简单实现tls https服务](https://www.cnblogs.com/davis12/p/16918591.html)

## 4. 数据库操作命令
```sql
-- 删除数据
DELETE FROM video_UserCommentLog;
-- 重置主键
UPDATE SQLITE_SEQUENCE SET seq = 0 WHERE name = 'video_UserCommentLog';
```

[gofound](http://192.168.0.9:5678/admin)

### 待开发功能
1. 使用Redis数据库记录页面浏览量日志数据