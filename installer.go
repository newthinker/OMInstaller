package main

import (
	"fmt"
	"github.com/newthinker/onemap-installer/web"
	"net/http"
)

func main() {
	////////////////////////////////////////////////////////////////
	fmt.Println("web test")

	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))
	http.Handle("/img/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/subconfig", web.SubHandler)
	http.HandleFunc("/sysconfig", web.SysConfig)
	http.HandleFunc("/syshandler", web.SysHandler)
	http.HandleFunc("/error", web.ErrHandler)

	err := http.ListenAndServe("192.168.80.98:8888", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		return
	}
	///////////////////////////////////////////////////////////////	
}
