package web

import (
	"encoding/json"
	"errors"
	"github.com/newthinker/onemap-installer/sys"
	"net/http"
	"path/filepath"
)

type sysController struct {
	//	sm sys.ServerMapping
	//	sc sys.SysConfig
}

type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

// Init each variable
func (this *sysController) Init() error {
	/*    if err:=sys.Init();err!=nil {
	        l.Error(err)
	        return err
	    }
		// 解析本地的配置文件
		basedir, err := filepath.Abs("./")
		if err != nil || basedir == "" {
			l.Error(err)
			return err
		}
		l.Debugf("Current dir is: %s", basedir)

		sm, err1 := sys.OpenSMConfig(basedir)
		sc, err2 := sys.OpenSCConfig(basedir)
		if err1 != nil || err2 != nil {
			l.Error(errors.New("Parse system config files failed"))
			return errors.New("Parse system config files failed")
		}

		this.sm = sm
		this.sc = sc
	*/
	return nil
}

// 处理系统配置页面
func (this *sysController) SysAction(w http.ResponseWriter, r *http.Request) {
	l.Messagef("method:", r.Method)

	w.Header().Set("content-type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("content-type", "text/json;charset=utf-8")

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
		err := r.ParseForm()
		if err != nil {
			l.Error(err)
			OutputJson(w, 1, err.Error(), nil)
			return
		}

		input := (r.Form["input"])[0]
		l.Debugf("Input params:%s", input)

		// 获取json数据流
		jsonstr := make(map[string]interface{})
		err = json.Unmarshal([]byte(input), &jsonstr)
		if err != nil {
			l.Error(err)
			OutputJson(w, 2, err.Error(), nil)
			return
		}

		// 获取系统运行路径
		var basepath string
		basepath, err = filepath.Abs("./")
		if err != nil || basepath == "" {
			l.Error(err)
			OutputJson(w, 3, err.Error(), nil)
			return
		}

		// 解析POST.json
        sd, arr_lo, err := sys.ParseSysSubmit(jsonstr)
		if err != nil {
			l.Error(err)
			OutputJson(w, 4, err.Error(), nil)
			return
		}

        // 进行分布式安装
        if err = sys.Process(sd, arr_lo); err!=nil {
            l.Error(err)
            OutputJson(w, 5, err.Error(), nil)
            return
        }

		l.Messagef("Distribute installing successfully")
		OutputJson(w, 0, "分布式安装成功", nil)
	}
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}

	w.Write(b)
	l.Debug(b)
}
