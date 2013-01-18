package main

import (
	"fmt"
	"github.com/newthinker/onemap-installer/log"
	"github.com/newthinker/onemap-installer/utl"
	"github.com/newthinker/onemap-installer/web"
	"net/http"
	"os/exec"
	"path/filepath"
)

func main() {
	////////////////////////////////////////////////////////////////
	// init log
	l, err := log.NewLog("inst.log", log.LogAll, log.DefaultBufSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	////////////////////////////////////////////////////////////////
	// get local ip
	ip, err := utl.GetNetIP()
	if err != nil {
		l.Errorf("Get local ip failed")
		return
	}
	l.Messagef("Get local ip: %s", ip)
	////////////////////////////////////////////////////////////////
	// install sshpass
	l.Message("Install the sshpass")
	base, err := filepath.Abs("./") // 获取系统当前路径
	fmt.Println("base:" + base)
	if err != nil || base == "" {
		l.Errorf("Get current directory failed")
		return
	}

	// whether installed
	if flag := utl.Exists(base + "/sshpass/bin/sshpass"); flag != true {
		if flag = utl.Exists(base + "/sshpass/Install.sh"); flag != true {
			l.Errorf("No sshpass software package")
			return
		}

		// exec the install script
		cmd := exec.Command("/bin/sh", base+"/sshpass/Install.sh", base+"/sshpass")
		err = cmd.Run()
		if err != nil {
			l.Errorf("Complier sshpass failed")
			return
		}

		// whether install successfully
		if flag := utl.Exists(base + "/sshpass/bin/sshpass"); flag != true {
			l.Errorf("Install sshpass failed")
			return
		}
	} else {
		l.Messagef("Sshpass is installed and go on")
	}

	cmd := exec.Command(base+"/sshpass/bin/sshpass", "-V")
	err = cmd.Run()
	if err != nil {
		l.Errorf("Sshpass isn't installed")
		return
	}

	////////////////////////////////////////////////////////////////
	l.Message("Listen and serve")
	web.Init(l)

	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))
	http.Handle("/images/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/subconfig", web.SubHandler)
	http.HandleFunc("/sysconfig", web.SysConfig)
	http.HandleFunc("/syshandler", web.SysHandler)
	http.HandleFunc("/error", web.ErrHandler)

	err = http.ListenAndServe(ip+":8888", nil)
	if err != nil {
		l.Errorf("Listen and serve failed: %s", err)
		return
	}
	///////////////////////////////////////////////////////////////	
}
