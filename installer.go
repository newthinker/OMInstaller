package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"github.com/newthinker/onemap-installer/log"
	"github.com/newthinker/onemap-installer/utl"
	"github.com/newthinker/onemap-installer/web"
	"net/http"
	"path/filepath"
	"runtime"

//	"syscall"
)

func main() {
	/*	
		// check the windows's version
		dll := syscall.MustLoadDLL("kernel32.dll")
		p := dll.MustFindProc("GetVersion")
		v, _, _ := p.Call()
		fmt.Printf("Windows version %d.%d (Build %d)\n", byte(v), uint8(v>>8), uint16(v>>16))
	*/
	////////////////////////////////////////////////////////////////
	// init log
	l, err := log.NewLog("inst.log", log.LogAll, log.DefaultBufSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	// get local ip
	ip, err := utl.GetNetIP()
	if err != nil {
		l.Errorf("Get local ip failed")
		return
	}
	l.Messagef("Get local ip: %s", ip)
	////////////////////////////////////////////////////////////////
	switch runtime.GOOS {
	case "windows":

	case "linux":
		// install sshpass
		l.Message("Start to install sshpass")
		base, err := filepath.Abs("./") // 获取系统当前路径
		fmt.Println("base:" + base)
		if err != nil || base == "" {
			l.Errorf("Get current directory failed")
			return
		}

		if err := utl.InstallSshpass(base); err != nil {
			l.Errorf("Install sshpass failed")
			return
		}

		if err := utl.CheckSshpass(base); err != nil {
			l.Errorf("Sshpass isn't installed")
			return
		}
		l.Message("Installing sshpass successfully")
	}

	////////////////////////////////////////////////////////////////
	l.Message("Listen and serve")
	web.Init(l)

	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))
	http.Handle("/images/", http.FileServer(http.Dir("template")))

	http.Handle("/json", websocket.Handler(web.JsonServer))

	http.HandleFunc("/subconfig", web.SubHandler)
	http.HandleFunc("/sysconfig", web.SysConfig)
	http.HandleFunc("/syshandler", web.SysHandler)
	http.HandleFunc("/error", web.ErrHandler)
	http.HandleFunc("/", web.MainHandler)

	err = http.ListenAndServe(ip+":8888", nil)
	if err != nil {
		fmt.Println(err)
		l.Errorf("Listen and serve failed: %s", err)
		return
	}
	///////////////////////////////////////////////////////////////	
}
