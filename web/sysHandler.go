package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type sysHandler struct {
}

// 处理用户菜单选择操作
func (this *sysHandler) SelectAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/sysconfig.html")
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		t, err := template.ParseFiles("template/sysconfig.html")
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
	}
}
