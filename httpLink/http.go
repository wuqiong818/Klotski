package httpLink

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"klotski/pojo"
	"net/http"
	"reflect"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //设置请求头
	fmt.Println("welcome!!!")
	fmt.Fprintln(w, "服务器返回的信息<h1>Welcome March</h1>") //打印
}

type User struct {
	ID    string `json:"userId"`
	Name  string `json:"userName"`
	Time  int    `json:"time"`
	Steps int    `json:"foot"`
}

// 创建一个 redis 客户端
var rdb = redis.NewClient(&redis.Options{
	Addr:     "8.141.88.60:6379",
	Password: "123456",
	DB:       0,
})

// 定义一个上下文
var ctx = context.Background()

// CreateUser 定义一个处理器函数，用于创建用户数据
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //设置请求头
	responsePkg := new(pojo.ResponsePkg)
	responsePkg.Type = pojo.CreateRoomType
	// 从请求体中读取用户数据
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//fmt.Println("data类型", reflect.TypeOf(data))
	// 将用户数据转换为 JSON 格式

	data, err := json.Marshal(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 将用户数据存储到 redis 中，以用户 ID 作为键
	err = rdb.Set(ctx, user.ID, data, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回成功的响应
	responsePkg.Code = 200
	responsePkg.Data = string(data)
	resp, err := json.Marshal(responsePkg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(resp)
}

// 定义一个处理器函数，用于获取用户数据
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //设置请求头
	responsePkg := new(pojo.ResponsePkg)
	responsePkg.Type = pojo.CreateRoomType
	// 从请求参数中获取用户 ID
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	// 从 redis 中获取用户数据，以用户 ID 作为键
	data, err := rdb.Get(ctx, id).Bytes()

	if err != nil {
		if err == redis.Nil {
			http.Error(w, "user not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	responsePkg.Code = 200
	responsePkg.Data = string(data)
	resp, err := json.Marshal(responsePkg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回成功的响应
	w.WriteHeader(200)
	w.Write(resp)
}

// 定义一个处理器函数，用于查询全部数据
func GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //设置请求头
	responsePkg := new(pojo.ResponsePkg)
	responsePkg.Type = pojo.CreateRoomType
	// 从 redis 中获取所有的用户 ID
	ids, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 从 redis 中批量获取所有的用户数据
	data, err := rdb.MGet(ctx, ids...).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("data类型", reflect.TypeOf(data))
	// 将用户数据转换为 JSON 格式
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回成功的响应
	responsePkg.Code = 200
	responsePkg.Data = string(jsonData)
	resp, err := json.Marshal(responsePkg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(resp)
}

// 定义一个处理器函数，用于更新用户数据
//func updateUser(w http.ResponseWriter, r *http.Request) {
//	// 从请求参数中获取用户 ID
//	id := r.URL.Query().Get("id")
//	if id == "" {
//		http.Error(w, "missing id", http.StatusBadRequest)
//		return
//	}
//	// 从请求体中读取用户数据
//	var user User
//	err := json.NewDecoder(r.Body).Decode(&user)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	// 将用户数据转换为 JSON 格式
//	data, err := json.Marshal(user)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	// 将用户数据更新到 redis 中，以用户 ID 作为键
//	err = rdb.Set(ctx, id, data, 0).Err()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	// 返回成功的响应
//	w.WriteHeader(http.StatusOK)
//	w.Write(data)
//}

// DeleteUser 定义一个处理器函数，用于删除用户数据
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //设置请求头
	// 从请求参数中获取用户 ID
	responsePkg := new(pojo.ResponsePkg)
	responsePkg.Type = pojo.CreateRoomType
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", 5000)
		return
	}

	// 从 redis 中删除用户数据，以用户 ID 作为键
	err := rdb.Del(ctx, id).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回成功的响应
	responsePkg.Code = 200
	responsePkg.Data = "OK"
	resp, err := json.Marshal(responsePkg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Write(resp)

}
