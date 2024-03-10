package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Engine() *gin.Engine {
	router := gin.Default()

	// 允许跨域
	//router.Use(cors.Default())
	router.NoRoute(NoRoute())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 性能监控
	//pprof.Register(router)

	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("template/*")

	router.GET("/", videoList)
	router.GET("/play", videoPlay)
	router.GET("/home", func(c *gin.Context) {
		c.HTML(http.StatusOK, "home.html", nil)
	})
	router.GET("/rename", videoRename)

	return router
}

func NoRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
		c.Abort()
	}
}
