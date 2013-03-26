package web

import (
	"github.com/newthinker/onemap-installer/sys"
	"html/template"
	"net/http"
)

type msgHandler struct {
}

// 处理用户菜单选择操作
func (this *msgHandler) SelectAction(w http.ResponseWriter, r *http.Request) {
	l.Messagef("method:", r.Method)

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/msglist.html")
		if err != nil {
			l.Error(err)
		}
		t.Execute(w, nil)

		go sys.FormatResult(0, "Test message", nil)
	} else if r.Method == "POST" {
		t, err := template.ParseFiles("template/msglist.html")
		if err != nil {
			l.Error(err)
		}
		t.Execute(w, nil)
	}
}
