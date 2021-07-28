package main

import (
	"fmt"
	_ "github.com/Berni-Shen/lion-go/oauth2/controller/user"
	"github.com/Berni-Shen/lion-go/utils/router"
)

func main() {
	fmt.Println("****** lionGO-OAuth2 Starting ******")
	c := router.Config{
		Address: "192.168.56.103",
		Port:    "6500",
	}
	router.ListenDefault(c)
}
