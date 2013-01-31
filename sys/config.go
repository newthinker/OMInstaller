package sys

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/newthinker/onemap-installer/utl"
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

////////////////////////////////////////////////////////
// SysConfig struct
type SysConfig struct {
	XMLName    xml.Name `xml:"root"`
	OneMapHome string   `xml:"OneMapHome,attr"`
	LayOut     Layout   `xml:"layout"`
	FileMap    Filemap  `xml:"filemapping"`
}

type Layout struct {
	Servers []ServerInfo `xml:"server"`
}

type ServerInfo struct {
	Srvname string     `xml:"name,attr"`
	Srvdesc string     `xml:"desc,attr"`
	Attrs   []AttrInfo `xml:"attr"`
}

type AttrInfo struct {
	Attrname  string `xml:"name,attr"`
	Attrdesc  string `xml:"desc,attr"`
	Encrypt   string `xml:"encrypt,attr"`
	Select    string `xml:"selects,attr"`
	Attrvalue string `xml:"value,attr"`
}

type Filemap struct {
	Containers []Container `xml:"container"`
}

type Container struct {
	Name    string      `xml:"name,attr"`
	Path    string      `xml:"path,attr"`
	Modules []ModuleMap `xml:"module"`
}

type ModuleMap struct {
	Name       string      `xml:"name,attr"`
	ServersMap []ServerMap `xml:"server"`
}

type ServerMap struct {
	Name    string `xml:"name,attr"`
	FileSet []File `xml:"file"`
}

type File struct {
	Path   string `xml:"path,attr"`
	KeySet []Key  `xml:"key"`
}

type Key struct {
	Template  string `xml:"template,attr"`
	Attribute string `xml:"attribute,attr"`
	Value     string `xml:",chardata"`
}

////////////////////////////////////////////////////////
// SysDeploy struct
type SysDeploy struct {
	XMLName xml.Name `xml:"root"`
	Nodes   []node   `xml:"nodes->node"`
}

type node struct {
	Nodename string `xml:"name,attr"`
	Attrs    []attr `xml:"attrs->attr"`
	Srvs     []srv  `xml:"servers->server"`
}

type attr struct {
	Attrname  string `xml:"name,attr"`
	Attrvalue string `xml:",chardata"`
}

type srv struct {
	Srvname string `xml:"name,attr"`
}

////////////////////////////////////////////////////////
/*/ update SysConfig.xml file except MonitorAgent module
func UpdateConfig(si *SysInfo, sc *SysConfig) error {
	flag := make(map[string]bool) // flag of whether update
	// initialize
	for i := 0; i < len(sc.LayOut.Servers); i++ {
		var srvinfo *ServerInfo = &(sc.LayOut.Servers[i])
		var srvtype string = srvinfo.Srvname
		flag[srvtype] = false
	}

	for i := 0; i < len(si.Machines); i++ {
		for j := 0; j < len(si.Machines[i].Servers); j++ {
			var si_srvinfo *ServerInfo = &(si.Machines[i].Servers[j])
			var si_srvtype string = si_srvinfo.Srvname // server type

			// don't set MonitorAgent server temperarily
			if si_srvtype == "agent" {
				flag[si_srvtype] = true
				continue
			}

			for k := 0; k < len(sc.LayOut.Servers); k++ {
				var sc_srvinfo *ServerInfo = &(sc.LayOut.Servers[k])
				var sc_srvtype string = sc_srvinfo.Srvname

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
		//        l.Warningf("%d modules(%s) didnot updated", strconv.Itoa(num), nomodules)
	}

	return nil
}

// update SysConfig.xml file's MonitorAgent module
func UpdateMdlAgent(mi *MachineInfo, sc *SysConfig) error {
	var flag bool = false
	for i := 0; i < len(mi.Servers); i++ {
		var mi_srvinfo *ServerInfo = &(mi.Servers[i])
		var mi_srvtype string = mi_srvinfo.Srvname

		if mi_srvtype != "agent" {
			continue
		} else {
			for j := 0; j < len(sc.LayOut.Servers); j++ {
				var sc_srvinfo *ServerInfo = &(sc.LayOut.Servers[j])
				var sc_srvtype string = sc_srvinfo.Srvname

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
}*/

// open the config files
func OpenSMConfig(filename string) (ServerMapping, error) {
	var sm ServerMapping

	// check the config file whether existed
	//	filename := basedir + "/conf/" + SERVER_MAPPING
	if flag := utl.Exists(filename); flag != true {
		err := fmt.Errorf("Config file(%s) isn't existed", filename)
		l.Error(err)
		return sm, err
	}

	file, err := os.Open(filename)
	if err != nil {
		l.Error(err)
		return sm, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		l.Error(err)
		return sm, err
	}
	if err := xml.Unmarshal([]byte(data), &sm); err != nil {
		l.Error(err)
		return sm, err
	}

	return sm, nil
}

/*
func OpenSIConfig(filename string) (SysInfo, error) {
	var si SysInfo

	// check the config file whether existed
	if flag := utl.Exists(filename); flag != true {
		err := fmt.Errorf("Config file(%s) isn't existed", filename)
		l.Error(err)
		return sm, err
	}

	file, err1 := os.Open(filename)
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
} */

func OpenSCConfig(filename string) (SysConfig, error) {
	var sc SysConfig

	// check the config file whether existed
	if flag := utl.Exists(filename); flag != true {
		err := fmt.Errorf("Config file(%s) isn't existed", filename)
		l.Error(err)
		return sc, err
	}

	file, err1 := os.Open(filename)
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

func OpenSDConfig(filename string) (SysDeploy, error) {
	var sd SysDeploy

	// check the config file whether existed
	if flag := utl.Exists(filename); flag != true {
		err := fmt.Errorf("Config file(%s) isn't existed", filename)
		l.Error(err)
		return sd, err
	}

	file, err1 := os.Open(filename)
	if err1 != nil {
		return sd, err1
	}
	data, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		return sd, err2
	}
	if err3 := xml.Unmarshal([]byte(data), &sd); err3 != nil {
		return sd, err3
	}

	return sd, nil
}

// 更新保存系统配置文件
func RefreshSysConfig(sc *SysConfig, conffile string) error {
	if sc == nil || conffile == "" {
		return errors.New("输入的配置文件为空、文件路径为空")
	}

	if flag := utl.Exists(conffile); flag != true {
		return errors.New("配置文件路径不存在")
	}

	output, err := xml.MarshalIndent(sc, "  ", "   ")
	if err != nil {
		return err
	}

	if utl.Exists(conffile) == true {
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

// Update or save SysDeploy.xml config file 
func RefreshSysDeploy(sd *SysDeploy, conffile string) error {
	if sd == nil || conffile == "" {
		return errors.New("Config file is null or file path is null")
	}

	if utl.Exists(conffile) == true {
		l.Warning("Config file is existed and delete it first")
		if err := os.Remove(conffile); err != nil {
			l.Errorf("Remove config file(%s) failed", conffile)
			return err
		}
	}

	output, err := xml.MarshalIndent(sd, "  ", "   ")
	if err != nil {
		return err
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

/// Reset the attribute of SysDeploy structure
func (sd *SysDeploy) ResetSysDeploy(newvalue string, tags []string) error {
	if len(tags) < 0 {
		l.Warning("No attribute will be reseted")
		return nil
	}

	return nil
}

/// Reset the deploy attribute of SysDeploy structure
func (n *node) ResetSysDeploy(newvalue int) {
	for i, v := range n.Attrs {
		if v.Attrname == "deploy" {
			n.Attrs[i].Attrvalue = strconv.Itoa(newvalue)
		}
	}
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
