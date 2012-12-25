package web

import (
	"errors"
	"fmt"
	"github.com/newthinker/onemap-installer/sys"
	"html/template"
	"log"
	"net/http"
	"os"
	//	"path"
	"path/filepath"
	"strings"
)

type subController struct {
}

// 处理用户菜单选择操作
func (this *subController) SelectAction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/subconfig.html")
		if err != nil {
			log.Println(err)
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		err := r.ParseForm() // 解析URL传递的参数

		// 如果分平台参数解析有问题，报告错误并返回
		if err != nil {
			OutputJson(w, 1, err.Error(), nil)
		} else { // 解析用户选择menu并进入下个页面
			fmt.Println(r.Form["selectValues"])

			// 进行分平台配置
			var base string
			base, err = filepath.Abs("./") // 获取系统当前路径
			fmt.Println("base:" + base)
			if err != nil || base == "" {
				OutputJson(w, 2, err.Error(), nil)
				return
			}

			var sqlfile string
			sqlfile, err = GetSqlFile(base)
			if err != nil || sqlfile == "" {
				OutputJson(w, 3, err.Error(), nil)
				return
			}
			sub := &sys.SubPlatform{SqlFile: sqlfile}
			sub.MenuMap = make(map[string]string)
			sub.RelMap = make(map[string]string)

			base = (r.Form["selectValues"])[0]
			fmt.Println(base)
			sub.SelID = strings.Split(base, ",")
			fmt.Println(sub.SelID)

			// 解析sql文件并初始化menuMap和relMap
			if err := sub.SPParseSQLFile(); err != nil {
				OutputJson(w, 4, "ERROR: 解析SQL文件失败", nil)
				return
			}

			// 更新sql文件
			if err = sub.SPUpdateSql(); err != nil {
				OutputJson(w, 5, "ERROR: 更新SQL文件失败", nil)
				return
			}

			// 进入系统参数配置页面
			http.Redirect(w, r, "/sysconfig", http.StatusFound)
		}
	}
}

// 查找GeoShareManager/Manager_Table_Data.sql
func GetSqlFile(basedir string) (string, error) {
	var filename string
	var err error

	if flag := sys.Exists(basedir); flag != true {
		return filename, errors.New("ERROR: 输入目录不存在")
	}

	//	subpath, err := sys.GetSubDir(basedir)
	//	if err != nil || subpath==""{
	//		return filename, errors.New("ERROR: 获取子目录失败")
	//	}

	//	for i := range subpath {
	//		filename := path.Base(subpath[i])
	//		filename = strings.ToUpper(filename)
	//		if strings.Index(filename, strings.ToUpper(sys.ONEMAP_NAME)) < 0 {
	//			continue
	//		}

	//		if strings.Index(filename, "_V") > 0 {
	//			filename = subpath[i]
	//			break
	//		}
	//	}

	filename = basedir + "/OneMap_Linux_V2.0/db/GeoShareManager/Manager_Table_Data.sql"
	fmt.Println("filename:" + filename)
	if flag := sys.Exists(filename); flag != true {
		return filename, errors.New("ERROR: SQL文件不存在")
	}

	// 生成一个临时文件夹
	if err = os.Mkdir(basedir+"/temp", 0755); err != nil {
		return filename, err
	}

	if err = sys.Copy(filename, basedir+"/temp"); err != nil {
		return filename, err
	}

	filename = basedir + "/temp/Manager_Table_Data.sql"

	return filename, nil
}
