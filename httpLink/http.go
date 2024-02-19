package httpLink

import (
	"fmt"
	"net/http"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8") //设置请求头
	fmt.Fprintln(w, "服务器返回的信息<b>你好!!!</b>")                    //打印
}
