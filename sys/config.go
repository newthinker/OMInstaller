package sys

import (
	"container/list"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

///////////////////////////////////////////////////////
// ServerMapping struct
type ServerMapping struct {
	XMLName     xml.Name `xml:"root"`
	Servers     []Server `xml:",any"`
	Description string   `xml:",innerxml"`
}

type Server struct {
	XMLName    xml.Name `xml:""`
	SrvDesc    string   `xml:"desc,attr"`
	ModuleList []Module `xml:",any"`
}

type Module struct {
	XMLName xml.Name `xml:""`
	MdlDesc string   `xml:"desc,attr"`
	MdlName string   `xml:",chardata"`
}

func (s *Server) AddModule(mdltype string, mdlname string, mdldesc string) {
	newm := Module{MdlName: mdlname, MdlDesc: mdldesc}
	newm.XMLName.Local = mdltype
	s.ModuleList = append(s.ModuleList, newm)
}

func (ss *ServerMapping) AddServer(srvname string, srvdesc string, mdltype []string, mdlnames []string, mdldesc []string) {
	news := Server{SrvDesc: srvdesc}
	news.XMLName.Local = srvname
	for v := range mdlnames {
		news.AddModule(mdltype[v], mdlnames[v], mdldesc[v])
	}

	ss.Servers = append(ss.Servers, news)
}

/*
func (s *Server) GetCurrentValue(mdlname string) (name string, e error) {
  for _, v := range s.ModuleList {
    if v.XMLName.Local == mdlname {
      name = v.ModuleName
      return
    }
  }
  e = errors.New(fmt.Sprintf("%s not found", mdlname))
  return
}
*/

// 将SrvMapping整理输出到map中
func (sm *ServerMapping) FormatSrvMapping() map[string]interface{} {
	srvlist := list.New() // 服务器数组

	for i := 0; i < len(sm.Servers); i++ {
		var srv *Server = &(sm.Servers[i])
		if srv.XMLName.Local == "" {
			continue
		}

		srvmap := make(map[string]interface{}) // 保存服务器信息
		srvmap["Srvname"] = srv.XMLName.Local
		srvmap["Srvdesc"] = srv.SrvDesc

		for j := 0; j < len(srv.ModuleList); j++ {
			mdlmap := make(map[string]string) // 保存模块信息

			var mdl Module = srv.ModuleList[j]
			if mdl.MdlName != "" {
				mdlmap["Modname"] = mdl.MdlName
				mdlmap["Moddesc"] = mdl.MdlDesc

				srvmap["Modules"] = mdlmap
			}
		}

		srvlist.PushBack(srvmap)
	}

	servers := make(map[string]interface{})
	servers["Server_modules"] = srvlist

	return servers
}

///////////////////////////////////////////////////////
// SysInfo struct
type SysInfo struct {
	XMLName  xml.Name      `xml:"root"`
	Machines []MachineInfo `xml:",any"`
	//	Description string        `xml:",innerxml"`
}

type MachineInfo struct {
	XMLName xml.Name     `xml:"machine"`
	Os      string       `xml:"os,attr"`
	Arch    string       `xml:"arch,attr"`
	Ip      string       `xml:"ip,attr"`
	User    string       `xml:"user,attr"`
	Pwd     string       `xml:"pwd,attr"`
	Omhome  string       `xml:"omhome,attr"`
	Web     string       `xml:"container,attr"`
	Servers []ServerInfo `xml:",any"`
}

type ServerInfo struct {
	XMLName xml.Name   `xml:""`
	Attrs   []AttrInfo `xml:",any"`
}

type AttrInfo struct {
	XMLName     xml.Name `xml:""`
	AttrName    string   `xml:"name,attr"`
	AttrEncrypt string   `xml:"encrypt,attr"`
	AttrSelect  string   `xml:"selects,attr"`
	Value       string   `xml:",chardata"`
}

func (s *ServerInfo) AddAttrInfo(attrname string,
	attrvalue string, desc string, encrypt string, selects string) {
	newa := AttrInfo{Value: attrvalue, AttrName: desc, AttrEncrypt: encrypt, AttrSelect: selects}
	newa.XMLName.Local = attrname
	s.Attrs = append(s.Attrs, newa)
}

func (ma *MachineInfo) AddServerInfo(srvtype string,
	attrarray []AttrInfo) {
	news := ServerInfo{}
	news.XMLName.Local = srvtype
	news.Attrs = attrarray
	ma.Servers = append(ma.Servers, news)
}

func (sm *SysInfo) AddMachineInfo(os string,
	arch string, ip string, user string, pwd string,
	servers []ServerInfo, omhome string, container string) {
	newm := MachineInfo{Os: os, Arch: arch, Ip: ip, User: user, Pwd: pwd, Omhome: omhome, Web: container}
	newm.Servers = servers
	sm.Machines = append(sm.Machines, newm)
}

////////////////////////////////////////////////////////
// SysConfig struct
type SysConfig struct {
	XMLName    xml.Name `xml:"root"`
	OneMapHome string   `xml:"OneMapHome,attr"`
	LayOut     Layout   `xml:""`
	FileMap    Filemap  `xml:""`
}

type Layout struct {
	XMLName xml.Name     `xml:"Layout"`
	Servers []ServerInfo `xml:",any"`
}

type Filemap struct {
	XMLName    xml.Name    `xml:"FileMapping"`
	Containers []Container `xml:",any"`
}

type Container struct {
	XMLName xml.Name    `xml:""`
	Path    string      `xml:"path,attr"`
	Modules []ModuleMap `xml:",any"`
}

type ModuleMap struct {
	XMLName    xml.Name    `xml:""`
	ServersMap []ServerMap `xml:",any"`
}

/// 目前只解析到服务器类型这一层 
type ServerMap struct {
	XMLName xml.Name `xml:""`
}

func (lo *Layout) AddServerInfo(srvtype string, attrarray []AttrInfo) {
	news := ServerInfo{}
	news.XMLName.Local = srvtype
	news.Attrs = attrarray
	lo.Servers = append(lo.Servers, news)
}

func (mm *ModuleMap) AddServerMap(srvname string) {
	newm := ServerMap{}
	newm.XMLName.Local = srvname
	mm.ServersMap = append(mm.ServersMap, newm)
}

func (c *Container) AddModuleMap(mdlname string, arrServers []ServerMap) {
	newm := ModuleMap{}
	newm.XMLName.Local = mdlname
	c.Modules = append(c.Modules, newm)
}

func (fm *Filemap) AddContainer(conname string, conpath string, arrmodule []ModuleMap) {
	newc := Container{Path: conpath}
	newc.XMLName.Local = conname
	newc.Modules = arrmodule
	fm.Containers = append(fm.Containers, newc)
}

// 将SysConfig中的输入参数整理输出
func (sc *SysConfig) FormatSysConfig() map[string]interface{} {
	srvlist := list.New()

	for i := 0; i < len(sc.LayOut.Servers); i++ {
		srvinfo := &(sc.LayOut.Servers[i])

		if srvinfo.XMLName.Local == "" {
			continue
		}

		srvmap := make(map[string]interface{})

		srvmap["Srvname"] = srvinfo.XMLName.Local

		lstparams := list.New() // 属性列表

		for j := 0; j < len(srvinfo.Attrs); j++ {
			attr := &(srvinfo.Attrs[j])

			if attr != nil && attr.XMLName.Local != "" && attr.AttrName != "" {
				attrmap := make(map[string]string)

				attrmap["Paramname"] = attr.XMLName.Local
				attrmap["Paramdesc"] = attr.AttrName

				// 判断是需要需要加密
				if attr.AttrEncrypt != "" {
					attrmap["Encrypt"] = "true"
				}

				lstparams.PushBack(attrmap)
			}
		}

		srvmap["Params"] = lstparams

		srvlist.PushBack(srvmap)
	}

	servers := make(map[string]interface{})
	servers["Server_params"] = srvlist

	return servers
}

////////////////////////////////////////////////////////
// update SysConfig.xml file except MonitorAgent module
func UpdateConfig(si *SysInfo, sc *SysConfig) error {
	flag := make(map[string]bool) // flag of whether update
	// initialize
	for i := 0; i < len(sc.LayOut.Servers); i++ {
		var srvinfo *ServerInfo = &(sc.LayOut.Servers[i])
		var srvtype string = srvinfo.XMLName.Local
		flag[srvtype] = false
	}

	for i := 0; i < len(si.Machines); i++ {
		for j := 0; j < len(si.Machines[i].Servers); j++ {
			var si_srvinfo *ServerInfo = &(si.Machines[i].Servers[j])
			var si_srvtype string = si_srvinfo.XMLName.Local // server type

			// don't set MonitorAgent server temperarily
			if si_srvtype == "agent" {
				flag[si_srvtype] = true
				continue
			}

			for k := 0; k < len(sc.LayOut.Servers); k++ {
				var sc_srvinfo *ServerInfo = &(sc.LayOut.Servers[k])
				var sc_srvtype string = sc_srvinfo.XMLName.Local

				// update the matching server info
				if si_srvtype == sc_srvtype {
					sc.LayOut.Servers[k] = *si_srvinfo

					flag[sc_srvtype] = true // update the flag
				}
			}
		}
	}

	// chekc whether all modules are updated
	var num int = 0
	var nomodules string
	for key, value := range flag {
		if value == false {
			if num > 0 {
				nomodules = nomodules + ","
			}
			nomodules = nomodules + key
			num++
		}
	}
	if num > 0 && nomodules != "" {
		msg := "WARN: " + strconv.Itoa(num) + "个模块(" + nomodules + ") 没有更新"
		fmt.Println(msg)
		//		return errors.New(msg)
	}

	return nil
}

// update SysConfig.xml file's MonitorAgent module
func UpdateMdlAgent(mi *MachineInfo, sc *SysConfig) error {
	var flag bool = false
	for i := 0; i < len(mi.Servers); i++ {
		var mi_srvinfo *ServerInfo = &(mi.Servers[i])
		var mi_srvtype string = mi_srvinfo.XMLName.Local

		if mi_srvtype != "agent" {
			continue
		} else {
			for j := 0; j < len(sc.LayOut.Servers); j++ {
				var sc_srvinfo *ServerInfo = &(sc.LayOut.Servers[j])
				var sc_srvtype string = sc_srvinfo.XMLName.Local

				if sc_srvtype == mi_srvtype {
					sc.LayOut.Servers[j] = *mi_srvinfo
				}
			}

			flag = true
		}
	}

	if !flag {
		return errors.New("WARN: 配置文件中没有监控代理模块参数信息")
	}

	return nil
}

// open the config files
func OpenSMConfig(basedir string) (ServerMapping, error) {
	var sm ServerMapping

	// check the base dir whether existed
	if flag := Exists(basedir); flag != true {
		msg := "ERROR: 安装目录(" + basedir + ")不存在"
		return sm, errors.New(msg)
	}

	file, err := os.Open(basedir + "/conf/" + SERVER_MAPPING)
	if err != nil {
		return sm, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return sm, err
	}
	if err := xml.Unmarshal([]byte(data), &sm); err != nil {
		return sm, err
	}

	return sm, nil
}

func OpenSIConfig(basedir string) (SysInfo, error) {
	var si SysInfo

	// check the base dir whether existed
	if flag := Exists(basedir); flag != true {
		msg := "ERROR: 输入目录(" + basedir + ")不存在"
		return si, errors.New(msg)
	}

	file, err1 := os.Open(basedir + "/conf/" + SYS_INFO)
	if err1 != nil {
		return si, err1
	}
	data, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return si, err2
	}
	if err3 := xml.Unmarshal([]byte(data), &si); err3 != nil {
		return si, err3
	}

	return si, nil
}

