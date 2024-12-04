FROM golang:1.22

ENV GOPROXY=https://goproxy.cn

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . ./
RUN go build -v -o /videosys

EXPOSE 8080

CMD ["/videosys"]