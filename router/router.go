package router

import (
	"fmt"
	"klotski/httpLink"
	"klotski/pojo"
	"klotski/socket"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run() {

	go startHttpServer()

	hub := pojo.NewHub()
	go startWebSocketServer(hub)
	go hub.Run()

	waitForTerminationSignal()

}

func startWebSocketServer(hub *pojo.HupCenter) {
	http.HandleFunc("/ws", socket.WsHandle(hub))
	//长连接
	if err := http.ListenAndServe(":8091", nil); err != nil {
		fmt.Println("start-websocket-server fail", err)
	}
}

func startHttpServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))

	http.HandleFunc("/welcome2", httpLink.Welcome)
	http.HandleFunc("/create", httpLink.CreateUser)
	http.HandleFunc("/delete", httpLink.DeleteUser)
	http.HandleFunc("/getId", httpLink.GetUser)
	http.HandleFunc("/getAll", httpLink.GetAll)
	//短连接
	if err := http.ListenAndServe(":8090", nil); err != nil {
		fmt.Println("start-http-server fail", err)
	}
}

func waitForTerminationSignal() {
	// 创建一个接收终止信号的通道
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan

	log.Println("Server has been gracefully stopped.")
}
