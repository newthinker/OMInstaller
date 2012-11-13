/*===================================================================
#    FileName: config.go
# Description: 进行安装前配置，包括构造OneMap文件包，分平台参数配置，
#              系统参数配置等；
#      Author: MichaelCho
#       Email: zuow11@gmail.com
#     WebSite: http://www.zone4cho.com
#  CreateTime: 2012.11.06
===================================================================*/
package onemap

import (
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

// copy OneMap modules from source directory to destination directory
// update: Add input params src and dst for common session. And 
//         not copy public part when copied. [zuow, 2012/11/08]
func (om *OMPInfo) OMCopy(src string, dst string) int {
	// first check the src and dst directory
	if src == "" || dst == "" {
		fmt.Println("ERROR: Invalid source or destination directory!")
		return 1
	}

	if (Exists(src)) != true {
		fmt.Printf("ERROR: Source directory(%s) isn't existed!", src)
		return 1
	}
	if (Exists(dst)) != true {
		if err := os.Mkdir(dst, 0755); err != nil {
			fmt.Printf("ERROR: Destination directory(%s) isn't existed and create failed!", dst)
			return 1
		}
	} else {
		if (Exists(dst + "/services")) == true {
			if err := os.RemoveAll(dst + "/services"); err != nil {
				fmt.Printf("ERROR: Remove directory %s failed!", dst+"/services")
				return 1
			}
		}
		if (Exists(dst + "/" + om.Container + "/webapps")) == true {
			if err := os.RemoveAll(dst + "/" + om.Container + "/webapps"); err != nil {
				fmt.Printf("ERROR: Remove directory %s failed!", dst+"/"+om.Container+"/webapps")
				return 1
			}
		}
	}

	// copy public
	var cmdline string
	if (Exists(dst + "/arcgis")) != true {
		cmdline = "cp -r " + src + "/arcgis/ " + dst
	}
	if (Exists(dst + "/bin")) != true {
		if cmdline != "" {
			cmdline = cmdline + ";"
		}
		cmdline = cmdline + "cp -r " + src + "/bin " + dst
	}
	if (Exists(dst + "/config")) != true {
		if cmdline != "" {
			cmdline = cmdline + ";"
		}
		cmdline = cmdline + "cp -r " + src + "/config " + dst
	}
	if (Exists(dst + "/java")) != true {
		if cmdline != "" {
			cmdline = cmdline + ";"
		}
		cmdline = cmdline + "cp -r " + src + "/java " + dst
	}
	if (Exists(dst + "/temp")) != true {
		if cmdline != "" {
			cmdline = cmdline + ";"
		}
		cmdline = cmdline + "cp -r " + src + "/temp " + dst
	}

	cmd := exec.Command(cmdline)
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR: Copy OneMap public part failed!")
		return 2
	}

	// copy modules
	if len(om.Apps) > 0 {
		// copy web container
		cmdline = "cp -r " + src + "/" + om.Container + " " + dst
		cmd = exec.Command(cmdline)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("ERROR: Copy OneMap web container %s failed!\n", om.Container)
			return 2
		}

		for i := 0; i < len(om.Apps); i++ {
			cmdline = "cp -r " + src + "/" + om.Container + "/webapps/" + om.Apps[i] + " " + dst + "/" + om.Container + "/webapps"

			switch om.Apps[i] {
			case "H2memDB":
				cmdline = cmdline + "; cp -r " + src + "/db " + dst
			case "GeoShareManager":
				cmdline = cmdline + "; cp -r " + src + "/example_data " + dst
			}

			cmd = exec.Command(cmdline)
			err = cmd.Run()
			if err != nil {
				fmt.Printf("WARN: Copy OneMap %s module failed!\n", om.Apps[i])
			}
		}
	}

	// copy services  
	if len(om.Services) > 0 {
		// create services directory
		if err = os.Mkdir(ONEMAP_NAME+"/services", 0755); err != nil {
			fmt.Println("ERROR: Create services directory failed!")
			return 3
		}

		for i := 0; i < len(om.Services); i++ {
			cmdline = "cp -r " + ONEMAP_NAME + om.Version + "/services/" + om.Services[i] + " " + ONEMAP_NAME + "/services"

			cmd = exec.Command(cmdline)
			err = cmd.Run()
			if err != nil {
				fmt.Printf("ERROR: Copy OneMap %s module failed!\n", om.Services[i])
				return 3
			}
		}
	}

	fmt.Println("MSG: Copy OneMap files successfully!")
	return 0
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
func (om *OMPInfo) OMGetInfo(mi *MachineInfo, sm *ServerMapping) int {
	if mi == nil || sm == nil {
		fmt.Println("ERROR: Input MachineInfo and SrvMapping object is nil!")
		return 1
	}

	// get attributes
	if mi.Os != "" {
		om.Os = mi.Os
	} else {
		fmt.Println("ERROR: Get machine's input param(os) is invalid!")
		return 2
	}
	if mi.Arch != "" {
		om.Arch = mi.Arch
	} else {
		fmt.Println("ERROR: Get machine's input param(arch) is invalid!")
		return 2
	}
	if mi.Ip != "" {
		om.Ip = mi.Ip
	} else {
		fmt.Println("ERROR: Get machine's input param(ip) is invalid!")
		return 2
	}
	if mi.User != "" {
		om.User = mi.Ip
	} else {
		fmt.Println("ERROR: Get machine's input param(user) is invalid!")
		return 2
	}
	if mi.Pwd != "" {
		om.Pwd = mi.Pwd
	} else {
		fmt.Println("ERROR: Get machine's input param(pwd) is invalid!")
		return 2
	}
	if mi.Omhome != "" {
		om.OMHome = mi.Omhome
	} else {
		fmt.Println("ERROR: Get machine's intput param(pwd) is invalid!")
		return 2
	}
	// get the version
	ret := om.OMGetVersion(".")
	if ret != 0 {
		fmt.Println("ERROR: Get OneMap version failed!")
		return 2
	}
	// get the container
	ret = om.OMGetContainer(".")
	if ret != 0 {
		fmt.Println("ERROR: Get OneMap container failed!")
		return 2
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
					var mdl *Module = &(srv.ModuleList[k])
					if mdl.XMLName.Local == "app" {
						om.Apps = append(om.Apps, mdl.ModuleName)
					} else if mdl.XMLName.Local == "srv" {
						om.Services = append(om.Services, mdl.ModuleName)
					}
				}

				// add the server type to the array
				om.Servers = append(om.Servers, srvtype)
			}
		}
	}

	return 0
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
		fmt.Printf("ERROR: Open input path(%s) failed!\n", path)
		return pn, err
	}

	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		fmt.Printf("ERROR: Read input path(%s) failed!\n", path)
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
func (om *OMPInfo) OMGetVersion(basedir string) int {
	if flag := Exists(basedir); flag != true {
		fmt.Printf("ERROR: Input directory(%s) isn't existed!\n", basedir)
		return 1
	}

	// get all sub directory name and search the onemap package directory
	subpath, err := getSubDir(basedir)
	if err != nil {
		fmt.Printf("ERROR: Get all the sub directory failed!\n", basedir)
		return 2
	}

	var base string // the onemap package directory
	for i := range subpath {
		var filename string = path.Base(subpath[i]) // get the base filename
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
		fmt.Println("ERROR: Invalid OneMap package!")
		return 3
	}

	// parse the path name and get the version
	var arr = strings.Split(base, "_")
	om.Version = arr[len(arr)-1]
	if om.Version == "" {
		fmt.Printf("ERROR: Invalid package name(%s) and have no version information!\n", base)
		return 4
	}

	return 0
}

