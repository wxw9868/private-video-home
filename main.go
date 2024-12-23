package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
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
	// 禁用控制台颜色，将日志写入文件时不需要控制台颜色。
	// gin.DisableConsoleColor()

	// 记录到文件。
	// f, _ := os.Create("gin.log")
	// gin.DefaultWriter = io.MultiWriter(f)

	// 如果需要同时将日志写入文件和控制台，请使用以下代码。
	// gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	r := gin.Default()
	//logger, _ := zap.NewProduction()
	//r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	//r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(middleware.GinCors())
	r.NoRoute(middleware.NoRoute())

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("assets", "assets")
	r.GET("ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) })

	pprof.Register(r) // 性能监控
	router.Router(r)

	addr := fmt.Sprintf("%s:%d", config.Config().System.Host, config.Config().System.Port)
	ginSwaggerURL := fmt.Sprintf("http://%s/swagger/doc.json", addr)

	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler, ginSwagger.URL(ginSwaggerURL), ginSwagger.DefaultModelsExpandDepth(-1)))

	//err := new(service.VideoService).ImportVideoData("D:/ta/video", "知念真紀", "秋元若菜")
	//err := new(service.ActressService).SaveActress()
	//err := new(service.ActressService).DownAvatar()
	//err := new(service.ActressService).CopyAvatar()
	//fmt.Println(err)

	//if err := gracehttp.Serve(&http.Server{Addr: addr, Handler: r}); err != nil {
	//	panic(err)
	//}

	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
