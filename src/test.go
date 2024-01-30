package test

import (
	"fmt"
	"net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("请求参数有", r.Form)
	fmt.Println("直接", r.Form["username"])
	fmt.Println("直接", r.Form["password"])
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(username)
	fmt.Println(password)
}
