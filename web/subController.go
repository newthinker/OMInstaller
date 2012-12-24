package web

import (
	"fmt"
	"github.com/newthinker/onemap-installer/sys"
	"html/template"
	"log"
	"net/http"
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
			log.Println(err)
		} else { // 解析用户选择menu并进入下个页面
			fmt.Println(r.Form["selValues"])

			// 进行分平台配置
			var base string
			base, err = filepath.Abs(".") // 获取系统当前路径
			if err != nil {
				log.Println(err)
				return
			}

			sub := &sys.SubPlatform{SqlFile: base + "/OneMap/db/GeoShareManager/Manager_Table_Data.sql"}
			sub.MenuMap = make(map[string]string)
			sub.RelMap = make(map[string]string)

			base = (r.Form["selValues"])[0]
			fmt.Println(base)
			sub.SelID = strings.Split(base, ",")
			fmt.Println(sub.SelID)

			// 解析sql文件并初始化menuMap和relMap
			if err := sub.SPParseSQLFile(); err != nil {
				fmt.Println("ERROR: Parse sql file failed!")
				return
			}

			// 更新sql文件
			if err = sub.SPUpdateSql(); err != nil {
				fmt.Println("ERROR: Update sql file failed!")
				return
			}

			// 进入系统参数配置页面
			http.Redirect(w, r, "/sysconfig", http.StatusFound)
		}
	}
}

// 查找GeoShareManager/Manager_Table_Data.sql
func GetSqlFile(basedir string) error {
    var base string

    if flag := sys.Exists(basedir); flag!=true {
        return errors.New("ERROR: 输入目录不存在")
    }

    subpath, err := getSubDir(basedir)
    if err!=nil {
        return errors.New("ERROR: 获取子目录失败")
    }

    for i:=range subpath {
        filename := path.Base(subpath[i])
        filename = strings.ToUpper(filename)
        if strings.Index(filename, strings.ToUpper(sys.ONEMAP_NAME))<0 {
            continue
        }

        if strings.Index(filename, "_") >0 {
            base = subpath[i]
            break
        }
    }

    filename = base + "/db/GeoShareManager/Manager_Table_Data.sql"
    if flag := sys.Exists(filename); flag!=true {
        return errors.New("ERROR: SQL文件不存在")
    }

	return nil
}
