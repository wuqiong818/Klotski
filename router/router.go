package router

import (
	"fmt"
	"klotski/httpLink"
	"klotski/socket"
	"net/http"
)

func Run() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))
	//httpLink.Handle("/pages/", httpLink.StripPrefix("/pages/", httpLink.FileServer(httpLink.Dir("views/pages"))))
	Group()
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println("listenAndServe", err)
		return
	}
}
func Group() {
	WebSocketGroup()
	HttpGroup()
}

func WebSocketGroup() {
	//当在地址栏中输入ws://127.0.0.1:8090/serverConnect，建立长连接
	//自动开启协程了吗？
	//为什么每一个访问的远程IP都不相同
	http.HandleFunc("/ws", socket.BuildConn)

}

func HttpGroup() {
	http.HandleFunc("/welcome", httpLink.Welcome)
	http.HandleFunc("/create", httpLink.CreateUser)
	http.HandleFunc("/delete", httpLink.DeleteUser)
	http.HandleFunc("/getId", httpLink.GetUser)
	http.HandleFunc("/getAll", httpLink.GetAll)
}
