package main

import (
	"fmt"
	"log"
)

func main() {
	router := Engine()

	ip, err := getLocalIP()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server Running")
	if err := router.Run(fmt.Sprintf("%s:%d", ip, 80)); err != nil {
		panic(err)
	}
}
