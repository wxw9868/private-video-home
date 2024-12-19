package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	metrics "github.com/bmc-toolbox/gin-go-metrics"
	metricsMiddleware "github.com/bmc-toolbox/gin-go-metrics/middleware"
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
	// Optional part to send metrics to Graphite,
	// as alternative you can send metrics from
	// rcrowley/go-metrics.DefaultRegistry yourself
	err := metrics.Setup(
		"graphite",  // clientType
		"localhost", // graphite host
		2003,        // graphite port
		"server",    // metrics prefix
		time.Minute, // graphite flushInterval
	)
	if err != nil {
		fmt.Printf("Failed to set up monitoring: %s\n", err)
		os.Exit(1)
	}

	// gin.SetMode(gin.ReleaseMode)

	// 强制日志颜色化
	gin.ForceConsoleColor()
	gin.DefaultWriter = colorable.NewColorableStdout()

	r := gin.Default()

	// argument to NewMetrics tells which variables need to be
	// expanded in metrics, more on that by link:
	// https://banzaicloud.com/blog/monitoring-gin-with-prometheus/
	p := metricsMiddleware.NewMetrics([]string{"expanded_parameter"})
	r.Use(p.HandlerFunc(
		nil,
		[]string{"/ping", "/api/ping"}, // ignore given URLs from stats
		true,                           // replace "/" with "_" in URLs to prevent splitting Graphite namespace
	))

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

	//err := new(service.VideoService).ImportVideoData("D:/ta", "浅川さおり", "向井瞳", "小島さとみ", "竹内柚葉", "森本ひとみ")
	//err := new(service.ActressService).SaveActress()
	//err := new(service.ActressService).DownAvatar()
	//err := new(service.ActressService).CopyAvatar()
	//fmt.Println(err)

	if err := r.Run(fmt.Sprintf("%s:%d", conf.Host, conf.Port)); err != nil {
		panic(err)
	}
}
