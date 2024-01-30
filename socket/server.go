package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var connSlice []*websocket.Conn

func ServerConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("conn错误", err)
		return
	}
	connSlice = append(connSlice, conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("conn.ReadMessage 读取信息错误", err)
			break
		}

		for i, _ := range connSlice {
			if conn != connSlice[i] {
				connSlice[i].WriteMessage(websocket.BinaryMessage, p)
			}
		}
		fmt.Println(messageType, string(p))
	}

	defer conn.Close()
	log.Println("服务关闭")
}
