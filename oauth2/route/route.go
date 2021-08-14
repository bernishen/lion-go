package route

import "github.com/bernishen/lion-go/utils/router"

var (
	config *router.Config
	Router *router.LionRouter
)

func init() {
	config = &router.Config{
		Address:     "192.168.56.101",
		Port:        "6500",
		RoutePrefix: "/api/oauth2",
	}
	Router = router.InitRouter(config)
}
