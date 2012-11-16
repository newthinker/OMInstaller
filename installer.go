package main

import (
	"encoding/xml"
	"fmt"
	"github.com/newthinker/onemap-installer/onemap"
	"io/ioutil"
	"os"
	"path/filepath"
)

var sm onemap.ServerMapping
var si onemap.SysInfo
var sc onemap.SysConfig

// open the config files
func openconfigs(basedir string) int {
	// check the base dir whether existed
	if flag := onemap.Exists(basedir); flag != true {
		fmt.Printf("ERROR: The input directory(%s) isn't existed!\n", basedir)
		return 1
	}

	file, err := os.Open(basedir + "/conf/" + onemap.SERVER_MAPPING)
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

	file, err = os.Open(basedir + "/conf/" + onemap.SYS_INFO)
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

	file, err = os.Open(basedir + "/conf/" + onemap.SYS_CONFIG)
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

func main() {
	//	var curdir string // current working directory

	/// start web service and get the input params into configs

	// parse the configs
	var basedir string
	var err error
	if basedir, err = filepath.Abs("./"); err != nil || basedir == "" {
		fmt.Println("ERROR: Current directory is invalid!")
		return
	}
	fmt.Printf("MSG: Current working dirctory is: %s\n", basedir)

	if ret := openconfigs(basedir); ret != 0 {
		fmt.Println("ERROR: Open and parse configs failed and exit!")
		return
	}

	// update the params except monitoragent module
	ret := onemap.UpdateConfig(&si, &sc)
	if ret != 0 {
		fmt.Println("ERROR: Update system config failed!")
		return
	}

	// package the OneMap installer package
	for i := 0; i < len(si.Machines); i++ {
		var om onemap.OMPInfo
		om.Basedir = basedir

		// get the info of the current machine
		var mi *(onemap.MachineInfo) = &(si.Machines[i])
		ret = om.OMGetInfo(mi, &sm)
		if ret != 0 {
			fmt.Printf("ERROR: Get the %d machine's info failed!", i+1)
			return
		}

		// package the onemap
		// create the onemap directory first
		dstdir := basedir + "/" + onemap.ONEMAP_NAME
		if flag := onemap.Exists(dstdir); flag != true { // create the onemap directory first
			fmt.Println("WARN: OneMap directory isn't existed!")
			if err := os.Mkdir(dstdir, 0755); err != nil {
				fmt.Println("ERROR: Make OneMap directory failed!")
				return
			}
		}
		srcdir := om.Basedir + "/" + onemap.ONEMAP_NAME + "_Linux_" + om.Version
		ret = om.OMCopy(srcdir, dstdir)
		if ret != 0 {
			fmt.Println("ERROR: Package onemap failed!")
			return
		}

		// update the monitoragent module
		ret = onemap.UpdateMdlAgent(mi, &sc)
		if ret != 0 {
			fmt.Printf("ERROR: Update the %d machine's monitoragent module failed!\n", i+1)
			return
		}

		// remote copy OneMap package
		srcdir = om.Basedir + "/" + onemap.ONEMAP_NAME
		dstdir = om.OMHome
		///////////////////test//////////////////////////
		om.Ip = "192.168.80.60"
		om.User = "root"
		om.Pwd = "dasiyebushuo"
		/////////////////////////////////////////////////
		ret = om.OMRemoteCopy(srcdir, dstdir)

		// remote exec the install bash script
		cmdline := "/bin/bash " + om.OMHome + "/install.sh"
		for i := range om.Servers {
			if i == 0 {
				cmdline += " "
			} else if i > 0 {
				cmdline += "|"
			}

			cmdline += om.Servers[i]
		}

		ret = om.OMRemoteExec(cmdline)
		if ret != 0 {
			fmt.Println("ERROR: Exec retmote command failed!")
			return
		}
	}
}
