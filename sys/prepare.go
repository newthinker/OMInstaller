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
	"strconv"
	"strings"
)

const (
	SERVER_MAPPING = "SrvMapping.xml" // 服务器类型映射文件名 
	SYS_INFO       = "SysInfo.xml"    // 服务器参数配置文件
	SYS_CONFIG     = "SysConfig.xml"  // 配置工具配置文件名
	SYS_DEPLOY     = "SysDeploy.xml"  // 系统部署配置文件
	ONEMAP_NAME    = "OneMap"         // OneMap directory name
	TOMCAT         = "Tomcat"         // Tomcat container
	WEBLOGIC       = "Weblogic"       // Weblogic container
)

// 安装程序的三种状态
const (
	INSTALL = iota // INSTALL =0
	UPDATE
	UNINSTALL
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

	Deploying int // 部署情况, [0:部署成功/1:部署失败/2:进行部署/3:暂不部署] 

	//	Basedir string // 当前安装包根目录

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
	omsc    *SysConfig
	omsm    *ServerMapping
	basedir string // 系统当前运行目录
)

/*
func Init() error {
	// 获取当前目录
	dir, err := utl.GetLocalDir()
	if err != nil {
		l.Error(err)
		return err
	}
	basedir = dir
	l.Debugf("Current directory is:%s", basedir)

	// open the config files
	sm, err1 := OpenSMConfig(basedir)
	sc, err2 := OpenSCConfig(basedir)
	if err1 != nil || err2 != nil {
		l.Error(errors.New("Parse system config files failed"))
		return errors.New("Parse system config files failed")
	}
	omsc = &sc
	omsm = &sm

	return nil
}
*/

func (om *OMPInfo) OMPackage() error {
	dstdir := basedir + "/" + ONEMAP_NAME
	l.Message("Make OneMap directory")
	if flag := utl.Exists(dstdir); flag == true {
		// 首先删除原来的
		if err := os.RemoveAll(dstdir); err != nil {
			l.Errorf("Remove the old OneMap package failed")
			return err
		}
		// 再创建新的空文件夹
		if err := os.Mkdir(dstdir, 0755); err != nil {
			l.Errorf("Make OneMap directory failed")
			return err
		}
	}

	l.Message("Package OneMap")
	srcdir := basedir + "/" + ONEMAP_NAME + "_Linux_" + om.Version
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
	if (utl.Exists(dst)) != true {
		if err := os.Mkdir(dst, 0755); err != nil {
			msg := "Destination directory (" + dst + ") isn't existed and create failed"
			return errors.New(msg)
		}
	} else {
		if (utl.Exists(dst + "/services")) == true {
			if err := os.RemoveAll(dst + "/services"); err != nil {
				msg := "Remove directory (" + dst + "/services) failed"
				return errors.New(msg)
			}
		}
		if (utl.Exists(dst + "/" + om.Container + "/webapps")) == true {
			if err := os.RemoveAll(dst + "/" + om.Container + "/webapps"); err != nil {
				msg := "Remove directory (" + dst + "/" + om.Container + "/webapps) failed"
				return errors.New(msg)
			}
		}
	}

	// copy public
	if (utl.Exists(dst + "/install.sh")) != true {
		if err := utl.Copy(src+"/install.sh", dst); err != nil {
			msg := "Copy install bash script failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/arcgis")) != true {
		if err := utl.Copy(src+"/arcgis", dst); err != nil {
			msg := "Copy directory (" + src + "/arcgis) failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/bin")) != true {
		if err := utl.Copy(src+"/bin", dst); err != nil {
			msg := "Copy directory(" + src + "/bin) failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/config")) != true {
		if err := utl.Copy(src+"/config", dst); err != nil {
			msg := "Copy directory(" + src + "/config) failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/java")) != true {
		if err := utl.Copy(src+"/java", dst); err != nil {
			msg := "Copy directory(" + src + "/java) failed"
			return errors.New(msg)
		}
	}
	if (utl.Exists(dst + "/temp")) != true {
		if err := utl.Copy(src+"/temp", dst); err != nil {
			msg := "Copy directory(" + src + "/temp) failed"
			return errors.New(msg)
		}
	}

	// copy modules
	if len(om.Apps) > 0 {
		// copy web container
		if err := utl.Copy(src+"/"+om.Container, dst); err != nil {
			msg := "Copy OneMap web container (" + om.Container + ") failed"
			return errors.New(msg)
		}

		for i := 0; i < len(om.Apps); i++ {
			if err := utl.Copy(src+"/webapps/"+om.Apps[i], dst+"/"+om.Container+"/webapps/"+om.Apps[i]); err != nil {
				msg := "Copy module (" + om.Apps[i] + ") failed"
				return errors.New(msg)
			}

			switch om.Apps[i] {
			case "H2memDB":
				if err := utl.Copy(src+"/db", dst); err != nil {
					return errors.New("Copy db directory failed")
				}
			case "GeoShareManager":
				if err := utl.Copy(src+"/example_data", dst); err != nil {
					return errors.New("Copy example data directory failed")
				}
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

			// 如果包含H2MemDB模块，且在临时目录中包含分平台数据更新文件，还需要更新分平台数据文件
			if om.Services[i] == "H2CommonMemDB" {
				// 首先拷贝db目录
				if err := utl.Copy(src+"/db", dst); err != nil {
					msg := "Copy OneMap db directory failed"
					return errors.New(msg)
				}
				// 更新数据库文件
				srcfile := basedir + "/temp/Manager_Table_Data.sql"
				dstfile := dst + "/db/GeoShareManager/Manager_Table_Data.sql"
				if utl.Exists(basedir+"/temp/Manager_Table_Data.sql") == true {
					if err := utl.Copy(srcfile, dstfile); err != nil {
						msg := "Update subplatform SQL file failed"
						return errors.New(msg)
					}

					// 删除临时目录
					os.RemoveAll(basedir + "/temp")
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
			om.Deploying, _ = strconv.Atoi(attvalue)
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
							om.ORCL_SID = attr.Value
						} else if attr.Attrname == "db_user" {
							om.ORCL_User = attr.Value
						} else if attr.Attrname == "system_user" {
							om.DB_User[0] = attr.Value
						} else if attr.Attrname == "system_pwd" {
							om.DB_PWD[0] = attr.Value
						} else if attr.Attrname == "manager_user" {
							om.DB_User[1] = attr.Value
						} else if attr.Attrname == "manager_pwd" {
							om.DB_PWD[1] = attr.Value
						} else if attr.Attrname == "portal_user" {
							om.DB_User[2] = attr.Value
						} else if attr.Attrname == "portal_pwd" {
							om.DB_PWD[2] = attr.Value
						} else if attr.Attrname == "geocoding_user" {
							om.DB_User[3] = attr.Value
						} else if attr.Attrname == "geocoding_pwd" {
							om.DB_PWD[3] = attr.Value
						} else if attr.Attrname == "geoportal_user" {
							om.DB_User[4] = attr.Value
						} else if attr.Attrname == "geoportal_pwd" {
							om.DB_PWD[4] = attr.Value
						} else if attr.Attrname == "sub_user" {
							om.DB_User[5] = attr.Value
						} else if attr.Attrname == "sub_pwd" {
							om.DB_PWD[5] = attr.Value
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
							om.AGS_Home = attr.Value
							/// AGS default log path
							attr.Value += "/server/user/log"
						}
					}
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
