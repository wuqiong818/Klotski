package socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"klotski/pojo"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub 保证这里的hub变量是唯一的，并且所有的Client对象操作的是同一个
var Hub = make(map[string]*pojo.Client)

func BuildConn(w http.ResponseWriter, r *http.Request) {
	var client pojo.Client
	client.Hub = Hub
	//1. 服务端与客户端建立连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("conn错误", err)
		return
	}
	fmt.Println("远程主机连接成功 IP为", conn.RemoteAddr())
	//2. 在信息中枢处根据消息类型进行特定的处理
	client.UserConn = conn

	for {
		_, p, err := client.UserConn.ReadMessage()
		if err != nil {
			log.Println("server.go conn.ReadMessage 读取信息错误", err)
			break
		}

		var requestPkg pojo.RequestPkg
		err = json.Unmarshal(p, &requestPkg)
		if err != nil {
			fmt.Println("websocket反序列化失败", err)
			return
		}
		switch requestPkg.Type {
		case pojo.CertificationType: //用户认证
			//1.读取数据
			var certificationInfo pojo.CertificationInfo
			err = json.Unmarshal([]byte(requestPkg.Data), &certificationInfo)
			//fmt.Println("certificationInfo", certificationInfo)
			//fmt.Println("UserId", certificationInfo.UserId)
			//2.进行用户验证
			isCer, _ := client.Certificate(&certificationInfo)
			//3.返回响应
			responsePkg := new(pojo.ResponsePkg)
			if isCer {
				responsePkg.Type = pojo.CertificationType
				responsePkg.Code = pojo.Success
			} else {
				responsePkg.Type = pojo.CertificationType
				responsePkg.Code = pojo.UserCertificationError
				responsePkg.Data = "用户认证失败,非法侵入，服务端主动关闭连接"
			}
			err := client.UserConn.WriteJSON(responsePkg)
			if err != nil {
				fmt.Println("服务端用户认证信息返回错误", err)
			}
			if !isCer {
				err := client.UserConn.Close()
				if err != nil {
					fmt.Println("用户认证失败，服务端主动关闭连接错误", err)
				}
			}
		case pojo.CreateRoomType: //创建房间号
			responsePkg := new(pojo.ResponsePkg)
			responsePkg.Type = pojo.CreateRoomType
			if client.UserCer {
				client.CreateRoom()
				responsePkg.Code = pojo.Success
				responsePkg.Data = client.RoomId

			} else {
				responsePkg.Code = pojo.UserCertificationError
				responsePkg.Data = "用户认证失败,请先完成用户认证"
			}
			err := client.UserConn.WriteJSON(responsePkg)
			if err != nil {
				fmt.Println("服务端返回房间信息error", err)
			}
		case pojo.JoinRoomType:
			fmt.Println("JoinRoomType")

		case pojo.RefreshScoreType:
			fmt.Println("RefreshScoreType")

		case pojo.DiscontinueQuitType:
			fmt.Println("DiscontinueQuitType")

		case pojo.GameOverType:
			fmt.Println("GameOverType")

		}
	}
	defer conn.Close()
	log.Println("服务关闭")
}
