package socket

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

func ClientConnect() {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://127.0.0.1:8080", nil)
	if err != nil {
		log.Println("客户端连接失败", err)
		return
	}

	go write(conn)

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取信息失败", err)
			break
		}
		fmt.Println(string(p))
	}

	defer conn.Close()
}
func write(conn *websocket.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println("客户端写入失败", err)
		}
		conn.WriteMessage(websocket.BinaryMessage, line)
	}
}
