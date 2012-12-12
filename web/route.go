package web

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
)

// 分平台处理器
func SubHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SubPlatform handler")

	sub := &subController{}
	controller := reflect.ValueOf(sub)
	method := controller.MethodByName("SelectAction")

	if !method.IsValid() {
		/// default controller
	}

	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})

}

// 参数配置处理器
func SysconfHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SysConfig handler")
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		t, _ := template.ParseFiles("template/sysconfig.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		err := r.ParseForm() // 解析URL传递的参数

		// 如果分平台参数解析有问题，报告错误并返回
		if err != nil {
			http.Redirect(w, r, "/sysconfig", http.StatusFound)
		} else {
			fmt.Println(r.Form[""])
			//			http.Redirect(w, r, "/sysconfig", http.StatusFound)
		}
	}
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
