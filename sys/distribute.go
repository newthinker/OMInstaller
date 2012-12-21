package sys

import (
	"fmt"
	"os"
)

// 进行分布式安装
func Distribute(basedir string, si *SysInfo, sc *SysConfig, sm *ServerMapping) error {
	// 读取SysInfo.xml文件并进行系统参数更新
	err := UpdateConfig(si, sc)
	if err != nil {
		fmt.Println("更新系统配置参数失败")
		return err
	}

	// 打包OenMap安装包并进行分发
	for i := 0; i < len(si.Machines); i++ {
		var om OMPInfo
		om.Basedir = basedir

		// get the info of the current machine
		var mi *(MachineInfo) = &(si.Machines[i])
		if err = om.OMGetInfo(mi, sm); err != nil {
			fmt.Printf("ERROR: Get the %d machine's info failed!", i+1)
			return err
		}

		// package the onemap
		// create the onemap directory first
		dstdir := basedir + "/" + ONEMAP_NAME
		if flag := Exists(dstdir); flag != true { // create the onemap directory first
			if err = os.Mkdir(dstdir, 0755); err != nil {
				fmt.Println("ERROR: Make OneMap directory failed!")
				return err
			}
		}
		srcdir := om.Basedir + "/" + ONEMAP_NAME + "_Linux_" + om.Version
		if err = om.OMCopy(srcdir, dstdir); err != nil {
			fmt.Println("ERROR: Package onemap failed!")
			return err
		}

		// update the monitoragent module
		if err = UpdateMdlAgent(mi, sc); err != nil {
			fmt.Printf("ERROR: Update the %d machine's monitoragent module failed!\n", i+1)
			return err
		}

		// remote copy OneMap package
		srcdir = om.Basedir + "/" + ONEMAP_NAME
		dstdir = om.OMHome
		///////////////////test//////////////////////////
		//om.Ip = "192.168.80.60"
		//om.User = "root"
		//om.Pwd = "dasiyebushuo"
		/////////////////////////////////////////////////
		if err = om.OMRemoteCopy(srcdir, dstdir); err != nil {
			fmt.Println("ERROR: Exec retmote copy failed!")
			return err
		}

		// remote exec the install bash script
		if err = om.OMRemoteExec(); err != nil {
			fmt.Println("ERROR: Exec retmote command failed!")
			return err
		}
	}

	return nil
}
