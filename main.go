package main

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wxw9868/video/config"
	"github.com/wxw9868/video/docs"
	"github.com/wxw9868/video/middleware"
	"github.com/wxw9868/video/router"
	"net/http"
)

// @title Video API
// @version 1.0
// @description This is a video server.

// @host 127.0.0.1:6000
// @host 192.168.0.9:6000
// @BasePath /
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

	// Set a lower memory limit for multipart forms (default is 32 MiB)
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

	//actresss := []string{"五月あおい"}
	//err := service.VideoImport("D:/ta", actresss)
	//fmt.Println(err)
	//actresss := "目々澤めぐ,瀬戸レイカ"
	//service.RepairVideoImport(strings.Split(actresss, ","))

	if err := r.Run(fmt.Sprintf("%s:%d", conf.Host, conf.Port)); err != nil {
		panic(err)
	}
}
