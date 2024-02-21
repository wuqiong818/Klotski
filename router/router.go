package router

import (
	"klotski/httpLink"
	"klotski/socket"
	"net/http"
)

func Run() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))
	//httpLink.Handle("/pages/", httpLink.StripPrefix("/pages/", httpLink.FileServer(httpLink.Dir("views/pages"))))
	Group()
	//http.ListenAndServe("172.19.22.102:8090", nil) //服务器私网 172.19.22.102
	//http.ListenAndServe("8.141.88.60:8090", nil) //公网
	http.ListenAndServe("localhost:8090", nil) //电脑本机的
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
