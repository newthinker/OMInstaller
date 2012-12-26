package main

import (
	//	"container/list"
	"encoding/xml"
	//	"errors"
	"fmt"
	"io/ioutil"
	"os"

//	"strconv"
)

///////////////////////////////////////////////////////
type ServerInfo struct {
	Name  string     `xml:"name,attr"`
	Attrs []AttrInfo `xml:"attr"`
}

type AttrInfo struct {
	Name    string `xml:"name,attr"`
	Desc    string `xml:"desc,attr"`
	Encrypt string `xml:"encrypt,attr"`
	Select  string `xml:"selects,attr"`
	Value   string `xml:"value,attr"`
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

/// 目前只解析到服务器类型这一层 
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
func main() {
	fmt.Println("SysConfig's info:")
	file, err := os.Open("../conf/test.xml")
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

	var sc SysConfig
	if err := xml.Unmarshal([]byte(data), &sc); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(sc.LayOut)

	fmt.Println("-------------------------------------")
	newsc := &SysConfig{}
	newsc.LayOut = Layout{}
	newsc.FileMap = Filemap{}

	for i := 0; i < len(sc.LayOut.Servers); i++ {
		newsc.LayOut.Servers = append(newsc.LayOut.Servers, sc.LayOut.Servers[i])
	}

	for i := 0; i < len(sc.FileMap.Containers); i++ {
		newsc.FileMap.Containers = append(newsc.FileMap.Containers, sc.FileMap.Containers[i])
	}

	fmt.Println(*newsc)

	fmt.Println("-------------------------------------")
	output, err := xml.MarshalIndent(newsc, "  ", "   ")
	os.Stdout.Write([]byte(output))
	if err != nil {
		fmt.Println(err)
		return
	}
}
