package web

import (
	"github.com/newthinker/onemap-installer/sys"
	"html/template"
	"net/http"
	"strconv"
)

type mainController struct {
}

func (this *mainController) SelectAction(w http.ResponseWriter, r *http.Request) {
	l.Messagef("Main page method:%s", r.Method)

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/index.html")
		if err != nil {
			l.Error(err)
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			l.Error(err)
			OutputJson(w, 1, err.Error(), nil)
		} else { // parse user's selection
			l.Debugf("The flag is:%s", (r.Form["flag"])[0])
			flag, err := strconv.Atoi((r.Form["flag"])[0])
			if err != nil {
				msg := "Invalid input params"
				l.Errorf(msg)
				OutputJson(w, 2, msg, nil)
			}
			if flag <= sys.MAINTAIN || flag > sys.SUBPLATFORM {
				OutputJson(w, 3, "Invalid input flag", nil)
			}

			switch flag {
			case sys.INSTALL:
				sys.SubFlag = false
				// reset the SysConfig struct
				sys.ResetSysConfig()

				// redirect to the sysconfig.html page
				http.Redirect(w, r, "/sysconfig", http.StatusFound)
			case sys.UPDATE: /// todo
				l.Debugf("The flag is %d", flag)
			case sys.UNINSTALL: /// todo
				l.Debugf("The flag is %d", flag)
			case sys.SUBPLATFORM:
				sys.SubFlag = true
				// redirect to the subconfig.html page
				http.Redirect(w, r, "/subconfig", http.StatusFound)
			}
		}
	}
}
