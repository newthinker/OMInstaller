package main

import (
	"fmt"
	"github.com/newthinker/onemap-installer/sys"
	"github.com/newthinker/onemap-installer/web"
	"net/http"
	"os/exec"
	"path/filepath"
)

func main() {
	////////////////////////////////////////////////////////////////
	// install sshpass
	base, err := filepath.Abs("./") // 获取系统当前路径
	fmt.Println("base:" + base)
	if err != nil || base == "" {
		fmt.Println("ERROR: Get current directory failed")
		return
	}

	// whether installed
	if flag := sys.Exists(base + "/sshpass/bin/sshpass"); flag != true {
		if flag = sys.Exists(base + "/sshpass/Install.sh"); flag != true {
			fmt.Println("ERROR: No sshpass software package")
			return
		}

		// exec the install script
		cmd := exec.Command("/bin/sh", base+"/sshpass/Install.sh", base+"/sshpass")
		err = cmd.Run()
		if err != nil {
			fmt.Println("ERROR: Complier sshpass failed")
			fmt.Println(err)
			return
		}

		// whether install successfully
		if flag := sys.Exists(base + "/sshpass/bin/sshpass"); flag != true {
			fmt.Println("ERROR: Install sshpass failed")
			return
		}
	} else {
		fmt.Println("MSG: Sshpass is installed and go on")
	}

	cmd := exec.Command(base+"/sshpass/bin/sshpass", "-V")
	err = cmd.Run()
	if err != nil {
		fmt.Println("ERROR: Sshpass isn't installed")
		return
	}

	////////////////////////////////////////////////////////////////
	fmt.Println("web test")

	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))
	http.Handle("/images/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/subconfig", web.SubHandler)
	http.HandleFunc("/sysconfig", web.SysConfig)
	http.HandleFunc("/syshandler", web.SysHandler)
	http.HandleFunc("/error", web.ErrHandler)

	err = http.ListenAndServe("192.168.80.98:8888", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		return
	}
	///////////////////////////////////////////////////////////////	
}
