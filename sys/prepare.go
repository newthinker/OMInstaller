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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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
	User      string   // 目标机器系统登录用户名
	Pwd       string   // 目标机器系统登录密码
	Basedir   string   // 当前安装包根目录
	Apps      []string // OneMap应用模块
	Services  []string // OneMap服务
	Servers   []string // OneMap server types
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
		fmt.Printf("ERROR: 源文件(%s)不存在\n", srcfile)
		return os.ErrNotExist
	}

	// check dstfile's parent path whether existed
	dir := filepath.Dir(dstfile)
	_, derr := os.Stat(dir)
	if os.IsNotExist(derr) {
		//fmt.Printf("WARN: Destination path(%s) isn't existed and then create it!\n", dstfile)

		if serr = os.MkdirAll(dir, 0755); serr != nil {
			fmt.Printf("ERROR: 生成目标文件夹失败(%s)\n", dir)
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
		fmt.Println("ERROR: 拷贝失败")
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
		return errors.New("ERROR: 输入/输出目录非法")
	}

	if (Exists(src)) != true {
		msg := "ERROR: 输入目录(" + src + ")不存在"
		return errors.New(msg)
	}
	if (Exists(dst)) != true {
		if err := os.Mkdir(dst, 0755); err != nil {
			msg := "ERROR: 目标文件夹(" + dst + ")不存在并且生成失败"
			return errors.New(msg)
		}
	} else {
		if (Exists(dst + "/services")) == true {
			if err := os.RemoveAll(dst + "/services"); err != nil {
				msg := "ERROR: 删除目录(" + dst + "/services)失败"
				return errors.New(msg)
			}
		}
		if (Exists(dst + "/" + om.Container + "/webapps")) == true {
			if err := os.RemoveAll(dst + "/" + om.Container + "/webapps"); err != nil {
				msg := "ERROR: 删除目录(" + dst + "/" + om.Container + "/webapps)失败"
				return errors.New(msg)
			}
		}
	}

	// copy public
	//var cmdline string
	if (Exists(dst + "/arcgis")) != true {
		if err := Copy(src+"/arcgis", dst); err != nil {
			msg := "ERROR: 拷贝目录(" + src + "/arcgis)失败"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/bin")) != true {
		if err := Copy(src+"/bin", dst); err != nil {
			msg := "ERROR: 拷贝目录(" + src + "/bin)失败"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/config")) != true {
		if err := Copy(src+"/config", dst); err != nil {
			msg := "ERROR: 拷贝目录(" + src + "/config)失败"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/java")) != true {
		if err := Copy(src+"/java", dst); err != nil {
			msg := "ERROR: 拷贝目录(" + src + "/java)失败"
			return errors.New(msg)
		}
	}
	if (Exists(dst + "/temp")) != true {
		if err := Copy(src+"/temp", dst); err != nil {
			msg := "ERROR: 拷贝目录(" + src + "/temp)失败"
			return errors.New(msg)
		}
	}

	// copy modules
	if len(om.Apps) > 0 {
		// copy web container
		if err := Copy(src+"/"+om.Container, dst); err != nil {
			msg := "ERROR: 拷贝OneMap服务器容器(" + om.Container + ")失败"
			return errors.New(msg)
		}

		for i := 0; i < len(om.Apps); i++ {
			if err := Copy(src+"/webapps/"+om.Apps[i], dst+"/"+om.Container+"/webapps/"+om.Apps[i]); err != nil {
				msg := "ERROR: 拷贝模块(" + om.Apps[i] + ")失败"
				return errors.New(msg)
			}

			switch om.Apps[i] {
			case "H2memDB":
				if err := Copy(src+"/db", dst); err != nil {
					return errors.New("ERROR: 拷贝数据库目录失败")
				}
			case "GeoShareManager":
				if err := Copy(src+"/example_data", dst); err != nil {
					return errors.New("ERROR: 拷贝示例数据失败")
				}
			}
		}
	}

	// copy services  
	if len(om.Services) > 0 {
		for i := 0; i < len(om.Services); i++ {
			if err := Copy(src+"/services/"+om.Services[i], dst+"/services/"+om.Services[i]); err != nil {
				msg := "ERROR: 拷贝OneMap服务(" + om.Services[i] + ")失败"
				return errors.New(msg)
			}

            // 如果包含H2MemDB模块，且在临时目录中包含分平台数据更新文件，还需要更新分平台数据文件
            if om.Services[i]=="H2CommonMemDB" {
                // 首先拷贝db目录
                if err:=Copy(src+"/db", dst); err!=nil {
                    msg := "ERROR: 拷贝OneMap数据库目录失败!"
                    return errors.New(msg)
                }
                // 更新数据库文件
                srcfile := om.Basedir+"/temp/Manager_Table_Data.sql"
                dstfile := dst+"/db/GeoShareManager/Manager_Table_Data.sql"
                if Exists(om.Basedir+"/temp/Manager_Table_Data.sql")==true {
                    if err:=Copy(srcfile, dstfile); err!=nil {
                        msg := "WARN: 更新分平台数据库文件失败!"
                        return errors.New(msg)
                    }

                    // 删除临时目录
                    os.RemoveAll(om.Basedir+"/temp")
                }
            }
		}
	}

	fmt.Println("MSG: 拷贝OneMap成功")
	return nil
}

func (om *OMPInfo) OMPackage() int {
	var flag bool = false

	// 判断OneMap文件夹是否存在
	var onemap_dir string = "./" + ONEMAP_NAME
	if flag = Exists(onemap_dir); flag != true {
		fmt.Println("WARN: OneMap directory isn't existed!")
		if err := os.Mkdir(onemap_dir, 0755); err != nil {
			fmt.Println("ERROR: Make OneMap directory failed!")
			return 1
		}
	}

	// 判断是否有安装的onemap模块
	if len(om.Apps) <= 0 && len(om.Services) <= 0 {
		fmt.Println("ERROR: no install modules")
		return 2
	}

	return 0
}

// parse the current machine info
func (om *OMPInfo) OMGetInfo(mi *MachineInfo, sm *ServerMapping) error {
	if mi == nil || sm == nil {
		return errors.New("ERROR: 输入参数为空")
	}

	// get attributes
	if mi.Os != "" {
		om.Os = mi.Os
	} else {
		return errors.New("ERROR: 获取输入参数(os)非法")
	}
	if mi.Arch != "" {
		om.Arch = mi.Arch
	} else {
		return errors.New("ERROR: 获取输入参数(arch)非法")
	}
	if mi.Ip != "" {
		om.Ip = mi.Ip
	} else {
		return errors.New("ERROR: 获取输入参数(ip)非法")
	}
	if mi.User != "" {
		om.User = mi.User
	} else {
		return errors.New("ERROR: 获取输入参数(user)非法")
	}
	if mi.Pwd != "" {
		om.Pwd = mi.Pwd
	} else {
		return errors.New("ERROR: 获取输入参数(pwd)非法")
	}
	if mi.Omhome != "" {
		om.OMHome = mi.Omhome
	} else {
		return errors.New("ERROR: 获取输入参数(pwd)非法")
	}

	filename, err := om.OMGetVersion(om.Basedir)
	if err != nil {
		return errors.New("ERROR: 获取OneMap版本信息失败")
	}

	// get the container
	if err = om.OMGetContainer(om.Basedir, filename); err != nil {
		return errors.New("ERROR: 获取OneMap容器信息失败")
	}

	// Get the web app modules name and services name
	for i := 0; i < len(mi.Servers); i++ {
		var srvinfo *ServerInfo = &(mi.Servers[i])
		var srvtype string = srvinfo.XMLName.Local
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
	}

	return nil
}

// list files in the path recursion
func getAllfiles(path string) {
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
		fmt.Printf("filepath.Walk() return %v\n", err)
	}
}

// list sub directory in current path
func getSubDir(path string) ([]string, error) {
	pn := []string{}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("ERROR: 打开输入路径(%s)失败\n", path)
		return pn, err
	}

	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		fmt.Printf("ERROR: 获取输入路径(%s)信息失败\n", path)
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
		msg := "ERROR: 输入路径(" + basedir + ")不存在"
		return base, errors.New(msg)
	}

	// get all sub directory name and search the onemap package directory
	subpath, err := getSubDir(basedir)
	if err != nil {
		msg := "ERROR: 获取子目录失败"
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
		msg := "WARN: OneMap安装包命名非法"
		//return base, errors.New(msg)
        fmt.Println(msg)
	}

	// parse the path name and get the version
	var arr = strings.Split(base, "_")
	om.Version = arr[len(arr)-1]
	if om.Version == "" {
		msg := "ERROR: 获取OneMap版本信息失败"
		return base, errors.New(msg)
	}

	return base, nil
}

