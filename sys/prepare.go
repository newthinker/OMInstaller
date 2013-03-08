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
	"github.com/newthinker/onemap-installer/utl"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	SERVER_MAPPING = "SrvMapping.xml" // 服务器类型映射文件名 
	//	SYS_INFO       = "SysInfo.xml"    // 服务器参数配置文件
	SYS_CONFIG  = "SysConfig.xml" // 配置工具配置文件名
	SYS_DEPLOY  = "SysDeploy.xml" // 系统部署配置文件
	ONEMAP_NAME = "OneMap"        // OneMap directory name
	TOMCAT      = "Tomcat"        // Tomcat container
	WEBLOGIC    = "Weblogic"      // Weblogic container
)

// 安装程序的三种状态
const (
	MAINTAIN    = iota // 0, 当前节点维持现状不做任何修改
	INSTALL            // 1, 当前节点需要进行安装操作
	UPDATE             // 2, 当前节点需要进行更新操作
	UNINSTALL          // 3, 当前节点需要进行卸载操作
	SUBPLATFORM        // 4, 
)

// 服务器基本信息
type SrvInfo struct {
	Os      string // 目标机器操作系统类型，[windos/linux]
	Arch    string // 目标机器系统架构，[386/amd64]
	MacName string // 机器名
	Domain  string // 域名
	Ip      string // 目标机器ip地址
	Root    string // 目标机器系统登录用户名
	Pwd     string // 目标机器系统登录密码
}

// 服务器集群信息
type Cluster struct {
	Type    string // 集群类型
	Enabled int    // 是否启用[0:启用/1:不启用]
	CIP     string // 集群IP
}

