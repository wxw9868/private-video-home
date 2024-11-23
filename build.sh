docker build -t my-video:v1.0.0 .

docker run -d -p 8080:8080 --name my-video \
  --mount type=bind,source=E:/video/assets,target=/go/src/app/assets \
  --mount type=bind,source=E:/video/database/sqlite,target=/go/src/app/database/sqlite \
  my-video:v1

docker run -d -p 8080:8080 --name my-video --mount type=bind,source=D:/video/assets,target=/go/src/app/assets --mount type=bind,source=D:/video/database,target=/go/src/app/database my-video:v1.0.0