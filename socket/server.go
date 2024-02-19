package socket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"klotski/pojo"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub 保证这里的hub变量是唯一的，并且所有的Client对象操作的是同一个
var Hub = make(map[string]map[string]*pojo.Client)

func BuildConn(w http.ResponseWriter, r *http.Request) {
	//1. 服务端与客户端建立连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("conn错误", err)
		return
	}
	fmt.Println("远程主机连接成功 IP为", conn.RemoteAddr())
	//进行client的初始化操作
	var client pojo.Client
	client.Hub = Hub
	client.UserConn = conn
	client.HealthCheck = time.Now().Add(10 * time.Second)

	//计时器 : 如果用户在7秒内没有完成用户认证，则直接断开连接
	time.AfterFunc(7*time.Second, func() {
		if !client.UserCer {
			err := conn.Close()
			if err != nil {
				fmt.Println("用户认证失败，关闭连接时遇到err", err)
				return
			}
		}
	})

	go func() {
		for {
			//对当前Hub中的连接进行心跳检测
			if client.HealthCheck.Before(time.Now()) {
				client.DeleteFromHub()
				err2 := client.UserConn.Close()
				if err2 != nil {
					fmt.Println("心跳检测中的关闭连接err", err)
				}
				break
			}
		}
	}()

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
		//2. 在信息中枢处根据消息类型进行特定的处理
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
			responsePkg.Type = pojo.CertificationType

			if isCer {
				responsePkg.Code = pojo.Success
			} else {
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
					return
				}
			}
		case pojo.CreateRoomType: //创建房间号,并将创建者加入房间
			responsePkg := new(pojo.ResponsePkg)
			responsePkg.Type = pojo.CreateRoomType
			if client.UserCer {
				client.CreateRoomCode()
				responsePkg.Code = pojo.Success
				responsePkg.Data = client.RoomId
				//创建房间的同时,加入房间
				client.JoinHub()
				fmt.Println(client.Hub)
			} else {
				responsePkg.Code = pojo.UserCertificationError
				responsePkg.Data = "用户认证失败,请先完成用户认证"
			}
			err := client.UserConn.WriteJSON(responsePkg)
			if err != nil {
				fmt.Println("服务端返回房间信息error", err)
				return
			}
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
			client.RoomId = uuidValue

			responsePkg := new(pojo.ResponsePkg)
			responsePkg.Type = pojo.JoinRoomType
			if client.RoomId != "" {
				flag := client.JoinHub()
				fmt.Println("JoinRoom", client.Hub)
				if flag {
					//第二个人加入房间后，向房间里面的双方互相广播对方的个人信息，发送两个read
					//注意，这里可能会存在程序漏洞，逻辑上的，当一个房间里面的人数多于两个人是出错
					responsePkg.Code = pojo.Success
					if roomMap, ok := client.Hub[client.RoomId]; ok {
						for userId, user := range roomMap {
							if userId != client.UserId {
								//向本客户端发送其他用户的信息
								//fmt.Println("向本人发信息")
								marshal1, _ := json.Marshal(user)
								responsePkg.Data = string(marshal1)
								err1 := client.UserConn.WriteJSON(responsePkg)
								if err1 != nil {
									fmt.Println("向本人发送除本人意外的其他信息失败 err", err)
								}
								//向其他客户端发送本客户端的信息
								//fmt.Println("向其他人发信息", user)
								marshal2, _ := json.Marshal(client)
								responsePkg.Data = string(marshal2)
								err2 := user.UserConn.WriteJSON(responsePkg)
								if err2 != nil {
									fmt.Println("向其他人返回的信息 err", err)
								}
							}
						}
					}
				} else {
					responsePkg.Code = pojo.JoinRoomError
					responsePkg.Data = "加入房间失败 失败原因，房间已满"
					err := client.UserConn.WriteJSON(responsePkg)
					if err != nil {
						fmt.Println("服务端返回加入房间信息error", err)
					}
					client.UserConn.Close()
				}
			} else {
				responsePkg.Code = pojo.JoinRoomError
				responsePkg.Data = "加入房间失败 错误原因RoomId未赋值"
				err := client.UserConn.WriteJSON(responsePkg)
				if err != nil {
					fmt.Println("服务端返回加入房间信息error", err)
				}
			}
		case pojo.RefreshScoreType:
			//什么是否进行分数更新，前端判断 type:RefreshScoreType, data:step、step、score
			//当用户的行为触发前端游戏机制的更新时，前端调用此接口，后端进行分数的转发 不需要做业务处理，直接转发即可
			//1.读取数据
			var gameData pojo.GameData
			err = json.Unmarshal([]byte(requestPkg.Data), &gameData)
			gameDataPkg := new(pojo.ResponsePkg)
			gameDataPkg.Type = pojo.RefreshScoreType
			if gameData.Score != "" && gameData.StepCount != "" && gameData.TakeTime != "" {
				gameDataPkg.Code = pojo.Success
				gameDataPkg.Data = requestPkg.Data
				user := client.QueryOtherUser()
				err := user.UserConn.WriteJSON(gameDataPkg)
				if err != nil {
					fmt.Println("转发游戏数据时 err", err)
					return
				}
			} else {
				gameDataPkg.Code = pojo.RefreshScoreError
				gameDataPkg.Data = ""
				err := client.UserConn.WriteJSON(gameDataPkg)
				if err != nil {
					fmt.Println("转发游戏数据时 err", err)
					return
				}
			}
		case pojo.DiscontinueQuitType:
			//用户自行离开或者网络异常 断开连接
			client.DeleteFromHub()
			err := client.UserConn.Close()
			if err != nil {
				fmt.Println("DiscontinueQuit err", err)
			}
		case pojo.GameOverType:
			//游戏结束类型好像没有太大用，游戏结束的时候的提醒，通过分数更新就可以实现了
			fmt.Println("GameOverType")

		case pojo.HeartCheckType:
			//开启一个协程遍历hub中的Client，进行健康检测，生命时间是否会过期，如果过期进行逻辑删除和关闭连接
			if requestPkg.Data == "PING" {
				client.ProlongLife()
				responsePkg := pojo.ResponsePkg{
					Type: pojo.HeartCheckType,
					Code: pojo.Success,
					Data: "PONG",
				}
				err := client.UserConn.WriteJSON(responsePkg)
				if err != nil {
					fmt.Println("心跳检测发送err", err)
				}
			}

		}
	}
	defer conn.Close()
	//log.Println("服务关闭")
}
