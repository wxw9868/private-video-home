docker build -t my-video:v1 .

docker run -d -p 8080:8080 --name my-video --mount type=bind,source=E:/video/assets,target=/usr/src/app/assets --mount type=bind,source=E:/video/database,target=/usr/src/app/database my-video:v1

docker stop 39fc7722df43