// get web container's name
func (om *OMPInfo) OMGetContainer(basedir string, subdirname string) error {
	if flag := Exists(basedir + "/" + subdirname); flag != true {
		msg := "ERROR: 输入目录(" + basedir + ")不存在"
		return errors.New(msg)
	}

	// get all sub directory name and search the onemap package directory
	subpath, err := getSubDir(basedir + "/" + subdirname)
	if err != nil {
		msg := "ERROR: 获取所有子目录失败"
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

		return errors.New("ERROR: 获取web容器路径失败")
	}

	// get the file/path name
	if base == "" {
		return errors.New("ERROR: 非法OneMap安装包")
	}

	// parse the path name and get the version
	om.Container = base

	return nil
}

// remote copy the OneMap package
func (om *OMPInfo) OMRemoteCopy(srcdir string, dstdir string) error {
	// check whether installed sshpass package
	cmd := exec.Command("sshpass", "-V")
	err := cmd.Run()
	if err != nil {
		return errors.New("ERROR: sshpass没有安装")
	}

	// check srcdir is a file or directory
	if flag := Exists(srcdir); flag != true {
		msg := "ERROR: 源文件或目录(" + srcdir + ")不存在"
		return errors.New(msg)
	}

	fi, _ := os.Stat(srcdir)
	if fi.IsDir() {
		cmd = exec.Command("sshpass", "-p", om.Pwd, "scp", "-r", srcdir, om.User+"@"+om.Ip+":"+dstdir)

		fmt.Printf("CMD: sshpass -p %s scp -r %s %s@%s:%s\n", om.Pwd, srcdir, om.User, om.Ip, dstdir)
	} else {
		cmd = exec.Command("sshpass", "-p", om.Pwd, "scp", srcdir, om.User+"@"+om.Ip+":"+dstdir)

		fmt.Printf("CMD: sshpass -p %s scp %s %s@%s:%s\n", om.Pwd, srcdir, om.User, om.Ip, dstdir)
	}
	err = cmd.Run()
	if err != nil {
		return errors.New("ERROR: 执行远程拷贝失败")
	}

	return nil
}

