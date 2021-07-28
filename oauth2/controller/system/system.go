package system

import (
	"encoding/json"
	"github.com/Berni-Shen/lion-go/oauth2/controller/system/domain"
	"github.com/Berni-Shen/lion-go/oauth2/service/sessionservice"
	"github.com/Berni-Shen/lion-go/utils/router"
	"github.com/Berni-Shen/lion-go/utils/websocket"
	ws "github.com/gorilla/websocket"
)

var wsService *websocket.WSService

func init() {
	wsService = websocket.InitWebSocket(20, callRoute)
	s := router.InitController("/system/{wsID}").
		Get(wsService.Handler, "wsID")
	router.Default().Register(s)
}

func callRoute(wsID string, msg *websocket.WSMessage) *websocket.WSReturn {
	switch msg.MessageType {
	case ws.TextMessage:
		f := initFunc(msg)
		if f == nil {
			return websocket.Return(returnFailed(301, "The parameter can not be parsed."))
		}
		switch f.name {
		case "Verify":
			return verify(wsID, f.data)
		case "SignOut":
			return signOut(f.data)
		default:
			return websocket.Return(returnFailed(301, "Found not function name of "+f.name+"."))
		}
	default:
		return websocket.Return(returnOK(nil))
	}
}

func verify(wsID string, param string) *websocket.WSReturn {
	pByte := []byte(param)
	var c domain.Client
	err := json.Unmarshal(pByte, c)
	if err != nil {
		return websocket.Return(returnFailed(301, err.Error()))
	}

	s, ex := sessionservice.VerifyGlobal(c.AccessToken)
	if ex != nil {
		return websocket.Return(returnFailed(301, err.Error()))
	}

	r, ok := s.Roles[c.SystemID]
	if !ok {
		return websocket.Return(returnFailed(301, "This user does not have this system permission."))
	}

	_, ex = sessionservice.VerifySystem(c.SystemID, c.AccessToken)
	if ex != nil {

	}

	sToken, ex := sessionservice.NewSystem(c.SystemID, c.AccessToken, &r)
	if ex != nil {
		return websocket.Return(returnFailed(301, ex.Message))
	}

	ret := []byte(sToken)
	return websocket.Return(returnOK(&ret))
}

func signOut(param string) *websocket.WSReturn {
	pByte := []byte(param)
	var c domain.Client
	err := json.Unmarshal(pByte, c)
	if err != nil {
		return websocket.Return(returnFailed(301, err.Error()))
	}

	//sessionservice.SignOutAll(c.AccessToken)
	return websocket.NoReturn()
}