// OneMap安装信息
type OMPInfo struct {
	SrvInfo

	Cluster

	Deploy int // 部署情况, [0:维持现状/1:安装/2:更新/3:卸载] 

	Version   string   // onemap版本号
	OMHome    string   // OneMap安装目录
	Container string   // OneMap web容器类型
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

var (
	omsc    *SysConfig     // SysConfig struct
	omsm    *ServerMapping // ServerMapping struct
	omsd    *SysDeploy     // SysDeploy struct
	basedir string         // the working directory
	status  int            // the process's status [MAINTAIN | INSTALL | UPDATE | UNINSTALL]
)

// Package the OneMap
func (om *OMPInfo) OMPackage() error {
	dstdir := filepath.FromSlash(basedir + "/" + ONEMAP_NAME)
	l.Message("Make OneMap directory")
	/*	if flag := utl.Exists(dstdir); flag == true {
			// 首先删除原来的
			if err := os.RemoveAll(dstdir); err != nil {
				l.Errorf("Remove the old OneMap package failed")
				return err
			}
			// 再创建新的空文件夹
			if err := os.Mkdir(dstdir, 0777); err != nil {
				l.Errorf("Make OneMap directory failed")
				return err
			}
		}
	*/
	l.Message("Package OneMap")
	srcdir := basedir
	// search the OneMap package
	subpath, err := utl.GetSubDir(basedir)
	if err != nil || len(subpath) <= 0 {
		return errors.New("Get sub directory failed")
	}
	for _, thepath := range subpath {
		l.Debugf("The subpath is %s", thepath)
		temp := path.Base(thepath)
		temp = strings.ToUpper(temp)
		if strings.Index(temp, strings.ToUpper(ONEMAP_NAME)) < 0 { // onemap flag
			continue
		}

		if strings.Index(temp, strings.ToUpper(CurOS)) < 0 { // windows or linux flag
			continue
		}

		if strings.Index(temp, "_V") > 0 { // v*.*flag
			srcdir += "/" + thepath
			break
		}
	}
	if err := om.OMCopy(srcdir, dstdir); err != nil {
		l.Errorf("Package OneMap failed")
		return err
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

	if (utl.Exists(src)) != true {
		msg := "Source directory (" + src + ") isn't existed"
		return errors.New(msg)
	}
	si, err := os.Stat(src)
	if err != nil {
		l.Error(err)
		return err
	}
	if (utl.Exists(dst)) == true {
		if err := os.RemoveAll(dst); err != nil {
			l.Errorf("Remove the old OneMap package failed")
			return err
		}
	}
	if err := os.Mkdir(dst, si.Mode()); err != nil {
		l.Errorf("Make OneMap directory failed")
		return err
	}
	// copy public
	var inst_script string
	switch CurOS {
	case "linux":
		inst_script = filepath.FromSlash(src + "/install.sh")
	case "windows":
		inst_script = filepath.FromSlash(src + "/install.bat")
	}
	l.Debugf(inst_script)
	l.Debugf(dst)
	if err := utl.Copy(inst_script, dst); err != nil {
		msg := "Copy install bash script failed"
		return errors.New(msg)
	}
	if (utl.Exists(dst + "/arcgis")) != true {
		if err := utl.Copy(src+"/arcgis", dst); err != nil {
			msg := "Copy directory arcgis directory failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/bin")) != true {
		if err := utl.Copy(src+"/bin", dst); err != nil {
			msg := "Copy directory bin directory failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/config")) != true {
		if err := utl.Copy(src+"/config", dst); err != nil {
			msg := "Copy directory config directory failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/java")) != true {
		if err := utl.Copy(src+"/java", dst); err != nil {
			msg := "Copy directory java directory failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/temp")) != true {
		if err := utl.Copy(src+"/temp", dst); err != nil {
			msg := "Copy directory temp directory failed"
			return errors.New(msg)
		}
	}

	if (utl.Exists(dst + "/example_data")) != true {
		if err := utl.Copy(src+"/example_data", dst); err != nil {
			msg := "Copy directory(" + src + "/example_data) failed"
			return errors.New(msg)
		}
	}

	// copy modules
	if len(om.Apps) > 0 {
		// copy web container
		if err := utl.Copy(src+"/"+om.Container, dst); err != nil {
			msg := "Copy OneMap web container " + om.Container + "  container failed"
			return errors.New(msg)
		}

		for i := 0; i < len(om.Apps); i++ {
			if err := utl.Copy(src+"/webapps/"+om.Apps[i], dst+"/"+om.Container+"/webapps/"+om.Apps[i]); err != nil {
				msg := "Copy module (" + om.Apps[i] + ") failed"
				return errors.New(msg)
			}
		}
	}

	// copy services  
	if len(om.Services) > 0 {
		for i := 0; i < len(om.Services); i++ {
			if err := utl.Copy(src+"/services/"+om.Services[i], dst+"/services/"+om.Services[i]); err != nil {
				msg := "Copy OneMap service (" + om.Services[i] + ") failed"
				return errors.New(msg)
			}
		}
	}

	// deal with the db server
	for _, srv := range om.Servers {
		if srv == "db" {
			if err := utl.Copy(src+"/db/Driver", dst+"/db/Driver"); err != nil {
				msg := "Copy OneMap Driver directory failed"
				return errors.New(msg)
			}

			if err := utl.Copy(src+"/db/GeoCoding", dst+"/db/GeoCoding"); err != nil {
				msg := "Copy OneMap GeoCoding directory failed"
				return errors.New(msg)
			}

			if err := utl.Copy(src+"/db/GeoPortal", dst+"/db/GeoPortal"); err != nil {
				msg := "Copy OneMap GeoPortal directory failed"
				return errors.New(msg)
			}

			if err := utl.Copy(src+"/db/GeoShareManager", dst+"/db/GeoShareManager"); err != nil {
				msg := "Copy OneMap GeoShareManager directory failed"
				return errors.New(msg)
			}

			if err := utl.Copy(src+"/db/Portal", dst+"/db/Portal"); err != nil {
				msg := "Copy OneMap Portal directory failed"
				return errors.New(msg)
			}

			// install subplatform module
			if SubFlag {
				// delete the original file(Manager_Table_Data.sql)
				sqlfile := filepath.FromSlash(dst + "/db/GeoShareManager/Manager_Table_Data.sql")
				if err := os.Remove(sqlfile); err != nil {
					msg := "Delete the file(" + sqlfile + ") failed!"
					l.Errorf(msg)
					return errors.New(msg)
				}
				// rename the bak file(Manager_Table_Data.sql)
				bakfile := filepath.FromSlash(dst + "/db/GeoShareManager/Manager_Table_Data.sql.bak")
				if err := os.Rename(bakfile, sqlfile); err != nil {
					msg := "Rename the file(" + bakfile + ") failed!"
					l.Errorf(msg)
					return errors.New(msg)
				}
				// copy the subplatform
				if err := utl.Copy(src+"/db/SubPlatform", dst+"/db"); err != nil {
					msg := "Copy OneMap SubPlatform directory failed"
					return errors.New(msg)
				}
			}
		}
	}

	l.Message("Copy OneMap files successfully")
	return nil
}

// parse the current machine info
func (om *OMPInfo) OMGetInfo(mi *node, sm *ServerMapping, lo *Layout) error {
	l.Message("Begin to init with machine's info")
	if mi == nil || sm == nil {
		msg := "Input MachineInfo and SrvMapping object is nil"
		return errors.New(msg)
	}

	// get attributes
	for i := 0; i < len(mi.Attrs); i++ {
		attname := mi.Attrs[i].Attrname
		attvalue := mi.Attrs[i].Attrvalue
		if attname == "" {
			l.Warningf("The %dth attribute name is null", i+1)
			continue
		}
		if attvalue == "" {
			l.Warningf("The %dth attribute value is null", i+1)
		}
		l.Debugf("The %dth attribute name is:%s, and value is:%s", i+1, attname, attvalue)

		switch attname {
		case "os":
			om.Os = attvalue
		case "arch":
			om.Arch = attvalue
		case "macname":
			om.MacName = attvalue
		case "domname":
			om.Domain = attvalue
		case "ip":
			om.Ip = attvalue
		case "user":
			om.Root = attvalue
		case "pwd":
			om.Pwd = attvalue
		case "omhome":
			om.OMHome = attvalue
		case "deploy":
			om.Deploy, _ = strconv.Atoi(attvalue)
		case "cluster_type":
			om.Type = attvalue
		case "cluster_enabled":
			om.Enabled, _ = strconv.Atoi(attvalue)
		case "cluster_ip":
			om.CIP = attvalue
		}
	}

	// Get the web app modules name and services name
	for i := 0; i < len(mi.Srvs); i++ {
		srvtype := mi.Srvs[i].Srvname

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
			for j := 0; j < len(lo.Servers); j++ {
				var srv *ServerInfo = &(lo.Servers[j])
				if srv.Srvname == srvtype {
					for _, attr := range srv.Attrs {
						if attr.Attrname == "db_sid" {
							om.ORCL_SID = attr.Attrvalue
						} else if attr.Attrname == "db_user" {
							om.ORCL_User = attr.Attrvalue
						} else if attr.Attrname == "system_user" {
							om.DB_User[0] = attr.Attrvalue
						} else if attr.Attrname == "system_pwd" {
							om.DB_PWD[0] = attr.Attrvalue
						} else if attr.Attrname == "manager_user" {
							om.DB_User[1] = attr.Attrvalue
						} else if attr.Attrname == "manager_pwd" {
							om.DB_PWD[1] = attr.Attrvalue
						} else if attr.Attrname == "portal_user" {
							om.DB_User[2] = attr.Attrvalue
						} else if attr.Attrname == "portal_pwd" {
							om.DB_PWD[2] = attr.Attrvalue
						} else if attr.Attrname == "geocoding_user" {
							om.DB_User[3] = attr.Attrvalue
						} else if attr.Attrname == "geocoding_pwd" {
							om.DB_PWD[3] = attr.Attrvalue
						} else if attr.Attrname == "geoportal_user" {
							om.DB_User[4] = attr.Attrvalue
						} else if attr.Attrname == "geoportal_pwd" {
							om.DB_PWD[4] = attr.Attrvalue
						} else if attr.Attrname == "sub_user" {
							om.DB_User[5] = attr.Attrvalue
						} else if attr.Attrname == "sub_pwd" {
							om.DB_PWD[5] = attr.Attrvalue
						}
					}
				}
			}
		} else if srvtype == "gis" {
			for j := 0; j < len(lo.Servers); j++ {
				var srv *ServerInfo = &(lo.Servers[j])
				if srv.Srvname == srvtype {
					for _, attr := range srv.Attrs {
						if attr.Attrname == "ags_log_path" {
							om.AGS_Home = attr.Attrvalue
							/// AGS default log path
							attr.Attrvalue += "/server/user/log"
						}
					}
				}
			}
		}
	}

	/// default params
	if CurOS == "linux" {
		om.OM_Group = "esri"
		om.OM_User = "esri"
		om.OM_PWD = "esri1234"
	}

	return nil
}

// get onemap's version
func (om *OMPInfo) OMGetVersion(basedir string) (string, error) {
	var base string // the onemap package directory
	if flag := utl.Exists(basedir); flag != true {
		msg := "Input directory(" + basedir + ") isn't existed"
		return base, errors.New(msg)
	}

	// get all sub directory name and search the onemap package directory
	subpath, err := utl.GetSubDir(basedir)
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
	if flag := utl.Exists(basedir + "/" + subdirname); flag != true {
		msg := "Input directory(" + basedir + ") isn't existed"
		return errors.New(msg)
	}

	// get all sub directory name and search the onemap package directory
	subpath, err := utl.GetSubDir(basedir + "/" + subdirname)
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