func OpenSCConfig(basedir string) (SysConfig, error) {
	var sc SysConfig

	// check the base dir whether existed
	if flag := Exists(basedir); flag != true {
		msg := "ERROR: 输入目录(" + basedir + ")不存在"
		return sc, errors.New(msg)
	}

	file, err1 := os.Open(basedir + "/conf/" + SYS_CONFIG)
	if err1 != nil {
		return sc, err1
	}
	data, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return sc, err2
	}
	if err3 := xml.Unmarshal([]byte(data), &sc); err3 != nil {
		return sc, err3
	}

	return sc, nil
}

// 更新保存系统配置文件
func RefreshSysConfig(sc *SysConfig, conffile string) error {
	if sc == nil || conffile == "" {
		return errors.New("输入的配置文件为空、文件路径为空")
	}

	if flag := Exists(conffile); flag != true {
		return errors.New("配置文件路径不存在")
	}

	output, err := xml.MarshalIndent(sc, "  ", "   ")
	fmt.Println([]byte(output))
	os.Stdout.Write([]byte(output))
	if err != nil {
		return err
	}

	if Exists(conffile) == true {
		if err = os.Remove(conffile); err != nil {
			return err
		}
	}

	file, err1 := os.Create(conffile)
	defer file.Close()
	if err1 != nil {
		return err1
	}

	_, err2 := file.Write([]byte(xml.Header))
	if err2 != nil {
		return err2
	}

	_, err3 := file.Write(output)
	if err3 != nil {
		return err3
	}

	return nil
}

