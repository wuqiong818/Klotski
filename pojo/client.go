package pojo

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Client struct {
	RoomId      string                        `json:"roomId"`
	UserId      string                        `json:"userId"`
	UserName    string                        `json:"userName"`
	UserProfile string                        `json:"userProfile"`
	UserCer     bool                          `json:"userCer"`
	UserConn    *websocket.Conn               `json:"-"`
	Hub         map[string]map[string]*Client `json:"-"` //第一个string-roomId 第二个string-userId
	Score       int                           `json:"score"`
	HealthCheck time.Time                     `json:"-"`
}

// Certificate 客服端认证
func (c *Client) Certificate(certificationInfo *CertificationInfo) (bool, error) {
	if certificationInfo.UserId != "" && certificationInfo.UserName != "" && certificationInfo.UserProfile != "" {
		c.UserId = certificationInfo.UserId
		c.UserName = certificationInfo.UserName
		c.UserProfile = certificationInfo.UserProfile
		c.UserCer = true
		return true, nil
	} else {
		c.UserCer = false
		err := errors.New("用户认证失败")
		return false, err
	}
}

// CreateRoomCode 创建房间 -- 用户连接存在并且用户认证通过
func (c *Client) CreateRoomCode() {
	if c.UserConn != nil && c.UserCer {
		uid := uuid.NewV4()
		c.RoomId = uid.String()
	}
}

// JoinHub 将连接加入中心 前提RoomId不为空, 加入房间的时候需要检测当前房间里面的人数
func (c *Client) JoinHub() (flag bool) {
	//	先查询是否存在此一个roomId key 加入房间
	if myMap, ok := c.Hub[c.RoomId]; ok {
		//检测人数
		count := len(myMap)
		if count == 1 {
			myMap[c.UserId] = c
			flag = true
		}
	} else { //创建房间
		myMap := make(map[string]*Client)
		myMap[c.UserId] = c     //userId
		c.Hub[c.RoomId] = myMap //roomId
		flag = true
	}
	return
}

// DeleteFromHub 逻辑删除 将连接从中心里删除 ,删除c它自己,前提roomId不为空
func (c *Client) DeleteFromHub() {
	if c.RoomId == "" {
		return
	}

	if value, ok1 := c.Hub[c.RoomId]; ok1 {
		if _, ok2 := value[c.UserId]; ok2 {
			delete(value, c.UserId)
		}
	}
	if len(c.Hub[c.RoomId]) == 0 {
		delete(c.Hub, c.RoomId)
	}
}

// QueryOtherUser 根据当前用户寻找另一位用户，返回user对象
func (c *Client) QueryOtherUser() *Client {
	if roomMap, ok := c.Hub[c.RoomId]; ok { //room
		for userId, user := range roomMap {
			if userId != c.UserId {
				return user
			}
		}
	}
	return nil
}

// ProlongLife 延长生命时间
func (c *Client) ProlongLife() {
	fmt.Println("延长时间")
	c.HealthCheck = time.Now().Add(10 * time.Second)
}
