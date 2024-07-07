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

	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Fatal(err)
	}

	ip = "0.0.0.0"
	port := 8080

	addr := fmt.Sprintf("%s:%d", ip, port)
	router := router.Engine(addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
	// if err := router.RunTLS(addr, "cert/server.pem", "cert/server.key"); err != nil {
	// 	log.Fatal(err)
	// }
}
