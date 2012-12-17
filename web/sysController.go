package web

import (
    "fmt"
    "encoding/json"
	"github.com/newthinker/onemap-installer/sys"
)

type sysController struct {
}

type Result struct {
    Ret int
    Reason string
    Data interface{}
}

// 处理系统配置页面
func (this *sysController) SysAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

    w.Header().Set("content-type", "application/json")

    // 将配置文件解析后传入前端显示
	if r.Method == "GET" {
        str, err := this.SysFormat()
        if err!=nil {
            OutputJson(w, 1, "解析配置文件失败", nil)
            return
        }

        OutputJson(w, 0, "", str)
	} else if r.Method == "POST" {      // 将前端输入参数传入后台解析

    }
}

// 组织上传json字符串
func (this *sysController) SysFormt(sm *ServerMapping, sc *SysConfig) (string, error) {
    
    return "", nil
}

// 解析下载json字符串
func (this *sysController) SysParse(json string) {error} {

    return nil
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
    out := &Result{ret, reason, i}
    b, err := json.Marshal(out)
    if err!= nil {
        return
    }

    w.Write(b)
}
