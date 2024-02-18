package router

import (
	"klotski/socket"
	test "klotski/src"
	"net/http"
)

func Run() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))
	//http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("views/pages"))))
	Group()
	http.ListenAndServe("127.0.0.1:8090", nil)
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

	http.HandleFunc("/hello", test.Hello)
	http.HandleFunc("/", test.Index)

}