// exec the remote command
func (om *OMPInfo) OMRemoteExec() error {
	// parse the remote command line
	if len(om.Servers) <= 0 {
		return errors.New("ERROR: 没有需要安装的模块")
	}

	// check whether installed sshpass package
	cmd := exec.Command("sshpass", "-V")
	err := cmd.Run()
	if err != nil {
		return errors.New("ERROR: Sshpass没有安装")
	}

	// exec the remote command line
	for i := 0; i < len(om.Servers); i++ {
		cmd = exec.Command("sshpass", "-p", om.Pwd, "ssh", om.User+"@"+om.Ip,
			"/bin/bash", om.OMHome+"/install.sh", om.Servers[i])
		err = cmd.Run()
		if err != nil {
			msg := "ERROR: 安装" + om.Servers[i] + "服务模块失败"
			return errors.New(msg)
		}
	}

	return nil
}

// search the input server params
func (om *OMPInfo) OMInputParams(sc *SysConfig) []string {
	var srvlist = []string{} // save the input server type

	// first check whether install any app or service
	if len(om.Apps) < 1 && len(om.Services) < 1 {
		fmt.Println("WARN: 没有需要安装的模块")
		return srvlist
	}

	// get the relative server type 
	// first the apps
	for i := range om.Apps {
		mdlname := om.Apps[i]

		for j := range sc.FileMap.Containers {
			var container *Container = &(sc.FileMap.Containers[j])
			conname := strings.ToUpper(container.XMLName.Local)
			if conname == strings.ToUpper(TOMCAT) || conname == strings.ToUpper(WEBLOGIC) {
				for k := range container.Modules {
					var module *ModuleMap = &(container.Modules[k])

					if mdlname == module.XMLName.Local {
						for l := range module.ServersMap {
							var srvname = module.ServersMap[l].XMLName.Local

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
			conname := strings.ToUpper(container.XMLName.Local)
			if conname == strings.ToUpper("SERVICES") {
				for k := range container.Modules {
					var module *ModuleMap = &(container.Modules[k])

					if srvname == module.XMLName.Local {
						for l := range module.ServersMap {
							var srvname = module.ServersMap[l].XMLName.Local

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
