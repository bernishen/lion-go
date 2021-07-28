package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var wsUpgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WSCallRoute is
type WSCallRoute func(funcID string, param *WSMessage) *WSReturn

type WSService struct {
	heartbeat int
	callRoute WSCallRoute
	connPool  map[string]*wsConnection
}

type WSMessage struct {
	MessageType int
	Data        []byte
}

type WSReturn struct {
	HasReturn bool
	Data      *WSMessage
}

type wsConnection struct {
	id                string
	connection        *websocket.Conn
	rQueue            chan *WSMessage
	wQueue            chan *WSMessage
	pQueue            chan *WSMessage
	heartbeatInterval time.Duration
	heartbeatNum      int
	callRoute         WSCallRoute
	isClosed          bool
	closeState        chan byte
	service           *WSService
	mutex             sync.Mutex
}

func InitWebSocket(heartbeatInterval int, callRoute WSCallRoute) *WSService {
	service := &WSService{
		heartbeat: heartbeatInterval,
		callRoute: callRoute,
		connPool:  make(map[string]*wsConnection),
	}
	return service
}

func Return(msg *WSMessage) *WSReturn {
	return &WSReturn{
		true,
		msg,
	}
}

func NoReturn() *WSReturn {
	return &WSReturn{
		false,
		nil,
	}
}

func (service *WSService) Handler(w http.ResponseWriter, r *http.Request, wsID string) {
	//   完成握手 升级为 WebSocket长连接，使用conn发送和接收消息。
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	wsconn := &wsConnection{
		id:      wsID,
		service: service,

		connection:        conn,
		heartbeatInterval: time.Duration(service.heartbeat),
		rQueue:            make(chan *WSMessage, 100),
		wQueue:            make(chan *WSMessage, 100),
		pQueue:            make(chan *WSMessage, 100),
	}

	service.connPool[wsconn.id] = wsconn
	go wsconn.read()
	go wsconn.write()
	go wsconn.heartbeat()
	go wsconn.funcHandler()

	//调用连接的WriteMessage和ReadMessage方法以一片字节发送和接收消息。实现如何回显消息：
	//p是一个[]字节，messageType是一个值为websocket.BinaryMessage或websocket.TextMessage的int。
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		log.Printf("Read from client msg:%s \n", msg)

		if err := conn.WriteMessage(messageType, msg); err != nil {
			//if err := conn.WriteMessage(1, []byte("今天。。。"));err != nil {
			log.Println("Writeing error...", err)
			return
		}
		log.Printf("Write msg to client: recved: %s \n", msg)
	}
}

func (service *WSService) WriteMsg(wsID string, messageType int, data []byte) {
	wsConn, ok := service.connPool[wsID]
	if !ok {
		log.Println("")
		return
	}
	msg := &WSMessage{
		MessageType: messageType,
		Data:        data,
	}
	select {
	case wsConn.wQueue <- msg:
	case <-wsConn.closeState:
		return
	}
}

func (service *WSService) CloseSocket(wsID string) {
	wsConn, ok := service.connPool[wsID]
	if !ok {
		log.Println("")
		return
	}
	wsConn.close()
}

func (wsConn *wsConnection) heartbeat() {
	for {
		select {
		case p := <-wsConn.pQueue:
			wsConn.heartbeatNum = 2
			if p.MessageType != websocket.PingMessage {
				continue
			}
			wsConn.connection.WriteMessage(websocket.PongMessage, []byte("Pong"))
		case <-time.After(time.Second * wsConn.heartbeatInterval * 2):
			if wsConn.heartbeatNum <= 0 {
				wsConn.close()
			}
			wsConn.heartbeatNum--
			wsConn.connection.WriteMessage(websocket.PingMessage, []byte("Ping"))
			continue
		}
	}
}

func (wsConn *wsConnection) read() {
	for {
		msgType, readDatas, err := wsConn.connection.ReadMessage()
		if err != nil {
			var errMsg = []byte("Reading error.|消息读取异常。--System:-->" + err.Error())
			wsConn.connection.WriteMessage(websocket.TextMessage, errMsg)
			continue
		}
		msg := &WSMessage{
			MessageType: msgType,
			Data:        readDatas,
		}
		switch msgType {
		case websocket.PingMessage:
			select {
			case wsConn.pQueue <- msg:
			case <-wsConn.closeState:
				return
			}
		case websocket.PongMessage:
			select {
			case wsConn.pQueue <- msg:
			case <-wsConn.closeState:
				return
			}
		default:
			select {
			case wsConn.rQueue <- msg:
			case <-wsConn.closeState:
				return
			}
		}
	}
}

func (wsConn *wsConnection) write() {
	for {
		select {
		case msg := <-wsConn.wQueue:
			if err := wsConn.connection.WriteMessage(msg.MessageType, msg.Data); err != nil {
				log.Println("发送消息给客户端发生错误", err.Error())
				if wsConn.isClosed {
					return
				}
			}
		case <-wsConn.closeState:
			return
		}
	}
}

func (wsConn *wsConnection) funcHandler() {
	for {
		select {
		case msg := <-wsConn.rQueue:
			ret := wsConn.callRoute(wsConn.id, msg)
			if ret == nil || ret.HasReturn == false {
				continue
			}
			select {
			case wsConn.wQueue <- ret.Data:
			case <-wsConn.closeState:
				return
			}
		case <-wsConn.closeState:
			return
		}
	}
}

func (wsConn *wsConnection) close() {
	log.Println("关闭连接被调用了")
	wsConn.connection.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if wsConn.isClosed == false {
		wsConn.isClosed = true
		// 删除这个连接的变量
		delete(wsConn.service.connPool, wsConn.id)
		close(wsConn.closeState)
	}
}
