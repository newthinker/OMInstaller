package web

import (
	"fmt"
	"net/http"
	//    "strings"
	//"reflect"
	"html/template"
	//"log"
)

// 分平台处理器
func SubHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SubPlatform handler")
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/subconfig.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm() // 解析URL传递的参数

		fmt.Println(r.Form)
	}

}

// 参数配置处理器
func SysconfHandler(w http.ResponseWriter, r *http.Request) {

}

// 404页面处理器 
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// 如果路径是"/"，就跳转到首页
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/subplatform", http.StatusFound)
	}

	// 如果访问路径不满足制定的路由，就读取显示404模板
	t, err := template.ParseFiles("template/404.html")
	if err != nil {
		fmt.Println(err)
	}

	t.Execute(w, nil)
}
