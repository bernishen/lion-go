package system

import (
	"encoding/json"
	"github.com/bernishen/lion-go/utils/websocket"
	ws "github.com/gorilla/websocket"
	"log"
)

type funcRun struct {
	funcID string
	name   string
	data   string
}

type retObj struct {
	funcID  string
	code    int
	message string
	data    interface{}
}

func initFunc(msg *websocket.WSMessage) *funcRun {
	var f funcRun
	err := json.Unmarshal(msg.Data, f)
	if err != nil {
		log.Printf("This is unvalid message of 'FuncRun'.")
		return nil
	}
	return &f
}

func returnOK(retData *[]byte) *websocket.WSMessage {
	dataObj := &retObj{
		code:    200,
		message: "successful",
		data:    retData,
	}
	data, _ := json.Marshal(dataObj)
	return &websocket.WSMessage{
		MessageType: ws.TextMessage,
		Data:        data,
	}
}

func returnFailed(code int, message string) *websocket.WSMessage {
	dataObj := &retObj{
		code:    code,
		message: message,
		data:    nil,
	}
	data, _ := json.Marshal(dataObj)
	return &websocket.WSMessage{
		MessageType: ws.TextMessage,
		Data:        data,
	}
}
