package onemap

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var sm ServerMapping
var si SysInfo
var sc SysConfig

func openconfigs() int {
	file, err := os.Open(SERVER_MAPPING)
	if err != nil {
		fmt.Printf("Open SrvMapping config file failed: %v\n", err)
		return 1
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Read SrvMapping config file failed: %v\n", err)
		return 2
	}
	if err := xml.Unmarshal([]byte(data), &sm); err != nil {
		fmt.Printf("Parse SrvMapping config file failed: %v\n", err)
		return 3
	}

	file, err = os.Open(SYS_INFO)
	if err != nil {
		fmt.Printf("Open SysInfo config file failed: %v\n", err)
		return 1
	}
	data, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Read SysInfo config file failed: %v\n", err)
		return 2
	}
	if err = xml.Unmarshal([]byte(data), &si); err != nil {
		fmt.Printf("Parse SysInfo config file failed: %v\n", err)
		return 3
	}

	file, err = os.Open(SYS_CONFIG)
	if err != nil {
		fmt.Printf("Open SysConfig config file failed: %v\n", err)
		return 1
	}
	data, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Read SysConfig file failed: %v\n", err)
		return 2
	}
	if err = xml.Unmarshal([]byte(data), &sc); err != nil {
		fmt.Printf("Parse SysConfig file failed: %v\n", err)
		return 3
	}

	return 0
}

// test parse config files
func TestParseConfigs(t *testing.T) {
	fmt.Println("Test parsing SrvMapping.xml")
	for i := 0; i < len(sm.Servers); i++ {
		fmt.Printf("The %d server is: %s\n", i+1, sm.Servers[i])
	}

	fmt.Println("Test parsing SysInfo.xml")
	for i := 0; i < len(si.Machines); i++ {
		fmt.Printf("The %d machine is %s\n", i+1, si.Machines[i])
	}

	fmt.Println("Test parsing SysConfig.xml")
	for i := 0; i < len(sc.LayOut.Servers); i++ {
		fmt.Printf("The %d server is %s\n", i+1, sc.LayOut.Servers[i])
	}
}

// test update config file
//func TestUpdateConfig(t *testing.T) {

//}
