package web

import (
	"encoding/json"
	"errors"
	"github.com/newthinker/onemap-installer/sys"
	"github.com/newthinker/onemap-installer/utl"
	"net/http"
	"path/filepath"
)

type sysController struct {
}

// 处理系统配置页面
func (this *sysController) SysAction(w http.ResponseWriter, r *http.Request) {
	l.Messagef("Sysconfig page method:", r.Method)

	w.Header().Set("content-type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("content-type", "text/json;charset=utf-8")

	rate := sys.BEGIN // process rate
	// 将配置文件解析后传入前端显示
	if r.Method == "GET" {
		sysmap, err := sys.SysFormat(sys.INSTALL)
		if err != nil {
			l.Error(errors.New("Format system params failed"))
			OutputJson(w, 3, "格式化系统参数失败", nil)
			return
		}

		OutputJson(w, 0, "", sysmap)
	} else if r.Method == "POST" { // 将前端输入参数传入后台解析
		go sys.FormatResult(rate, utl.NowString()+"Begin to the process", nil)

		err := r.ParseForm()
		if err != nil {
			l.Error(err)
			go sys.FormatResult(sys.BREAK, utl.NowString()+err.Error(), nil)
			return
		}

		input := (r.Form["input"])[0]
		l.Debugf("Input params:%s", input)

		// 获取json数据流
		rate += sys.GET_JSON
		go sys.FormatResult(rate, utl.NowString()+"Get the json string from the front-end", nil)
		jsonstr := make(map[string]interface{})
		err = json.Unmarshal([]byte(input), &jsonstr)
		if err != nil {
			l.Error(err)
			go sys.FormatResult(sys.BREAK, utl.NowString()+err.Error(), nil)
			return
		}

		// 获取系统运行路径
		rate += sys.GET_WORKINGDIR
		go sys.FormatResult(rate, utl.NowString()+"Get the working directory", nil)
		var basepath string
		basepath, err = filepath.Abs("./")
		if err != nil || basepath == "" {
			l.Error(err)
			go sys.FormatResult(sys.BREAK, utl.NowString()+err.Error(), nil)
			return
		}

		// 解析POST.json
		rate += sys.PARSE_JSON
		go sys.FormatResult(rate, utl.NowString()+"Parse the json string", nil)
		sd, arr_lo, err := sys.ParseSysSubmit(jsonstr)
		if err != nil {
			l.Error(err)
			go sys.FormatResult(sys.BREAK, utl.NowString()+err.Error(), nil)
			return
		}

		// do the main process
		rate += sys.MAIN_PROCESS
		go sys.FormatResult(rate, utl.NowString()+"Do the main process", nil)
		if err = sys.Process(sd, arr_lo); err != nil {
			l.Error(err)
			go sys.FormatResult(sys.BREAK, utl.NowString()+err.Error(), nil)
			return
		}

		l.Messagef("Done the main process successfully")
		go sys.FormatResult(sys.END, utl.NowString()+"Done the main process successfully", nil)
		//		OutputJson(w, 0, "Done the main process successfully", nil)
	}
}

// output post message to the front-end
func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := sys.Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}

	w.Write(b)
	l.Debug(b)
}
