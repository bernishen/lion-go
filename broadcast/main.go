package main

import (
	"fmt"
	"github.com/bernishen/lion-go/utils/router"
)

func main() {
	fmt.Println("****** lionGO-Broadcast Starting ******")
	c := router.Config{
		Address: "192.168.56.101",
		Port:    "6501",
	}
	router.ListenDefault(c)
}
