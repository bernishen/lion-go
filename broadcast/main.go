package main

import (
	"fmt"
	"github.com/Berni-Shen/lion-go/utils/router"
)

func main() {
	fmt.Println("****** lionGO-Broadcast Starting ******")
	c := router.Config{
		Address: "192.168.56.103",
		Port:    "6501",
	}
	router.ListenDefault(c)
}
