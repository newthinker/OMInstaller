/*===================================================================
#    FileName: config.go
# Description: 进行安装前配置，包括构造OneMap文件包，分平台参数配置，
#              系统参数配置等；
#      Author: MichaelCho
#       Email: zuow11@gmail.com
#     WebSite: http://www.zone4cho.com
#  CreateTime: 2012.11.06
===================================================================*/
package sys

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
    "errors"
)

const (
	SERVER_MAPPING = "SrvMapping.xml" // 服务器类型映射文件名 
	SYS_INFO       = "SysInfo.xml"    // 服务器参数配置文件
	SYS_CONFIG     = "SysConfig.xml"  // 配置工具配置文件名
	ONEMAP_NAME    = "OneMap"         // OneMap directory name
	TOMCAT         = "Tomcat"         // Tomcat container
	WEBLOGIC       = "Weblogic"       // Weblogic container
)

type OMPInfo struct {
	Version   string   // onemap版本号
	OMHome    string   // OneMap安装目录
	Container string   // OneMap web容器类型
	Os        string   // 目标机器操作系统类型，[windos/linux]
	Arch      string   // 目标机器系统架构，[386/amd64]
	Ip        string   // 目标机器ip地址
	Root      string   // 目标机器系统登录用户名
	Pwd       string   // 目标机器系统登录密码
	Basedir   string   // 当前安装包根目录
	Apps      []string // OneMap应用模块
	Services  []string // OneMap服务
	Servers   []string // OneMap server types

	OM_Group string // OneMap 组名
	OM_User  string // OneMap 系统用户
	OM_PWD   string // OneMap 系统密码

	ORCL_User string    // oracle系统帐号
	ORCL_SID  string    // SID
	DB_User   [6]string // 数据库用户, 6个：system,geoshare_platform, geoshare_portal, geo_coding, geo_portal, geoshare_sub_platform
	DB_PWD    [6]string // 用户密码

	AGS_Home string // AGS home目录
}

// 判断文件或者路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

// Copy file or dirctory
func Copy(srcfile string, dstfile string) error {
	// first check the srcfile whether exist
	fi, serr := os.Stat(srcfile)
	if os.IsNotExist(serr) {
		return os.ErrNotExist
	}

	// check dstfile's parent path whether existed
	dir := filepath.Dir(dstfile)
	_, derr := os.Stat(dir)
	if os.IsNotExist(derr) {

		if serr = os.MkdirAll(dir, 0755); serr != nil {
			return serr
		}
	}

	// check the srcfile is file or directory
	if fi.IsDir() {
		cmd := exec.Command("cp", "-r", srcfile, dstfile)
		serr = cmd.Run()
	} else {
		cmd := exec.Command("cp", srcfile, dstfile)
		serr = cmd.Run()
	}
	// exec the copy comand
	if serr != nil {
		return serr
	}

	return nil
}

