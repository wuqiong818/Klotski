package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", hello)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodPost {
	//	fmt.Println("Method Not Allowed")
	//	return
	//}
	r.ParseForm()
	fmt.Println("请求参数有", r.Form)
	fmt.Println("直接", r.Form["username"])
	fmt.Println("直接", r.Form["password"])
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username)
	fmt.Println(password)
	//params := r.URL.Query()
	//name := params.Get("name")
}
