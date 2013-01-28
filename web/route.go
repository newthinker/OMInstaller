package web

import (
	"github.com/newthinker/onemap-installer/log"
	"github.com/newthinker/onemap-installer/sys"
	"html/template"
	"net/http"
	"reflect"
)

var (
	l *(log.Logger)
)

func Init(logger *(log.Logger)) {
	l = logger
	sys.Init(l)
}

// 分平台处理器
func SubHandler(w http.ResponseWriter, r *http.Request) {
	l.Message("SubPlatform handler")

	sub := &subController{}
	controller := reflect.ValueOf(sub)
	method := controller.MethodByName("SelectAction")

	if !method.IsValid() {
		l.Errorf("Invalid input params")
		OutputJson(w, 1, "输入参数非法", nil)
		return
	}

	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
}

// SysConfig页面
func SysConfig(w http.ResponseWriter, r *http.Request) {
	l.Message("SysHandler handler")

	sys := &sysHandler{}
	controller := reflect.ValueOf(sys)
	method := controller.MethodByName("SelectAction")

	if !method.IsValid() {
		l.Errorf("Invalid input params")
		OutputJson(w, 1, "输入参数非法", nil)
		return
	}

	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
}

// 参数配置处理器
func SysHandler(w http.ResponseWriter, r *http.Request) {
	l.Message("SysConfig handler")

	sys := &sysController{}
	/*	if err := sys.Init(); err != nil {
			l.Errorf("Sys module init failed")
			OutputJson(w, 2, "系统初始化错误!", nil)
			return
		}
	*/
	controller := reflect.ValueOf(sys)
	method := controller.MethodByName("SysAction")

	if !method.IsValid() {
		l.Errorf("Invalid input params")
		OutputJson(w, 1, "非法输入参数!", nil)
		return
	}

	requestValue := reflect.ValueOf(r)
	responseValue := reflect.ValueOf(w)
	method.Call([]reflect.Value{responseValue, requestValue})
}

// sysconfig页面中参数检查错误提醒页面
func ErrHandler(w http.ResponseWriter, r *http.Request) {
	l.Message("Error handler")

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/error.html")
		if err != nil {
			l.Error(err)
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
		l.Error(err)
	}

	t.Execute(w, nil)
}
