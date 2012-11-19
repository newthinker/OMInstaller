package main

import (
	"encoding/xml"
	"errors"
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
func openconfigs(basedir string) error {
	// check the base dir whether existed
	if flag := onemap.Exists(basedir); flag != true {
		msg := "ERROR: The input directory(" + basedir + ") isn't existed!"
		return errors.New(msg)
	}

	file, err := os.Open(basedir + "/conf/" + onemap.SERVER_MAPPING)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if err := xml.Unmarshal([]byte(data), &sm); err != nil {
		return err
	}

	file, err = os.Open(basedir + "/conf/" + onemap.SYS_INFO)
	if err != nil {
		return err
	}
	data, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if err = xml.Unmarshal([]byte(data), &si); err != nil {
		return err
	}

	file, err = os.Open(basedir + "/conf/" + onemap.SYS_CONFIG)
	if err != nil {
		return err
	}
	data, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	if err = xml.Unmarshal([]byte(data), &sc); err != nil {
		return err
	}

	return nil
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

	if err = openconfigs(basedir); err != nil {
		fmt.Println("ERROR: Open and parse configs failed and exit!")
		return
	}

	// update the params except monitoragent module
	if err = onemap.UpdateConfig(&si, &sc); err != nil {
		fmt.Println("ERROR: Update system config failed!")
		return
	}

	// package the OneMap installer package
	for i := 0; i < len(si.Machines); i++ {
		var om onemap.OMPInfo
		om.Basedir = basedir

		// get the info of the current machine
		var mi *(onemap.MachineInfo) = &(si.Machines[i])
		if err = om.OMGetInfo(mi, &sm); err != nil {
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
		if err = om.OMCopy(srcdir, dstdir); err != nil {
			fmt.Println("ERROR: Package onemap failed!")
			return
		}

		// update the monitoragent module
		if err = onemap.UpdateMdlAgent(mi, &sc); err != nil {
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
		if err := om.OMRemoteCopy(srcdir, dstdir); err != nil {
			fmt.Println("ERROR: Exec retmote copy failed!")
			return
		}

		// remote exec the install bash script
		if err := om.OMRemoteExec(); err != nil {
			fmt.Println("ERROR: Exec retmote command failed!")
			return
		}
	}
}
