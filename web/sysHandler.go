package web

import (
	"html/template"
	"net/http"
)

type sysHandler struct {
}

// 处理用户菜单选择操作
func (this *sysHandler) SelectAction(w http.ResponseWriter, r *http.Request) {
	l.Messagef("method:", r.Method)

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/sysconfig.html")
		if err != nil {
			l.Error(err)
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		t, err := template.ParseFiles("template/sysconfig.html")
		if err != nil {
			l.Error(err)
		}
		t.Execute(w, nil)
	}
}
