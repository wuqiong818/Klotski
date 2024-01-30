package router

import (
	"klotski/socket"
	test "klotski/src"
	"net/http"
)

func Run() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))
	http.Handle("/pages/", http.StripPrefix("/pages/", http.FileServer(http.Dir("views/pages"))))
	Group()
	http.ListenAndServe("127.0.0.1:8090", nil)
}
func Group() {
	WebSocketGroup()
	TestGroup()
}

func WebSocketGroup() {
	http.HandleFunc("/serverConnect", socket.ServerConnect)
	//http.HandleFunc("/clientConnect", socket.ClientConnect)
}

func TestGroup() {
	http.HandleFunc("/hello", test.Hello)
	//http.HandleFunc("/clientConnect", socket.ClientConnect)
}