// copy OneMap modules from source directory to destination directory
// update: Add input params src and dst for common session. And 
//         not copy public part when copied. [zuow, 2012/11/08]
func (om *OMPInfo) OMCopy(src string, dst string) error {
	// first check the src and dst directory
	if src == "" || dst == "" {
        msg := "Invalid source or destination directory"
		return errors.New(msg)
	}

	if (Exists(src)) != true {
		msg := "Source directory (" + src + ") isn't existed"
		return errors.New(msg)
	}
	if (Exists(dst)) != true {
		if err := os.Mkdir(dst, 0755); err != nil {
			msg := "Destination directory (" + dst + ") isn't existed and create failed"
			return errors.New(msg)
		}
	} else {
		if (Exists(dst + "/services")) == true {
			if err := os.RemoveAll(dst + "/services"); err != nil {
				msg := "Remove directory (" + dst + "/services) failed"
				return errors.New(msg)
			}
		}
		if (Exists(dst + "/" + om.Container + "/webapps")) == true {
			if err := os.RemoveAll(dst + "/" + om.Container + "/webapps"); err != nil {
				msg := "Remove directory (" + dst + "/" + om.Container + "/webapps) failed"
				return errors.New(msg)
			}
		}
	}

	// copy public
	if (Exists(dst + "/install.sh")) != true {
		if err := Copy(src+"/install.sh", dst); err != nil {
			msg := "Copy install bash script failed"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/arcgis")) != true {
		if err := Copy(src+"/arcgis", dst); err != nil {
			msg := "Copy directory (" + src + "/arcgis) failed"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/bin")) != true {
		if err := Copy(src+"/bin", dst); err != nil {
			msg := "Copy directory(" + src + "/bin) failed"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/config")) != true {
		if err := Copy(src+"/config", dst); err != nil {
			msg := "Copy directory(" + src + "/config) failed"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/java")) != true {
		if err := Copy(src+"/java", dst); err != nil {
			msg := "Copy directory(" + src + "/java) failed"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/temp")) != true {
		if err := Copy(src+"/temp", dst); err != nil {
			msg := "Copy directory(" + src + "/temp) failed"
			return errors.New(msg)
		}
	}

	// copy modules
	if len(om.Apps) > 0 {
		// copy web container
		if err := Copy(src+"/"+om.Container, dst); err != nil {
			msg := "Copy OneMap web container (" + om.Container + ") failed"
			return errors.New(msg)
		}

		for i := 0; i < len(om.Apps); i++ {
			if err := Copy(src+"/webapps/"+om.Apps[i], dst+"/"+om.Container+"/webapps/"+om.Apps[i]); err != nil {
				msg := "Copy module (" + om.Apps[i] + ") failed"
				return errors.New(msg)
			}

			switch om.Apps[i] {
			case "H2memDB":
				if err := Copy(src+"/db", dst); err != nil {
					return errors.New("Copy db directory failed")
				}
			case "GeoShareManager":
				if err := Copy(src+"/example_data", dst); err != nil {
					return errors.New("Copy example data directory failed")
				}
			}
		}
	}

	// copy services  
	if len(om.Services) > 0 {
		for i := 0; i < len(om.Services); i++ {
			if err := Copy(src+"/services/"+om.Services[i], dst+"/services/"+om.Services[i]); err != nil {
				msg := "Copy OneMap service (" + om.Services[i] + ") failed"
				return errors.New(msg)
			}

			// 如果包含H2MemDB模块，且在临时目录中包含分平台数据更新文件，还需要更新分平台数据文件
			if om.Services[i] == "H2CommonMemDB" {
				// 首先拷贝db目录
				if err := Copy(src+"/db", dst); err != nil {
					msg := "Copy OneMap db directory failed"
					return errors.New(msg)
				}
				// 更新数据库文件
				srcfile := om.Basedir + "/temp/Manager_Table_Data.sql"
				dstfile := dst + "/db/GeoShareManager/Manager_Table_Data.sql"
				if Exists(om.Basedir+"/temp/Manager_Table_Data.sql") == true {
					if err := Copy(srcfile, dstfile); err != nil {
						msg := "Update subplatform SQL file failed"
						return errors.New(msg)
					}

					// 删除临时目录
					os.RemoveAll(om.Basedir + "/temp")
				}
			}
		}
	}

	l.Message("Copy OneMap files successfully")
	return nil
}

func (om *OMPInfo) OMPackage() int {
	var flag bool = false

	// 判断OneMap文件夹是否存在
	var onemap_dir string = "./" + ONEMAP_NAME
	if flag = Exists(onemap_dir); flag != true {
		l.Warning("OneMap directory isn't existed")
		if err := os.Mkdir(onemap_dir, 0755); err != nil {
			l.Errorf("Make OneMap directory failed")
			return 1
		}
	}

	// 判断是否有安装的onemap模块
	if len(om.Apps) <= 0 && len(om.Services) <= 0 {
		l.Errorf("No install modules")
		return 2
	}

	return 0
}

// parse the current machine info
func (om *OMPInfo) OMGetInfo(mi *MachineInfo, sm *ServerMapping) error {
    l.Message("Begin to init with machine's info")
	if mi == nil || sm == nil {
        msg := "Input MachineInfo and SrvMapping object is nil"
		return errors.New(msg)
	}

	// get attributes
	if mi.Os != "" {
		om.Os = mi.Os
	} else {
        msg := "Get machine's input param(os) is invalid"
		return errors.New(msg)
	}
	if mi.Arch != "" {
		om.Arch = mi.Arch
	} else {
        msg := "Get machine's input param(arch) is invalid"
		return errors.New(msg)
	}
	if mi.Ip != "" {
		om.Ip = mi.Ip
	} else {
        msg := "Get machine's input param(ip) is invalid"
		return errors.New(msg)
	}
	if mi.User != "" {
		om.Root = mi.User
	} else {
        msg := "Get machine's input param(user) is invalid"
		return errors.New(msg)
	}
	if mi.Pwd != "" {
		om.Pwd = mi.Pwd
	} else {
        msg := "Get machine's input param(pwd) is invalid"
		return errors.New(msg)
	}
	if mi.Omhome != "" {
		om.OMHome = mi.Omhome
	} else {
        msg := "Get machine's input param(OneMap Home) is invalid"
		return errors.New(msg)
	}

	filename, err := om.OMGetVersion(om.Basedir)
	if err != nil {
        msg := "Get OneMap version failed"
		return errors.New(msg)
	}

	// get the container
	if err = om.OMGetContainer(om.Basedir, filename); err != nil {
        msg := "Get OneMap container failed"
		return errors.New(msg)
	}

	// Get the web app modules name and services name
	for i := 0; i < len(mi.Servers); i++ {
		var srvinfo *ServerInfo = &(mi.Servers[i])
		var srvtype string = srvinfo.Name
		if srvtype == "" {
			continue
		}

		for j := 0; j < len(sm.Servers); j++ {
			var srv *Server = &(sm.Servers[j])

			if srvtype == srv.XMLName.Local {
				for k := 0; k < len(srv.ModuleList); k++ {
					var mdl Module = srv.ModuleList[k]
					if mdl.XMLName.Local == "app" {
						om.Apps = append(om.Apps, mdl.MdlName)
					} else if mdl.XMLName.Local == "srv" {
						om.Services = append(om.Services, mdl.MdlName)
					}
				}

				// add the server type to the array
				om.Servers = append(om.Servers, srvtype)
			}
		}

		// 获取其它参数信息
		if srvtype == "db" {
			for _, attr := range srvinfo.Attrs {
				if attr.Name == "db_sid" {
					om.ORCL_SID = attr.Value
				} else if attr.Name == "db_user" {
					om.ORCL_User = attr.Value
				} else if attr.Name == "system_user" {
					om.DB_User[0] = attr.Value
				} else if attr.Name == "system_pwd" {
					om.DB_PWD[0] = attr.Value
				} else if attr.Name == "manager_user" {
					om.DB_User[1] = attr.Value
				} else if attr.Name == "manager_pwd" {
					om.DB_PWD[1] = attr.Value
				} else if attr.Name == "portal_user" {
					om.DB_User[2] = attr.Value
				} else if attr.Name == "portal_pwd" {
					om.DB_PWD[2] = attr.Value
				} else if attr.Name == "geocoding_user" {
					om.DB_User[3] = attr.Value
				} else if attr.Name == "geocoding_pwd" {
					om.DB_PWD[3] = attr.Value
				} else if attr.Name == "geoportal_user" {
					om.DB_User[4] = attr.Value
				} else if attr.Name == "geoportal_pwd" {
					om.DB_PWD[4] = attr.Value
				} else if attr.Name == "sub_user" {
					om.DB_User[5] = attr.Value
				} else if attr.Name == "sub_pwd" {
					om.DB_PWD[5] = attr.Value
				}
			}
		} else if srvtype == "gis" {
			for k := 0; k < len(srvinfo.Attrs); k++ {
				attr := &(srvinfo.Attrs[k])
				if attr.Name == "ags_log_path" {
					om.AGS_Home = attr.Value
					/// AGS default log path
					attr.Value += "/server/user/log"
				}
			}
		}
	}

	/// default params
	om.OM_Group = "esri"
	om.OM_User = "esri"
	om.OM_PWD = "esri1234"

	return nil
}

// list files in the path recursion
func GetAllfiles(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		return nil
	})

	if err != nil {
		l.Debugf("filepath.Walk() return %v", err)
	}
}

// list sub directory in current path
func GetSubDir(path string) ([]string, error) {
	pn := []string{}

	f, err := os.Open(path)
	if err != nil {
		l.Errorf("Open input path(%s) failed", path)
		return pn, err
	}

	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		l.Errorf("Read input path(%s) failed", path)
		return pn, err
	}

	for _, fileinfo := range list {
		if fileinfo == nil {
			continue
		}
		if fileinfo.IsDir() {
			var pathname string = fileinfo.Name()
			if pathname != "" {
				pn = append(pn, pathname)
			}
		}
	}

	return pn, err
}

// get onemap's version
func (om *OMPInfo) OMGetVersion(basedir string) (string, error) {
	var base string // the onemap package directory
	if flag := Exists(basedir); flag != true {
		msg := "Input directory(" + basedir + ") isn't existed"
		return base, errors.New(msg)
	}

	// get all sub directory name and search the onemap package directory
	subpath, err := GetSubDir(basedir)
	if err != nil {
		msg := "Get all the sub directory failed"
		return base, errors.New(msg)
	}

	for i := range subpath {
		filename := path.Base(subpath[i]) // get the base filename
		filename = strings.ToUpper(filename)
		if strings.Index(filename, strings.ToUpper(ONEMAP_NAME)) < 0 {
			continue
		}

		if strings.Index(filename, "_V") > 0 {
			base = subpath[i]
			break
		}
	}

	// get the file/path name
	if base == "" {
		msg := "Invalid OneMap package"
		return base, errors.New(msg)
	}

	// parse the path name and get the version
	var arr = strings.Split(base, "_")
	om.Version = arr[len(arr)-1]
	if om.Version == "" {
		msg := "Invalid package name(" + base + ") and have no version information"
		return base, errors.New(msg)
	}

	return base, nil
}

// get web container's name
func (om *OMPInfo) OMGetContainer(basedir string, subdirname string) error {
	if flag := Exists(basedir + "/" + subdirname); flag != true {
		msg := "Input directory(" + basedir + ") isn't existed"
		return errors.New(msg)
	}

	// get all sub directory name and search the onemap package directory
	subpath, err := GetSubDir(basedir + "/" + subdirname)
	if err != nil {
		msg := "Get all the sub directory failed"
		return errors.New(msg)
	}

	var base string // the onemap package directory
	for i := range subpath {
		base = path.Base(subpath[i]) // get the base filename
		if ((strings.Index(strings.ToUpper(base), strings.ToUpper(TOMCAT))) >= 0) ||
			((strings.Index(strings.ToUpper(base), strings.ToUpper(WEBLOGIC))) > 0) {
			break
		} else {
			continue
		}

		return errors.New("Get container's path failed")
	}

	// get the file/path name
	if base == "" {
        msg := "Invalid OneMap package"
		return errors.New(msg)
	}

	// parse the path name and get the version
	om.Container = base

	return nil
}

// remote copy the OneMap package
func (om *OMPInfo) OMRemoteCopy(srcdir string, dstdir string) error {
	// check whether installed sshpass package
	cmd := exec.Command(om.Basedir+"/sshpass/bin/sshpass", "-V")
	err := cmd.Run()
	if err != nil {
		return errors.New("sshpass isn't installed")
	}

	// check srcdir is a file or directory
	if flag := Exists(srcdir); flag != true {
		msg := "Source file or directory " + srcdir + " isn't existed"
		return errors.New(msg)
	}

	fi, _ := os.Stat(srcdir)
	if fi.IsDir() {
		cmd = exec.Command(om.Basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "scp", "-r", srcdir, om.Root+"@"+om.Ip+":"+dstdir)

		l.Debugf("sshpass -p %s scp -r %s %s@%s:%s", om.Pwd, srcdir, om.Root, om.Ip, dstdir)
	} else {
		cmd = exec.Command(om.Basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "scp", srcdir, om.Root+"@"+om.Ip+":"+dstdir)

		l.Debugf("sshpass -p %s scp %s %s@%s:%s", om.Pwd, srcdir, om.Root, om.Ip, dstdir)
	}
	err = cmd.Run()
	if err != nil {
		return errors.New("Exec remote copy command failed")
	}

	return nil
}

// exec the remote command
func (om *OMPInfo) OMRemoteExec() error {
	// parse the remote command line
	if len(om.Servers) <= 0 {
        msg := "No install modules"
		return errors.New(msg)
	}

	// check whether installed sshpass package
	cmd := exec.Command(om.Basedir+"/sshpass/bin/sshpass", "-V")
	err := cmd.Run()
	if err != nil {
        msg := "Sshpass isn't installed"
		return errors.New(msg)
	}

	// service flag
	flag_ma := true  // monitoragent service
	flag_h2 := false // h2memdb service
	flag_mq := false // activemq service
	flag_om := false // onemap service

	// exec the remote command line to install the OneMap
	for i := 0; i < len(om.Servers); i++ {
		cmd = exec.Command(om.Basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
			"/bin/bash", om.OMHome+"/install.sh", om.Servers[i])
		l.Debugf("sshpass -p %s ssh %s@%s /bin/bash %s/install.sh %s", om.Pwd, om.Root, om.Ip,
			om.OMHome, om.Servers[i])
		err = cmd.Run()
		if err != nil {
			msg := "Install " + om.Servers[i] + " module failed"
			return errors.New(msg)
		}

		if (flag_ma == false) && ((om.Servers[i] == "gis") || (om.Servers[i] == "web") || (om.Servers[i] == "token")) {
			flag_om = true
		}
		if (flag_h2 == false) && (om.Servers[i] == "main") {
			flag_h2 = true
			flag_om = true
		}
		if (flag_mq == false) && (om.Servers[i] == "msg") {
			flag_mq = true
		}

		if flag_ma == true {
			cmd = exec.Command("nohup", om.Basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/monitoragent", "start", ">/dev/null", "2>&1", "&")
			err = cmd.Run()
			if err != nil {
				l.Errorf("Start up service monitoragent failed")
			}

			flag_ma = false // only run one time
		}
		if flag_h2 == true {
			cmd = exec.Command("nohup", om.Basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/h2memdb", "start", ">/dev/null", "2>&1", "&")
			err = cmd.Run()
			if err != nil {
				l.Errorf("Start up service h2memdb failed")
			}

			flag_h2 = false
		}
		if flag_mq == true {
			cmd = exec.Command(om.Basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/activemq", "start")
			err = cmd.Run()
			if err != nil {
				l.Errorf("Start up service activemq failed")
			}

			flag_mq = false
		}
		if flag_om == true {
			cmd = exec.Command("nohup", om.Basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/onemap", "start", ">/dev/null", "2>&1", "&")
			err = cmd.Run()
			if err != nil {
				l.Errorf("Start up service onemap failed")
			}

			flag_om = false
		}
	}

	return nil
}

// search the input server params
func (om *OMPInfo) OMInputParams(sc *SysConfig) []string {
	var srvlist = []string{} // save the input server type

	// first check whether install any app or service
	if len(om.Apps) < 1 && len(om.Services) < 1 {
		l.Warning("No installed modules")
		return srvlist
	}

	// get the relative server type 
	// first the apps
	for i := range om.Apps {
		mdlname := om.Apps[i]

		for j := range sc.FileMap.Containers {
			var container *Container = &(sc.FileMap.Containers[j])
			conname := strings.ToUpper(container.Name)
			if conname == strings.ToUpper(TOMCAT) || conname == strings.ToUpper(WEBLOGIC) {
				for k := range container.Modules {
					var module *ModuleMap = &(container.Modules[k])

					if mdlname == module.Name {
						for l := range module.ServersMap {
							var srvname = module.ServersMap[l].Name

							var flag bool = true

							for m := range srvlist {
								if srvname == srvlist[m] {
									flag = false
									break
								}
							}

							if flag { // remove the repeat server type
								srvlist = append(srvlist, srvname)
							}
						}
					}
				}
			}
		}
	}

	for i := range om.Services {
		srvname := om.Services[i]

		for j := range sc.FileMap.Containers {
			var container *Container = &(sc.FileMap.Containers[j])
			conname := strings.ToUpper(container.Name)
			if conname == strings.ToUpper("SERVICES") {
				for k := range container.Modules {
					var module *ModuleMap = &(container.Modules[k])

					if srvname == module.Name {
						for l := range module.ServersMap {
							var srvname = module.ServersMap[l].Name

							var flag bool = true

							for m := range srvlist {
								if srvname == srvlist[m] {
									flag = false
									break
								}
							}

							if flag { // remove the repeat server type
								srvlist = append(srvlist, srvname)
							}
						}
					}
				}
			}
		}
	}

	return srvlist
}
