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
        OutputJson(w, 1, "输入参数非法", nil);
        return
	}

	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})

}

// SysConfig页面
func SysConfig(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SysHandler handler")

	sys := &sysHandler{}
	controller := reflect.ValueOf(sys)
	method := controller.MethodByName("SelectAction")

	if !method.IsValid() {
        OutputJson(w, 1, "输入参数非法", nil);
        return
	}

	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})

}

// 参数配置处理器
func SysHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SysConfig handler")

	sys := &sysController{}
	if err := sys.Init(); err != nil {
		OutputJson(w, 2, "系统初始化错误!", nil)
		return
	}
	controller := reflect.ValueOf(sys)
	method := controller.MethodByName("SysAction")

	if !method.IsValid() {
        OutputJson(w, 1, "非法输入参数!", nil)
        return
	}

	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})

}

// sysconfig页面中参数检查错误提醒页面
func ErrHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Error handler")

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/error.html")
		if err != nil {
			fmt.Println(err)
		}
		t.Execute(w, nil)
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