// get web container's name
func (om *OMPInfo) OMGetContainer(basedir string) int {
	if flag := Exists(basedir); flag != true {
		fmt.Printf("ERROR: Input directory(%s) isn't existed!\n", basedir)
		return 1
	}

	// get all sub directory name and search the onemap package directory
	subpath, err := getSubDir(basedir)
	if err != nil {
		fmt.Printf("ERROR: Get all the sub directory failed!\n", basedir)
		return 2
	}

	var base string // the onemap package directory
	for i := range subpath {
		var filename string = path.Base(subpath[i]) // get the base filename
		filename = strings.ToUpper(filename)
		if (strings.Index(filename, strings.ToUpper(TOMCAT)) < 0) && (strings.Index(filename, strings.ToUpper(WEBLOGIC)) < 0) {
			continue
		}
	}

	// get the file/path name
	if base == "" {
		fmt.Println("ERROR: Invalid OneMap package!")
		return 3
	}

	// parse the path name and get the version
	om.Container = base

	return 0
}

// remote copy the OneMap package
func (om *OMPInfo) OMRemoteCopy(srcdir string, dstdir string) int {
	// check whether installed sshpass package
	cmdline := "sshpass -V"
	cmd := exec.Command(cmdline)
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR: sshpass isn't installed!")
		return 1
	}

	// check srcdir is a file or directory
	if flag := Exists(srcdir); flag != true {
		fmt.Printf("ERROR: Source file or directory(%s) isn't existed!\n", srcdir)
		return 2
	}
	fi, _ := os.Stat(srcdir)
	if fi.IsDir() {
		cmdline = "sshpass -p " + om.Pwd + " scp -r " + srcdir + om.User + "@" + om.Ip + ":" + dstdir
	} else {
		cmdline = "sshpass -p " + om.Pwd + " scp " + srcdir + om.User + "@" + om.Ip + ":" + dstdir
	}
	cmd = exec.Command(cmdline)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: Exec remote copy command(%s) failed!\n", cmdline)
		return 3
	}

	return 0
}

// exec the remote command
func (om *OMPInfo) OMRemoteExec(rcmd string) int {
	//check the cmdline
	if rcmd == "" {
		fmt.Println("ERROR: The cmd is null!")
		return 1
	}

	// check whether installed sshpass package
	cmdline := "sshpass -V"
	cmd := exec.Command(cmdline)
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR: sshpass isn't installed!")
		return 2
	}

	cmdline = "sshpass -p " + om.Pwd + " ssh " + om.User + "@" + om.Ip + " " + rcmd
	cmd = exec.Command(cmdline)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("ERROR: Exec remote command(%s) failed!", cmdline)
		return 3
	}

	return 0
}
