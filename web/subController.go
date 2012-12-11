package web

import (
    "net/http"
    "html/template"
    "log"
)

type subController struct {
    SelectID []int          // 用户选择的所有菜单
}


// 处理用户菜单选择操作
func (this *subController)SelectAction(w http.ResponseWriter, r http.Request) {
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/subconfig.html")
        if err!=nil {
            log.Println(err)
        }
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		err := r.ParseForm() // 解析URL传递的参数

		// 如果分平台参数解析有问题，报告错误并返回
		if err != nil {
            log.Println(err)
		} else {    // 解析用户选择menu并进入下个页面
			fmt.Println(r.Form["selValues"])

            this.SelectID = r.Form["selValues"]

            // 进入系统参数配置页面
			http.Redirect(w, r, "/sysconfig", http.StatusFound)
		}

}

