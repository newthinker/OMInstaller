package main

import (
	"encoding/json"
	"fmt"
	"github.com/newthinker/onemap-installer/log"
	"github.com/newthinker/onemap-installer/sys"
	"github.com/newthinker/onemap-installer/utl"
	"github.com/newthinker/onemap-installer/web"
	"net/http"
	"os/exec"
	"path/filepath"
)

func main() {
	////////////////////////////////////////////////////////////////
	// init log
	l, err := log.NewLog("inst.log", log.LogAll, log.DefaultBufSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	////////////////////////////////////////////////////////////////
	// get local ip
	ip, err := utl.GetNetIP()
	if err != nil {
		l.Errorf("Get local ip failed")
		return
	}
	l.Messagef("Get local ip: %s", ip)
	////////////////////////////////////////////////////////////////
	// install sshpass
	l.Message("Install the sshpass")
	base, err := filepath.Abs("./") // 获取系统当前路径
	fmt.Println("base:" + base)
	if err != nil || base == "" {
		l.Errorf("Get current directory failed")
		return
	}

	// whether installed
	if flag := utl.Exists(base + "/sshpass/bin/sshpass"); flag != true {
		if flag = utl.Exists(base + "/sshpass/Install.sh"); flag != true {
			l.Errorf("No sshpass software package")
			return
		}

		// exec the install script
		cmd := exec.Command("/bin/sh", base+"/sshpass/Install.sh", base+"/sshpass")
		err = cmd.Run()
		if err != nil {
			l.Errorf("Complier sshpass failed")
			return
		}

		// whether install successfully
		if flag := utl.Exists(base + "/sshpass/bin/sshpass"); flag != true {
			l.Errorf("Install sshpass failed")
			return
		}
	} else {
		l.Messagef("Sshpass is installed and go on")
	}

	cmd := exec.Command(base+"/sshpass/bin/sshpass", "-V")
	err = cmd.Run()
	if err != nil {
		l.Errorf("Sshpass isn't installed")
		return
	}

	////////////////////////////////////////////////////////////////
	/// post.json test
    str := `{"data":[{
		"base": [{
			"Attrname": "os",
			"Attrvalue": "linux"
		},
		{
			"Attrname": "arch",
			"Attrvalue": "amd64"
		},
		{
			"Attrname": "ip",
			"Attrvalue": "127.0.0.1"
		},
		{
			"Attrname": "user",
			"Attrvalue": "admin"
		},
		{
			"Attrname": "pwd",
			"Attrvalue": "admin"
		},
		{
			"Attrname": "omhome",
			"Attrvalue": "/opt/OneMap"
		}],
		"params": [{
			"Srvname": "db",
			"Attrs": [{
				"Attrname": "db_ip",
				"Attrvalue": "数据库服务器IP"
			},
			{
				"Attrname": "db_port",
				"Attrvalue": "数据库服务器PORT"
			},
			{
				"Attrname": "db_sid",
				"Attrvalue": "数据库服务器实例名"
			},
			{
				"Attrname": "manager_user",
				"Attrvalue": "GeoShareManager数据库用户名"
			},
			{
				"Attrname": "manager_pwd",
				"Attrvalue": "GeoShareManager数据库密码",
				"Encrypt": "true"
			},
			{
				"Attrname": "portal_user",
				"Attrvalue": "GeoSharePortal数据库用户名"
			},
			{
				"Attrname": "portal_pwd",
				"Attrvalue": "GeoSharePortal数据库密码",
				"Encrypt": "true"
			},
			{
				"Attrname": "geocoding_user",
				"Attrvalue": "GeoCoding数据库用户名"
			},
			{
				"Attrname": "geocoding_pwd",
				"Attrvalue": "GeoCoding数据库密码",
				"Encrypt": "true"
			},
			{
				"Attrname": "geoportal_user",
				"Attrvalue": "GeoPortal数据库用户名"
			},
			{
				"Attrname": "geoportal_pwd",
				"Attrvalue": "GeoPortal数据库密码",
				"Encrypt": "true"
			},
			{
				"Attrname": "db_h2",
				"Attrvalue": "内存数据库IP"
			},
			{
				"Attrname": "sub_user",
				"Attrvalue": "分平台数据库用户名"
			}]
		},
		{
			"Srvname": "gis",
			"Attrs": [{
				"Attrname": "gis_ip",
				"Attrvalue": "GIS服务器IP"
			},
			{
				"Attrname": "gis_port",
				"Attrvalue": "GIS服务器端口"
			},
			{
				"Attrname": "gis_ags_ip",
				"Attrvalue": "ArcGISServer服务器IP"
			},
			{
				"Attrname": "gis_ags_user",
				"Attrvalue": "ArcGISServer服务器用户名"
			},
			{
				"Attrname": "gis_ags_pwd",
				"Attrvalue": "ArcGISServer服务器密码",
				"Encrypt": "true"
			},
			{
				"Attrname": "gis_sharekey",
				"Attrvalue": "sharekey",
				"Encrypt": "true"
			},
			{
				"Attrname": "ags_log_path",
				"Attrvalue": "ArcGISServer日志路径"
			},
			{
				"Attrname": "gis_basemap_type",
				"Attrvalue": "底图服务类型"
			},
			{
				"Attrname": "gis_services_port",
				"Attrvalue": "RemoteServices系统的HTTPSPORT"
			},
			{
				"Attrname": "wmts_ip",
				"Attrvalue": "WMTS系统的IP"
			},
			{
				"Attrname": "wmts_port",
				"Attrvalue": "WMTS系统的PORT"
			},
			{
				"Attrname": "sysrest_ip",
				"Attrvalue": "SysRest服务IP"
			},
			{
				"Attrname": "sysrest_port",
				"Attrvalue": "SysRest服务PORT"
			}]
		},
		{
			"Srvname": "main",
			"Attrs": [{
				"Attrname": "main_ip",
				"Attrvalue": "运维服务器IP"
			},
			{
				"Attrname": "main_port",
				"Attrvalue": "运维服务器端口"
			},
			{
				"Attrname": "main_upload_level",
				"Attrvalue": "运维文件上传目录深度"
			},
			{
				"Attrname": "main_upload_num",
				"Attrvalue": "运维文件上传最大值"
			},
			{
				"Attrname": "geocoding_ip",
				"Attrvalue": "GeoCoding系统IP"
			},
			{
				"Attrname": "geocoding_port",
				"Attrvalue": "GeoCoding系统PORT"
			},
			{
				"Attrname": "geoportal_ip",
				"Attrvalue": "GeoPortal系统IP"
			},
			{
				"Attrname": "geoportal_portal",
				"Attrvalue": "GeoPortal系统PORT"
			},
			{
				"Attrname": "sub_upload_num",
				"Attrvalue": "分平台文件上传最大值"
			},
			{
				"Attrname": "sub_syn_code",
				"Attrvalue": "0"
			},
			{
				"Attrname": "file_ip",
				"Attrvalue": "FileServices系统IP"
			},
			{
				"Attrname": "file_port",
				"Attrvalue": "FileServices系统的PORT"
			},
			{
				"Attrname": "tile_ip",
				"Attrvalue": "TileServices系统IP"
			},
			{
				"Attrname": "tile_port",
				"Attrvalue": "TileServices系统PORT"
			},
			{
				"Attrname": "aggregator_ip",
				"Attrvalue": "Aggregator系统IP"
			},
			{
				"Attrname": "aggregator_port",
				"Attrvalue": "Aggregator系统PORT"
			}]
		},
		{
			"Srvname": "agent",
			"Attrs": [{
				"Attrname": "agent_ip",
				"Attrvalue": "监控代理服务器IP"
			},
			{
				"Attrname": "agent_port",
				"Attrvalue": "监控代理服务器PORT"
			}]
		},
		{
			"Srvname": "msg",
			"Attrs": [{
				"Attrname": "jms_ip",
				"Attrvalue": "JMS服务器IP"
			},
			{
				"Attrname": "jms_port",
				"Attrvalue": "JSM服务器PORT"
			}]
		},
		{
			"Srvname": "token",
			"Attrs": [{
				"Attrname": "token_ip",
				"Attrvalue": "Token服务器IP"
			},
			{
				"Attrname": "token_port",
				"Attrvalue": "Token服务器PORT"
			}]
		},
		{
			"Srvname": "web",
			"Attrs": [{
				"Attrname": "web_ip",
				"Attrvalue": "门户IP"
			},
			{
				"Attrname": "web_port",
				"Attrvalue": "门户PORT"
			},
			{
				"Attrname": "web_user",
				"Attrvalue": "门户默认用户名"
			},
			{
				"Attrname": "web_pwd",
				"Attrvalue": "门户默认用户名密码"
			},
			{
				"Attrname": "flex_ip",
				"Attrvalue": "Flex发布服务器IP"
			},
			{
				"Attrname": "flex_port",
				"Attrvalue": "Flex发布服务器PORT"
			},
			{
				"Attrname": "sl_ip",
				"Attrvalue": "Silverlight发布服务器IP"
			},
			{
				"Attrname": "sl_port",
				"Attrvalue": "Silverlight发布服务器PORT"
			},
			{
				"Attrname": "ria_ip",
				"Attrvalue": "快速搭建服务器IP"
			},
			{
				"Attrname": "ria_port",
				"Attrvalue": "快速搭建服务器PORT"
			},
			{
				"Attrname": "iis_root",
				"Attrvalue": "Silverlight文件生成路径"
			},
			{
				"Attrname": "web_csw",
				"Attrvalue": "0"
			}]
		}]
	}]}`
	jsonstr := make(map[string]interface{})
	err = json.Unmarshal([]byte(str), &jsonstr)
	if err != nil {
		l.Error(err)
		fmt.Println(err)
		return
	}
	sd, arrlo, err := sys.ParseSysSubmit(jsonstr)
	if err != nil {
		l.Error(err)
		fmt.Println(err)
		return
	}
	l.Debug(sd)
	l.Debug(arrlo)

	////////////////////////////////////////////////////////////////
	l.Message("Listen and serve")
	web.Init(l)

	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))
	http.Handle("/images/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/subconfig", web.SubHandler)
	http.HandleFunc("/sysconfig", web.SysConfig)
	http.HandleFunc("/syshandler", web.SysHandler)
	http.HandleFunc("/error", web.ErrHandler)

	err = http.ListenAndServe(ip+":8888", nil)
	if err != nil {
		fmt.Println(err)
		l.Errorf("Listen and serve failed: %s", err)
		return
	}
	///////////////////////////////////////////////////////////////	
}
