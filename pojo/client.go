package pojo

import (
	"encoding/json"
	"fmt"
)

type Client struct {
	User *User
	Hub  *HupCenter
}

func (c *Client) CertificationProcess(requestPkg RequestPkg) {
	//1.读取数据
	var certificationInfo CertificationInfo
	err := json.Unmarshal([]byte(requestPkg.Data), &certificationInfo)
	//2.进行用户验证
	isCer, _ := c.User.Certificate(&certificationInfo)

	//3.返回响应
	responsePkg := &ResponsePkg{
		Type: CertificationType,
	}

	if isCer {
		responsePkg.Code = Success
	} else {
		responsePkg.Code = UserCertificationError
		responsePkg.Data = "用户认证失败,非法侵入，服务端主动关闭连接"
	}
	err = c.User.UserConn.WriteJSON(responsePkg)
	if err != nil {
		fmt.Println("服务端用户认证信息返回错误", err)
	}
	if !isCer {
		err := c.User.UserConn.Close()
		if err != nil {
			fmt.Println("用户认证失败，服务端主动关闭连接错误", err)
			return
		}
	}
}

func (c *Client) CreateRoomProcess() {
	responsePkg := new(ResponsePkg)
	responsePkg.Type = CreateRoomType
	if c.User.UserCer {
		c.User.CreateRoomCode()
		responsePkg.Code = Success
		responsePkg.Data = c.User.RoomId
		//创建房间的同时,加入房间
		c.Hub.JoinHub(c)
		fmt.Println("加入房间信息", c.Hub.ClientsMap)
	} else {
		responsePkg.Code = UserCertificationError
		responsePkg.Data = "用户认证失败,请先完成用户认证"
	}
	err := c.User.UserConn.WriteJSON(responsePkg)
	if err != nil {
		fmt.Println("服务端返回房间信息error", err)
		return
	}
}

func (c *Client) JoinRoomProcess(uuid string) {
	c.User.RoomId = uuid
	responsePkg := new(ResponsePkg)
	responsePkg.Type = JoinRoomType

	if c.User.RoomId != "" {
		flag := c.Hub.JoinHub(c)
		fmt.Println("加入房间", c.Hub.ClientsMap)
		if flag {
			//第二个人加入房间后，向房间里面的双方互相广播对方的个人信息，发送两个read
			responsePkg.Code = Success
			if roomMap, ok := c.Hub.ClientsMap[c.User.RoomId]; ok {
				for userId, mapClient := range roomMap { //*Client
					if userId != c.User.UserId {
						//向本客户端发送其他用户的信息
						marshal1, _ := json.Marshal(mapClient.User)
						responsePkg.Data = string(marshal1)
						err1 := c.User.UserConn.WriteJSON(responsePkg)
						if err1 != nil {
							fmt.Println("向本人发送除本人意外的其他信息失败 err", err1)
						}
						//向其他客户端发送本客户端的信息
						marshal2, _ := json.Marshal(c.User)
						responsePkg.Data = string(marshal2)
						err2 := mapClient.User.UserConn.WriteJSON(responsePkg)
						if err2 != nil {
							fmt.Println("向其他人返回的信息 err", err2)
						}
					}
				}
			}
		} else {
			responsePkg.Code = JoinRoomError
			responsePkg.Data = "加入房间失败 失败原因，房间已满"
			err := c.User.UserConn.WriteJSON(responsePkg)
			if err != nil {
				fmt.Println("服务端返回加入房间信息error", err)
			}
			c.User.UserConn.Close()
		}
	} else {
		responsePkg.Code = JoinRoomError
		responsePkg.Data = "加入房间失败 错误原因RoomId未赋值"
		err := c.User.UserConn.WriteJSON(responsePkg)
		if err != nil {
			fmt.Println("服务端返回加入房间信息error", err)
		}
	}
}

func (c *Client) RefreshScoreProcess(requestPkg RequestPkg) {
	var gameData GameData
	err := json.Unmarshal([]byte(requestPkg.Data), &gameData)
	if err != nil {
		fmt.Println("读取游戏中的信息失败")
		return
	}
	gameDataPkg := new(ResponsePkg)
	gameDataPkg.Type = RefreshScoreType
	if gameData.Score != "" && gameData.StepCount != "" && gameData.TakeTime != "" {
		//游戏数据没有在本地进行存储，选择了直接进行转发
		gameDataPkg.Code = Success
		gameDataPkg.Data = requestPkg.Data
		otherClient := c.Hub.QueryOtherUser(c)
		err := otherClient.User.UserConn.WriteJSON(gameDataPkg)
		if err != nil {
			fmt.Println("转发游戏数据时 err", err)
			return
		}
	} else {
		gameDataPkg.Code = RefreshScoreError
		gameDataPkg.Data = ""
		err := c.User.UserConn.WriteJSON(gameDataPkg)
		if err != nil {
			fmt.Println("转发游戏数据时 err", err)
			return
		}
	}
}

func (c *Client) DiscontinueQuitProcess() {
	//用户自行离开或者网络异常 断开连接
	c.Hub.DeleteFromHub(c)
	fmt.Println("退出房间", c.Hub.ClientsMap)
	err := c.User.UserConn.Close()
	if err != nil {
		fmt.Println("DiscontinueQuit err", err)
	}
}

func (c *Client) HeartCheckProcess() {
	c.User.ProlongLife()
	responsePkg := ResponsePkg{
		Type: HeartCheckType,
		Code: Success,
		Data: "PONG",
	}
	err := c.User.UserConn.WriteJSON(responsePkg)
	if err != nil {
		fmt.Println("心跳检测发送err", err)
	}
}

/*

func (c *Client) CreateRoomProcess(requestPkg RequestPkg) {

}

*/

/*// Certificate 客服端身份认证
func (c *Client) Certificate(certificationInfo *CertificationInfo) (bool, error) {
	if certificationInfo.UserId != "" && certificationInfo.UserName != "" && certificationInfo.UserProfile != "" {
		c.User.UserId = certificationInfo.UserId
		c.User.UserName = certificationInfo.UserName
		c.User.UserProfile = certificationInfo.UserProfile
		c.User.UserCer = true
		return true, nil
	} else {
		c.User.UserCer = false
		err := errors.New("用户认证失败")
		return false, err
	}
}
*/

/*// CreateRoomCode 创建房间 -- 用户连接存在并且用户认证通过
func (c *Client) CreateRoomCode() {
	if c.UserConn != nil && c.UserCer {
		uid := uuid.NewV4()
		c.RoomId = uid.String()
	}
}

// ProlongLife 延长生命时间
func (c *Client) ProlongLife() {
	fmt.Println("延长时间")
	c.HealthCheck = time.Now().Add(10 * time.Second)
}
*/
