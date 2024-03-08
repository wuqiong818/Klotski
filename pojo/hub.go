package pojo

import (
	"fmt"
	"klotski/pkg"
	"time"
)

type HupCenter struct {
	ClientsMap map[string]map[string]*Client `json:"-"` //第一个string-roomId 第二个string-userId
	Register   chan *Client
	UnRegister chan *Client
}

// NewHub 初始化一个hub
func NewHub() *HupCenter {
	return &HupCenter{
		ClientsMap: make(map[string]map[string]*Client),
		Register:   make(chan *Client, 1),
		UnRegister: make(chan *Client, 1),
	}
}

// Run 用户向hub中的逻辑注册、删除、心跳检测全方法
func (h *HupCenter) Run() {
	checkTicker := time.NewTicker(time.Duration(pkg.HeartCheckSecond) * time.Second)
	defer checkTicker.Stop()

	for {
		select {
		case client := <-h.Register:
			//先查询是否存在此一个roomId key
			if myMap, ok := client.Hub.ClientsMap[client.User.RoomId]; ok { //有,加入房间
				//检测人数
				if len(myMap) == 1 {
					myMap[client.User.UserId] = client
				}
			} else { //没有,创建房间
				myMap := make(map[string]*Client)
				myMap[client.User.UserId] = client                //userId
				client.Hub.ClientsMap[client.User.RoomId] = myMap //roomId
			}
			fmt.Println("有人加入房间:", client.Hub.ClientsMap)
		case client := <-h.UnRegister:
			client.User.Close()
			if value, ok1 := client.Hub.ClientsMap[client.User.RoomId]; ok1 {
				if _, ok2 := value[client.User.UserId]; ok2 {
					delete(value, client.User.UserId)
				}
			}
			if len(client.Hub.ClientsMap[client.User.RoomId]) == 0 {
				delete(client.Hub.ClientsMap, client.User.RoomId)
			}
		case <-checkTicker.C:
			for _, roomMap := range h.ClientsMap {
				//fmt.Println(roomMap)
				for _, client := range roomMap {
					//fmt.Println(client)
					//fmt.Println(client.User.HealthCheck)
					if client.User.HealthCheck.Before(time.Now()) {
						h.UnRegister <- client
					}
				}
			}
			fmt.Println(time.Now().Format(time.DateTime), h.ClientsMap)
		}
	}
}

// QueryOtherUser 读操作 -- 根据当前用户寻找另一位用户，返回user对象
func (h *HupCenter) QueryOtherUser(c *Client) *Client {
	if roomMap, ok := h.ClientsMap[c.User.RoomId]; ok { //room
		for userId, user := range roomMap {
			if userId != c.User.UserId {
				return user
			}
		}
	}
	return nil
}

/*
//---使用锁，操作一个多线程共享的Map---//

type HupCenter struct {
	ClientsMap map[string]map[string]*Client `json:"-"` //第一个string-roomId 第二个string-userId
	mutex      sync.RWMutex

}
// JoinHub  写操作 --将连接加入中心 前提RoomId不为空, 加入房间的时候需要检测当前房间里面的人数
func (h *HupCenter) JoinHub(c *Client) (flag bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	//先查询是否存在此一个roomId key
	if myMap, ok := c.Hub.ClientsMap[c.User.RoomId]; ok { //有,加入房间
		//检测人数
		if len(myMap) == 1 {
			myMap[c.User.UserId] = c
			flag = true
		}
	} else { //没有,创建房间
		myMap := make(map[string]*Client)
		myMap[c.User.UserId] = c                //userId
		c.Hub.ClientsMap[c.User.RoomId] = myMap //roomId
		flag = true
	}
	return
}

// DeleteFromHub 写操作 --逻辑删除 将传入的参数c从hub连接池中删除
func (h *HupCenter) DeleteFromHub(c *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if c.User.RoomId == "" {
		return
	}

	if value, ok1 := c.Hub.ClientsMap[c.User.RoomId]; ok1 {
		if _, ok2 := value[c.User.UserId]; ok2 {
			delete(value, c.User.UserId)
		}
	}

	if len(c.Hub.ClientsMap[c.User.RoomId]) == 0 {
		delete(c.Hub.ClientsMap, c.User.RoomId)
	}
}

// QueryOtherUser 读操作 -- 根据当前用户寻找另一位用户，返回user对象
func (h *HupCenter) QueryOtherUser(c *Client) *Client {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if roomMap, ok := h.ClientsMap[c.User.RoomId]; ok { //room
		for userId, user := range roomMap {
			if userId != c.User.UserId {
				return user
			}
		}
	}
	return nil
}
*/
