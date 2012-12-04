package web

import (
    "net/http"
    "strings"
    "reflect"
    "html/template"
    "log"
)

// 分平台处理器
func subHandler(w http.ResponseWriter, r *http.Request) {
    
}

// 参数配置处理器
func sysconfHandler(w http.ResponseWriter, r *Http.Request) {

}

// 404页面处理器 
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    // 如果路径是"/"，就跳转到首页
    if r.URL.Path == "/" {
        http.Redirect(w, r, "/template/html/index.html", http.StatusFound)
    }

    // 如果访问路径不满足制定的路由，就读取显示404模板
    t, err := template.ParseFiles("template/html/404.html")
    if (err!=nil) {
        log.Println(err)
    }

    t.Execute(w, nil)
}

