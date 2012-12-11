package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/newthinker/onemap-installer/sys"
	"github.com/newthinker/onemap-installer/web"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var sm sys.ServerMapping
var si sys.SysInfo
var sc sys.SysConfig

// open the config files
func openconfigs(basedir string) error {
	// check the base dir whether existed
	if flag := sys.Exists(basedir); flag != true {
		msg := "ERROR: The input directory(" + basedir + ") isn't existed!"
		return errors.New(msg)
	}

	file, err := os.Open(basedir + "/conf/" + sys.SERVER_MAPPING)
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

	file, err = os.Open(basedir + "/conf/" + sys.SYS_INFO)
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

	file, err = os.Open(basedir + "/conf/" + sys.SYS_CONFIG)
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
	////////////////////////////////////////////////////////////////
	fmt.Println("web test")

	http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.FileServer(http.Dir("template")))
	http.Handle("/img/", http.FileServer(http.Dir("template")))

	http.HandleFunc("/subconfig", web.SubHandler)
	err := http.ListenAndServe("192.168.80.98:8888", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
		return
	}

	////////////////////////////////////////////////////////////////

	//	var curdir string // current working directory

	/// start web service and get the input params into configs

	// parse the configs
	var basedir string
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
	if err = sys.UpdateConfig(&si, &sc); err != nil {
		fmt.Println("ERROR: Update system config failed!")
		return
	}

	// package the OneMap installer package
	for i := 0; i < len(si.Machines); i++ {
		var om sys.OMPInfo
		om.Basedir = basedir

		// get the info of the current machine
		var mi *(sys.MachineInfo) = &(si.Machines[i])
		if err = om.OMGetInfo(mi, &sm); err != nil {
			fmt.Printf("ERROR: Get the %d machine's info failed!", i+1)
			return
		}

		// package the onemap
		// create the onemap directory first
		dstdir := basedir + "/" + sys.ONEMAP_NAME
		if flag := sys.Exists(dstdir); flag != true { // create the onemap directory first
			fmt.Println("WARN: OneMap directory isn't existed!")
			if err := os.Mkdir(dstdir, 0755); err != nil {
				fmt.Println("ERROR: Make OneMap directory failed!")
				return
			}
		}
		srcdir := om.Basedir + "/" + sys.ONEMAP_NAME + "_Linux_" + om.Version
		if err = om.OMCopy(srcdir, dstdir); err != nil {
			fmt.Println("ERROR: Package onemap failed!")
			return
		}

		// update the monitoragent module
		if err = sys.UpdateMdlAgent(mi, &sc); err != nil {
			fmt.Printf("ERROR: Update the %d machine's monitoragent module failed!\n", i+1)
			return
		}

		// remote copy OneMap package
		srcdir = om.Basedir + "/" + sys.ONEMAP_NAME
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
