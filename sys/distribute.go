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
		fmt.Println("ERROR: Update system config params failed")
		return err
	}

	// 打包OenMap安装包并进行分发
	for i := 0; i < len(si.Machines); i++ {
		fmt.Printf("Distribute installing the %d machine：\n", i+1)

		var om OMPInfo
		om.Basedir = basedir

		// get the info of the current machine
		var mi *(MachineInfo) = &(si.Machines[i])
		fmt.Printf("MSG: Get the %dth machine's info\n", i+1)
		if err = om.OMGetInfo(mi, sm); err != nil {
			fmt.Printf("ERROR: Get the %d machine's info failed!", i+1)
			return err
		}

		// 打包onemap
		// 如果不存在就创建OneMap文件夹
		dstdir := basedir + "/" + ONEMAP_NAME
		fmt.Println("MSG: Make OneMap directory")
		if flag := Exists(dstdir); flag == true {
			// 首先删除原来的
			if err = os.RemoveAll(dstdir); err != nil {
				fmt.Println("ERROR: Remove the old onemap package failed")
				return err
			}
			// 再创建新的空文件夹
			if err = os.Mkdir(dstdir, 0755); err != nil {
				fmt.Println("ERROR: Make OneMap directory failed!")
				return err
			}
		}

		fmt.Println("MSG: Package onemap")
		srcdir := om.Basedir + "/" + ONEMAP_NAME + "_Linux_" + om.Version
		if err = om.OMCopy(srcdir, dstdir); err != nil {
			fmt.Println("ERROR: Package onemap failed!")
			return err
		}

		// 更新monitoragent模块
		fmt.Printf("MSG: Update the %d machine's monitoragent module\n", i+1)
		if err = UpdateMdlAgent(mi, sc); err != nil {
			fmt.Printf("ERROR: Update the %d machine's monitoragent module failed!\n", i+1)
			return err
		}

		// 将配置参数写入配置文件中
		fmt.Println("MSG: Update the local config file")
		srcfile := om.Basedir + "/" + ONEMAP_NAME + "/config/SystemConfig/SysConfig.xml"
		if err = RefreshSysConfig(sc, srcfile); err != nil {
			return err
		}

		// 更新脚本文件
		fmt.Println("MSG: Update the local install script")
		srcfile = om.Basedir + "/" + ONEMAP_NAME + "/install.sh"
		if err = UpdateScritp(&om, srcfile); err != nil {
			fmt.Println("ERROR: Update the local install script failed")
			return err
		}

		// 远程拷贝OneMap package
		srcdir = om.Basedir + "/" + ONEMAP_NAME
		dstdir = om.OMHome

		fmt.Println("MSG: Exec the remote copy")
		if err = om.OMRemoteCopy(srcdir, dstdir); err != nil {
			fmt.Println("ERROR: Exec retmote copy failed!")
			return err
		}

		// remote exec the install bash script
		fmt.Println("MSG: Exec the remote install script")
		if err = om.OMRemoteExec(); err != nil {
			fmt.Println("ERROR: Exec retmote command failed!")
			return err
		}

		fmt.Println("MSG: Distribute installing successfully")
	}

	return nil
}
