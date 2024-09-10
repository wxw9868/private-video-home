package main

import (
	"fmt"
	"net/http"

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

// @title Score Admin API
// @version 1.0
// @description This is a score home server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 0.0.0.0:8080
// @BasePath /
func main() {
	// gin.SetMode(gin.ReleaseMode)

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
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler,
		ginSwagger.URL(fmt.Sprintf("http://%s:%d/swagger/doc.json", "127.0.0.1", conf.Port)),
		ginSwagger.DefaultModelsExpandDepth(-1)),
	)

	if err := r.Run(fmt.Sprintf("%s:%d", conf.Host, conf.Port)); err != nil {
		panic(err)
	}
}
