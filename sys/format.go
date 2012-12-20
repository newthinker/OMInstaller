package sys

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
)

///////////////////////////////
/// get.json中Server_modules
type MdlGet struct {
	Modname string
	Moddesc string
}

type SrvMdl struct {
	Srvname string
	Srvdesc string
	Modules []MdlGet
}

type SrvsMdl struct {
	Server_modules []SrvMdl
}

/// get.json中Server_params
type ParamGet struct {
	Paramname string
	Paramdesc string
	Encrypt   string
	Selects   string
}

type SrvParam struct {
	Srvname string
	Params  []ParamGet
}

type SrvsParam struct {
	Server_params []SrvParam
}

///////////////////////////////
/// post.json 
type SrvBase struct {
	Os        string
	Arch      string
	Ip        string
	User      string
	Pwd       string
	Omhome    string
	Container string
}

type ParamPost struct {
	Paramname  string
	Paramvalue string
	Encrypt    string
	Selects    string
}

type SrvPost struct {
	Srvname string
	Params  []ParamPost
}

type ServerParams struct {
	Server_base   SrvBase
	Server_params []SrvPost
}

///////////////////////////////

// 将SrvMapping整理输出到map中
func FormatSrvMapping(sm ServerMapping) SrvsMdl {
	var servers SrvsMdl // 服务器数组

	for i := 0; i < len(sm.Servers); i++ {
		var srv *Server = &(sm.Servers[i])
		if srv.XMLName.Local == "" {
			continue
		}

		var srvmap SrvMdl // 保存服务器信息
		srvmap.Srvname = srv.XMLName.Local
		srvmap.Srvdesc = srv.SrvDesc

		for j := 0; j < len(srv.ModuleList); j++ {
			var mdl Module = srv.ModuleList[j]
			if mdl.MdlName != "" {
				var mdlmap MdlGet
				mdlmap.Modname = mdl.MdlName
				mdlmap.Moddesc = mdl.MdlDesc

				srvmap.Modules = append(srvmap.Modules, mdlmap)
			}
		}

		servers.Server_modules = append(servers.Server_modules, srvmap)
	}

	//	result := make(map[string]interface{})
	//	result["Server_modules"] = servers

	return servers
}

// 将SysConfig整理输出到map中
func FormatSysConfig(sc SysConfig) SrvsParam {
	var servers SrvsParam

	for i := 0; i < len(sc.LayOut.Servers); i++ {
		srvinfo := &(sc.LayOut.Servers[i])

		if srvinfo.XMLName.Local == "" {
			continue
		}

		var srvmap SrvParam
		srvmap.Srvname = srvinfo.XMLName.Local

		for j := 0; j < len(srvinfo.Attrs); j++ {
			attr := &(srvinfo.Attrs[j])

			if attr != nil && attr.XMLName.Local != "" && attr.AttrName != "" {
				var param ParamGet

				param.Paramname = attr.XMLName.Local
				param.Paramdesc = attr.AttrName

				// 判断是否需要加密
				if attr.AttrEncrypt != "" {
					param.Encrypt = "true"
				}

				// 判断是否是select框 
				if attr.AttrSelect != "" {
					param.Selects = attr.AttrSelect
				}

				srvmap.Params = append(srvmap.Params, param)
			}
		}

		servers.Server_params = append(servers.Server_params, srvmap)
	}

	return servers
}

// 解析POST的JSON结构
func ParseSysSubmit(jsonstr interface{}, basepath string) error {
	postmap := jsonstr.(map[string]interface{})

	for _, v := range postmap {
		switch vv := v.(type) {
		case string:
		case int:
			if vv != 0 {
				return errors.New("前端返回码错误，请检查!")
			}
		case []interface{}:
			/// 暂时输出到SysInfo.xml文件中
			sysinfo := &SysInfo{}

			// 开始解析数据体部分
			for i, s := range vv {
				fmt.Printf("解析第%d个服务器参数：\n", i)
				srvparams := s.(ServerParams)

				machine := &MachineInfo{}
				machine.Os = srvparams.Server_base.Os
				machine.Arch = srvparams.Server_base.Arch
				machine.Ip = srvparams.Server_base.Ip
				machine.User = srvparams.Server_base.User
				machine.Pwd = srvparams.Server_base.Pwd
				machine.Omhome = srvparams.Server_base.Omhome
				machine.Web = srvparams.Server_base.Container

				for _, sp := range srvparams.Server_params {
					srvpost := sp //.(SrvPost)
					server := &ServerInfo{}
					server.XMLName.Local = srvpost.Srvname

					for _, p := range sp.Params {
						attr := &AttrInfo{}
						attr.XMLName.Local = p.Paramname
						attr.Value = p.Paramvalue

						if p.Encrypt != "" {
							attr.AttrEncrypt = p.Encrypt
						}
						if p.Selects != "" {
							attr.AttrSelect = p.Selects
						}

						server.Attrs = append(server.Attrs, *attr)
					}

					machine.Servers = append(machine.Servers, *server)
				}

				sysinfo.Machines = append(sysinfo.Machines, *machine)
			}

			// 输出到SysInfo.xml文件中
			output, err1 := xml.MarshalIndent(sysinfo, " ", "    ")

			if err1 != nil {
				return err1
			}

			sysconfig := basepath + "/conf/" + SYS_INFO

			// 如果配置文件已存在，先将其删除
			if Exists(sysconfig) == true {
				if err4 := os.Remove(sysconfig); err4 != nil {
					return err4
				}
			}
			// 将输出字符流写入文件中
			file, err2 := os.Open(sysconfig)
			defer file.Close()
			if err2 != nil {
				return err2
			}

			_, err3 := file.Write(output)
			if err3 != nil {
				return err3
			}
		}
	}

	return nil
}
