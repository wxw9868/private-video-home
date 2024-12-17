package main

import (
	"fmt"
	"net/http"

	"github.com/wxw9868/video/service"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wxw9868/video/config"
	"github.com/wxw9868/video/docs"
	"github.com/wxw9868/video/middleware"
	"github.com/wxw9868/video/router"
)

//	@title			私人视频 API
//	@version		1.0
//	@description	This is a private video server.

//	@contact.name	API Support
//	@contact.email	wxw9868gmail.com

//	@host		192.168.0.9:8080
//	@BasePath	/

func main() {
	// gin.SetMode(gin.ReleaseMode)

	// 强制日志颜色化
	gin.ForceConsoleColor()
	gin.DefaultWriter = colorable.NewColorableStdout()

	r := gin.Default()

	//logger, _ := zap.NewProduction()
	//r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	//r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(middleware.GinCors()) // 允许跨域
	r.NoRoute(middleware.NoRoute())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/assets", "./assets")

	pprof.Register(r) // 性能监控
	router.Router(r)

	conf := config.Config().System
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", conf.Host, conf.Port)),
		ginSwagger.DefaultModelsExpandDepth(-1)),
	)

	//err := new(service.VideoService).ImportVideoData("D:/ta", "美咲愛", "笠木日向", "中田みなみ", "当麻叶美")
	//err := new(service.ActressService).SaveActress()
	//err := new(service.ActressService).DownAvatar()
	err := new(service.ActressService).CopyAvatar()
	fmt.Println(err)

	if err := r.Run(fmt.Sprintf("%s:%d", conf.Host, conf.Port)); err != nil {
		panic(err)
	}
}
