package onemap

import (
	"encoding/xml"
	"fmt"

//  "errors"
)

/////////////////////////////////////////////////////////////////////////
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

/////////////////////////////////////////////////////////////////////////////
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

/////////////////////////////////////////////////////////////////////////////
// SysConfig struct
type SysConfig struct {
	XMLName    xml.Name `xml:"root"`
	OneMapHome string   `xml:"OneMapHome,attr"`
	LayOut     Layout   `xml:",any"`
}

type Layout struct {
	XMLName xml.Name     `xml:"Layout"`
	Servers []ServerInfo `xml:",any"`
}

func (lo *Layout) AddServerInfo(srvtype string, attrarray []AttrInfo) {
	news := ServerInfo{}
	news.XMLName.Local = srvtype
	news.Attrs = attrarray
	lo.Servers = append(lo.Servers, news)
}

////////////////////////////////////////////////////////
// update SysConfig.xml file except MonitorAgent module
func UpdateConfig(si *SysInfo, sc *SysConfig) int {
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
		fmt.Printf("WARN: There are %d modules(%s) not updated!", num, nomodules)
		return 1
	}

	return 0
}

// update SysConfig.xml file's MonitorAgent module
func UpdateMdlAgent(mi *MachineInfo, sc *SysConfig) int {
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
		fmt.Println("WARN: There's no MonitorAgent module in config files!")
		return 1
	}

	return 0
}

/*///////////////////////////////////////////////////////
func main() {
	file, err := os.Open("SrvMapping.xml")
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
	file, err = os.Open("SysInfo.xml")
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
	file, err = os.Open("SysConfig.xml")
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
}
*/
