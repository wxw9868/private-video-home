package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/wxw9868/video/config"
	"github.com/wxw9868/video/middleware"
	"github.com/wxw9868/video/router"
)

func main() {
	var ids []uint
	fmt.Println(len(ids))
	actresss := make([]string, len(ids))
	fmt.Println(actresss)
	return

	gin.SetMode(gin.ReleaseMode)

	// 强制日志颜色化
	gin.ForceConsoleColor()
	gin.DefaultWriter = colorable.NewColorableStdout()

	r := gin.Default()

	//r := gin.New()
	//logger, _ := zap.NewProduction()
	//r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	//r.Use(ginzap.RecoveryWithZap(logger, true))
	r.Use(middleware.GinCors()) // 允许跨域
	r.NoRoute(middleware.NoRoute())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/assets", "./assets")

	pprof.Register(r) // 性能监控
	router.Router(r)

	conf := config.Config().System
	if err := r.Run(fmt.Sprintf("%s:%d", conf.Host, conf.Port)); err != nil {
		panic(err)
	}
}
