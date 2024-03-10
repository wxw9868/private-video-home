package main

import (
	"fmt"
	"log"

	"github.com/wxw9868/video/api"
)

func main() {
	router := api.Engine()

	ip, err := api.GetLocalIP()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server Running")
	if err := router.Run(fmt.Sprintf("%s:%d", ip, 80)); err != nil {
		panic(err)
	}
}
