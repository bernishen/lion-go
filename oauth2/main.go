package main

import (
	"fmt"
	_ "github.com/bernishen/lion-go/oauth2/controller/security"
	_ "github.com/bernishen/lion-go/oauth2/controller/user"
	"github.com/bernishen/lion-go/oauth2/route"
	"github.com/bernishen/lion-go/utils/router"
)

func main() {
	defer func(){ // 必须要先声明defer，否则不能捕获到panic异常
		fmt.Println("c")
		if err:=recover();err!=nil{
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
		fmt.Println("d")
	}()
	fmt.Println("****** lionGO-OAuth2 Starting ******")
	router.ListenDefault(route.Router)
}
