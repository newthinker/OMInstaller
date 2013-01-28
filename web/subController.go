package web

import (
	"errors"
	"github.com/newthinker/onemap-installer/sys"
	"github.com/newthinker/onemap-installer/utl"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type subController struct {
}

// 处理用户菜单选择操作
func (this *subController) SelectAction(w http.ResponseWriter, r *http.Request) {
	l.Messagef("method:", r.Method)

	if r.Method == "GET" {
		t, err := template.ParseFiles("template/subconfig.html")
		if err != nil {
			l.Error(err)
		}
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		err := r.ParseForm() // 解析URL传递的参数

		// 如果分平台参数解析有问题，报告错误并返回
		if err != nil {
			l.Error(err)
			OutputJson(w, 1, err.Error(), nil)
		} else { // 解析用户选择menu并进入下个页面
			//			fmt.Println(r.Form["selectValues"])
			l.Messagef("Subplatform select nodes:%s", r.Form["selectValues"])

			// 进行分平台配置
			var base string
			base, err = filepath.Abs("./") // 获取系统当前路径
			l.Debugf("base path:%s", base)
			if err != nil || base == "" {
				l.Error(err)
				OutputJson(w, 2, err.Error(), nil)
				return
			}

			var sqlfile string
			sqlfile, err = GetSqlFile(base)
			if err != nil || sqlfile == "" {
				l.Error(err)
				OutputJson(w, 3, err.Error(), nil)
				return
			}
			sub := &sys.SubPlatform{SqlFile: sqlfile}
			sub.MenuMap = make(map[string]string)
			sub.RelMap = make(map[string]string)

			base = (r.Form["selectValues"])[0]
			l.Debugf("Subplatform select nodes:%s", base)
			sub.SelID = strings.Split(base, ",")

			// 解析sql文件并初始化menuMap和relMap
			if err := sub.SPParseSQLFile(); err != nil {
				l.Error(errors.New("Parse SQL file failed"))
				OutputJson(w, 4, "ERROR: 解析SQL文件失败", nil)
				return
			}

			// 更新sql文件
			if err = sub.SPUpdateSql(); err != nil {
				l.Error(errors.New("Update SQL file failed"))
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

	if flag := utl.Exists(basedir); flag != true {
		return filename, errors.New("Input directory isn't existed")
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
	l.Debugf("filename:%s", filename)
	if flag := utl.Exists(filename); flag != true {
		return filename, errors.New("SQL file isn't existed")
	}

	// 生成一个临时文件夹
	tempdir := basedir + "/temp"
	if flag := utl.Exists(tempdir); flag == true {
		err = os.RemoveAll(tempdir)
		if err != nil {
			return filename, errors.New("Delete temp directory failed")
		}
	}
	if err = os.Mkdir(basedir+"/temp", 0755); err != nil {
		return filename, err
	}

	if err = utl.Copy(filename, basedir+"/temp"); err != nil {
		return filename, err
	}

	filename = basedir + "/temp/Manager_Table_Data.sql"

	return filename, nil
}
