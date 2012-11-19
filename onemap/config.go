package onemap

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
)

///////////////////////////////////////////////////////
// SysMapping struct
type ServerMapping struct {
	XMLName     xml.Name `xml:"root"`
	Servers     []Server `xml:",any"`
	Description string   `xml:",innerxml"`
}

type Server struct {
	XMLName    xml.Name `xml:""`
	ModuleList []Module `xml:",any"`
}

type Module struct {
	XMLName    xml.Name `xml:""`
	ModuleName string   `xml:",chardata"`
}

func (s *Server) AddModule(mdltype string, mdlname string) {
	newm := Module{ModuleName: mdlname}
	newm.XMLName.Local = mdltype
	s.ModuleList = append(s.ModuleList, newm)
}

func (ss *ServerMapping) AddServer(srvname string, mdltype []string, mdlnames []string) {
	news := Server{}
	news.XMLName.Local = srvname
	for v := range mdlnames {
		news.AddModule(mdltype[v], mdlnames[v])
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

///////////////////////////////////////////////////////
// SysInfo struct
type SysInfo struct {
	XMLName     xml.Name      `xml:"root"`
	Machines    []MachineInfo `xml:",any"`
	Description string        `xml:",innerxml"`
}

type MachineInfo struct {
	XMLName xml.Name     `xml:"machine"`
	Os      string       `xml:"os,attr"`
	Arch    string       `xml:"arch,attr"`
	Ip      string       `xml:"ip,attr"`
	User    string       `xml:"user,attr"`
	Pwd     string       `xml:"pwd,attr"`
	Omhome  string       `xml:"omhome,attr"`
	Servers []ServerInfo `xml:",any"`
}

type ServerInfo struct {
	XMLName xml.Name   `xml:""`
	Attrs   []AttrInfo `xml:",any"`
}

type AttrInfo struct {
	XMLName  xml.Name `xml:""`
	AttrName string   `xml:"name,attr"`
	Value    string   `xml:",chardata"`
}

func (s *ServerInfo) AddAttrInfo(attrname string,
	attrvalue string) {
	newa := AttrInfo{Value: attrvalue}
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
	servers []ServerInfo, omhome string) {
	newm := MachineInfo{Os: os, Arch: arch, Ip: ip, User: user, Pwd: pwd, Omhome: omhome}
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
			nomodules = nomodules + "," + key
			num++
		}
	}
	if num > 0 && nomodules != "" {
		msg := "WARN: There are " + strconv.Itoa(num) + " modules(" + nomodules + ") not updated!"
		fmt.Println(msg)
		return errors.New(msg)
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
		return errors.New("WARN: There's no MonitorAgent module in config files!")
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
	file, err = os.Open("../conf/SysConfig.xml")
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