/*///////////////////////////////////////////////////////
func main() {
	file, err := os.Open("../conf/SrvMapping.xml")
	defer file.Close()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	var sm ServerMapping
	if err := xml.Unmarshal([]byte(data), &sm); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("-------------------------------------")
	fmt.Println("ServerMampping's server:")
	for i := 0; i < len(sm.Servers); i++ {
		fmt.Printf("The %d server is:%s\n", i+1, sm.Servers[i])
	}

	fmt.Println("-------------------------------------")
	file, err = os.Open("../conf/SysInfo.xml")
	defer file.Close()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	data, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Println("SysInfo's info:")

	var si SysInfo
	if err := xml.Unmarshal([]byte(data), &si); err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(si.Machines); i++ {
		fmt.Printf("The %d machine is %s\n", i+1, si.Machines[i])
	}

	fmt.Println("-------------------------------------")
	fmt.Println("SysConfig's info:")
	file, err := os.Open("../conf/SysConfig.xml")
	defer file.Close()
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	data, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	var sc SysConfig
	if err := xml.Unmarshal([]byte(data), &sc); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The OneMapHome is:%s\n", sc.OneMapHome)
	for i := 0; i < len(sc.LayOut.Servers); i++ {
		fmt.Printf("The %d Server is:%s\n", i+1, sc.LayOut.Servers[i])
	}

	fmt.Println("-------------------------------------")
	fmt.Printf("The FileMap is:%s\n", sc.FileMap)

}
*/
