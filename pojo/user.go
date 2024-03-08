package pojo

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"klotski/pkg"
	"sync"
	"time"
)

type User struct {
	RoomId      string          `json:"roomId"`
	UserId      string          `json:"userId"`
	UserName    string          `json:"userName"`
	UserProfile string          `json:"userProfile"`
	UserCer     bool            `json:"userCer"` //是否验证通过
	UserConn    *websocket.Conn `json:"-"`
	HealthCheck time.Time       `json:"-"`
	Closed      bool
	mutex       sync.Mutex
}

// Certificate 客服端身份认证
func (u *User) Certificate(certificationInfo *CertificationInfo) (bool, error) {
	if certificationInfo.UserId != "" && certificationInfo.UserName != "" && certificationInfo.UserProfile != "" {
		u.UserId = certificationInfo.UserId
		u.UserName = certificationInfo.UserName
		u.UserProfile = certificationInfo.UserProfile
		u.UserCer = true
		return true, nil
	} else {
		u.UserCer = false
		err := errors.New("用户认证失败")
		return false, err
	}
}

// CreateRoomCode 创建房间 -- 用户连接存在并且用户认证通过
func (u *User) CreateRoomCode() {
	if u.UserConn != nil && u.UserCer {
		uid := uuid.NewV4()
		u.RoomId = uid.String()
	}
}

// Close 关闭连接
func (u *User) Close() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	if !u.Closed {
		if err := u.UserConn.Close(); err != nil {
			fmt.Println("关闭用户连接 err", err)
		}
		u.Closed = true
	}
}

// ProlongLife 延长生命时间
func (u *User) ProlongLife() {
	u.HealthCheck = time.Now().Add(time.Duration(pkg.HeartCheckSecond) * time.Second)
}
