#!/bin/bash

# 指定要检查的镜像ID或名称
IMAGE_NAME="redis:latest"

# 使用docker image inspect获取镜像信息
INSPECT_OUTPUT=$(docker image inspect "$IMAGE_NAME")

# 检查docker命令是否成功执行
if [ $? -eq 0 ]; then
    # 使用jq提取并打印镜像的ID和创建时间
    echo "镜像ID: $(echo "$INSPECT_OUTPUT" | jq -r '.[0].Id')"
    echo "创建时间: $(echo "$INSPECT_OUTPUT" | jq -r '.[0].Created')"
else
    echo "无法获取镜像信息，请检查镜像名称或Docker服务状态"
fi


docker image inspect redis:latest

docker pull redis:latest


docker run -d -p 8080:8080 --name my-video-v1.0.1 --mount type=bind,source=E:/video/assets,target=/usr/src/app/assets --mount type=bind,source=E:/video/database/sqlite,target=/usr/src/app/database/sqlite my-video:v1.0.1