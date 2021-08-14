package security

import (
	"github.com/bernishen/lion-go/oauth2/route"
	"github.com/bernishen/lion-go/oauth2/service/rsaservice"
	"github.com/bernishen/lion-go/oauth2/service/securityservice"
	"github.com/bernishen/exception"
	"github.com/bernishen/lion-go/utils/router"
	"strings"
)

func init() {
	s := router.InitController("/security").
		Get(publicKey)
	route.Router.Register(s)
	s0 := router.InitController("/securitytest").
		Get(privateKey)
	route.Router.Register(s0)
}

func publicKey(ctx *router.Context) (*interface{}, *exception.Exception) {
	cid := strings.Trim(ctx.Request.Header.Get("cid"), " ")
	if cid == "" {
		return nil, exception.NewException(exception.Error, 1001, "The clientID(cid) is null.")
	}

	v := strings.Trim(ctx.Request.Header.Get("version"), " ")
	if v == "" {
		v = "lastest"
	}

	pub, ex := securityservice.CreateKey(cid, v)
	if ex != nil {
		return nil, ex
	}

	var ret interface{} = pub
	return &ret, nil
}

func privateKey(ctx *router.Context) (*interface{}, *exception.Exception) {
	cid := strings.Trim(ctx.Request.Header.Get("cid"), " ")
	if cid == "" {
		return nil, exception.NewException(exception.Error, 1001, "The clientID(cid) is null.")
	}
	pvt, ex := rsaservice.FindKeyByClient(cid)
	if ex != nil {
		return nil, ex
	}
	var ret interface{} = pvt
	return &ret, nil
}
