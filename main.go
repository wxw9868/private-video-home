package main

import (
	"fmt"
	"log"

	"github.com/wxw9868/video/initialize/db"
	"github.com/wxw9868/video/router"
	"github.com/wxw9868/video/utils"
)

func main() {
	db.RegisterTables()

	router := router.Engine()

	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Fatal(err)
	}

	if err := router.Run(fmt.Sprintf("%s:%d", ip, 8080)); err != nil {
		panic(err)
	}
}
