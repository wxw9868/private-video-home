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

	dev := false
	dev = true

	ip, err := utils.GetLocalIP()
	if err != nil {
		log.Fatal(err)
	}
	port := 80
	if dev {
		ip = "127.0.0.1"
		port = 8081
	}
	ip = "192.168.0.9"
	ip = "0.0.0.0"
	port = 8080

	addr := fmt.Sprintf("%s:%d", ip, port)
	router := router.Engine(addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
	// if err := router.RunTLS(addr, "cert/server.pem", "cert/server.key"); err != nil {
	// 	log.Fatal(err)
	// }
}
