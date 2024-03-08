package socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"klotski/pkg"
	"klotski/pojo"
	"log"
	"net/http"
	"time"
)

//// Hub 保证这里的hub变量是唯一的，并且所有的Client对象操作的是同一个
//var Hub = make(map[string]map[string]*pojo.Client)

func WsHandle(hub *pojo.HupCenter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//将短连接 升级成 长连接-建立WebSocket通信
		upgrade := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrade.Upgrade(w, r, nil)
		if err != nil {
			log.Println("conn错误", err)
			return
		}

		fmt.Println("远程主机连接成功 IP为", conn.RemoteAddr())
		//进行client的初始化操作
		client := &pojo.Client{User: &pojo.User{}, Hub: &pojo.HupCenter{}} //非字段不要为nil
		client.Hub = hub
		client.User.UserConn = conn
		client.User.HealthCheck = time.Now().Add(time.Duration(pkg.HeartCheckSecond) * time.Second) //健康时间

		//计时器 : 如果用户规定秒内没有完成用户认证，则直接断开连接
		time.AfterFunc(time.Duration(pkg.UserAuthSecond)*time.Second, func() {
			if !client.User.UserCer {
				fmt.Println("用户认证失败，关闭连接")
				client.User.Close()
			}
		})

		//接受信息 根据信息类型进行分别处理
		Controller(client)
	}
}

func Controller(client *pojo.Client) {
	defer func() {
		client.Hub.UnRegister <- client
	}()

	for {
		_, p, err := client.User.UserConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				//用户主动给关闭连接后的输出
				log.Println("WebSocket closed")
			} else {
				//服务器主动断开走这个，从一个断开的连接中读取信息
				log.Println("server.go conn.ReadMessage 读取信息错误", err)
			}
			return
		}

		var requestPkg pojo.RequestPkg
		err = json.Unmarshal(p, &requestPkg)
		if err != nil {
			fmt.Println("websocket反序列化失败", err)
			return
		}

		//2. 在信息中枢处根据消息类型进行特定的处理
		switch requestPkg.Type {
		case pojo.CertificationType:
			//用户认证
			client.CertificationProcess(requestPkg)

		case pojo.CreateRoomType:
			//创建房间号,并将创建者加入房间
			client.CreateRoomProcess()

		case pojo.JoinRoomType:
			//1.加入房间的前提,先建立连接
			//2.完成用户认证
			//3.发送消息类型和房间号 Type uuid
			//只有完成上述步骤,才可以加入房间
			var data map[string]interface{}
			err = json.Unmarshal([]byte(requestPkg.Data), &data)
			if err != nil {
				fmt.Println("解析 JSON 失败：", err)
				return
			}
			uuidValue, ok := data["uuid"].(string)
			if !ok {
				fmt.Println("uuid 字段不存在或不是字符串类型")
				return
			}
			client.JoinRoomProcess(uuidValue)

		case pojo.RefreshScoreType:
			//什么是否进行分数更新，前端判断 type:RefreshScoreType, data:step、step、score
			//当用户的行为触发前端游戏机制的更新时，前端调用此接口，后端进行分数的转发 不需要做业务处理，直接转发即可
			fmt.Println("游戏交换中数据", client)
			client.RefreshScoreProcess(requestPkg)

		case pojo.DiscontinueQuitType:
			client.DiscontinueQuitProcess()

		case pojo.GameOverType:
			//游戏结束类型好像没有太大用，游戏结束的时候的提醒，通过分数更新就可以实现了
			fmt.Println("GameOverType")

		case pojo.HeartCheckType:
			//开启一个协程遍历hub中的Client，进行健康检测，生命时间是否会过期，如果过期进行逻辑删除和关闭连接
			if requestPkg.Data == "PING" {
				client.HeartCheckProcess()
			}
		}

	}

}
