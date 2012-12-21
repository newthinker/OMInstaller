package web

import (
	//	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/newthinker/onemap-installer/sys"
	"net/http"
	"path/filepath"
)

type sysController struct {
	sm sys.ServerMapping
	sc sys.SysConfig
}

type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

// Init each variable
func (this *sysController) Init() error {
	// 解析本地的配置文件
	basedir, err := filepath.Abs("./")
	if err != nil || basedir == "" {
		//		OutputJson(w, 1, "获取当前路径失败!", nil)
		return err
	}
	fmt.Println("Current directory:" + basedir)

	sm, err1 := sys.OpenSMConfig(basedir)
	sc, err2 := sys.OpenSCConfig(basedir)
	if err1 != nil || err2 != nil {
		//		OutputJson(w, 2, "解析系统配置文件失败", nil)
		return errors.New("解析系统配置文件失败")
	}

	this.sm = sm
	this.sc = sc

	return nil
}

// 处理系统配置页面
func (this *sysController) SysAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

	w.Header().Set("content-type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("content-type", "text/json;charset=utf-8")

	// 将配置文件解析后传入前端显示
	if r.Method == "GET" {
		sysmap, err := this.SysFormat()
		if err != nil {
			OutputJson(w, 3, "格式化系统参数失败", nil)
			return
		}

		//		for k1, v1 := range sysmap {
		//			fmt.Println(k1)

		//			//			lstmap := list.New()
		//			lstmap := v1.(*list.List)
		//			for e := lstmap.Front(); e != nil; e = e.Next() {
		//				map2 := (e.Value).(map[string]interface{})
		//				for k2, v2 := range map2 {
		//					fmt.Println(k2)
		//					fmt.Println(v2)
		//				}
		//			}

		//			fmt.Println(v1)
		//		}

		//		sysmap := make(map[string]interface{})
		//		sysmap["1"] = "a"
		//		sysmap["2"] = "b"
		//		lstmap := make([](map[string]string))
		//		map1 := map[string]string{"5": "e"}
		//		lstmap = append(lstmap, map1)
		//		map2 := map[string]string{"6": "f"}
		//		lsttest.Maps = append(lsttest.Maps, map2)
		//		map3 := map[string]string{"7": "g"}
		//		lsttest.Maps = append(lsttest.Maps, map3)

		//		sysmap["3"] = lsttest

		OutputJson(w, 0, "", sysmap)
	} else if r.Method == "POST" { // 将前端输入参数传入后台解析
		err := r.ParseForm()
		if err != nil {
			OutputJson(w, 1, err.Error(), nil)
			return
		}

		input := (r.Form["input"])[0]
		//		fmt.Println(input)

		// 获取json数据流
		jsonstr := make(map[string]interface{})
		err = json.Unmarshal([]byte(input), &jsonstr)
		if err != nil {
			OutputJson(w, 2, err.Error(), nil)
			return
		}

		// 获取系统运行路径
		var basepath string
		basepath, err = filepath.Abs("./")
		if err != nil || basepath == "" {
			OutputJson(w, 3, err.Error(), nil)
			return
		}

		err = sys.ParseSysSubmit(jsonstr, basepath, &(this.sc), &(this.sm))
		if err != nil {
			OutputJson(w, 4, err.Error(), nil)
			return
		}

		/// 进行分布式安装 
		//		err = sys.Distribute(basepath, &(this.sc), &(this.sm))
		//		if err != nil {
		//			OutputJson(w, 5, "分布式安装失败", nil)
		//			return
		//		}

		//////////////////////////////////////////////
		OutputJson(w, 0, "分布式安装成功", nil)
	}
}

// 组织上传json字符串
func (this *sysController) SysFormat() (map[string]interface{}, error) {
	//	if sm == nil || sc == nil {
	//		return nil, errors.New("Init error!")
	//	}
	//	fmt.Printf("SM:%s\n", this.sm)
	//	fmt.Printf("SC:%s\n", this.sc)

	//	map1 := this.sm.FormatSrvMapping()
	//	map2 := this.sc.FormatSysConfig()

	srvsmdl := sys.FormatSrvMapping(this.sm)
	srvsparam := sys.FormatSysConfig(this.sc)

	result := make(map[string]interface{})

	if len(srvsmdl.Server_modules) > 0 && len(srvsparam.Server_params) > 0 {
		result["Server_modules"] = srvsmdl.Server_modules
		result["Server_params"] = srvsparam.Server_params
	}

	// fmt.Println(map1)
	//	fmt.Println(map2)

	//	if len(map1) <= 0 || len(map2) <= 0 {
	//		return nil, errors.New("没有输入配置参数")
	//	}

	//	for k, v := range map2 {
	//		map1[k] = v
	//	}

	return result, nil
}

func OutputJson(w http.ResponseWriter, ret int, reason string, i interface{}) {
	out := Result{ret, reason, i}
	b, err := json.Marshal(out)
	if err != nil {
		return
	}

	w.Write(b)

	//	fmt.Printf("%s\n", b)
}
