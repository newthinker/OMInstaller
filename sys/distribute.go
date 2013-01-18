package sys

import (
	"os"
	"github.com/newthinker/onemap-installer/log"
	"github.com/newthinker/onemap-installer/utl"
)

var (
    l *(log.Logger)
)

func Init(logger *(log.Logger)) {
    l = logger
    if l!=nil {
        l.Debug("Init successfully")
    }
}

// 进行分布式安装
func Distribute(basedir string, si *SysInfo, sc *SysConfig, sm *ServerMapping) error {
	// 读取SysInfo.xml文件并进行系统参数更新
    l.Message("Begin distribute installing")
	err := UpdateConfig(si, sc)
	if err != nil {
        l.Errorf("Update system config params failed")
		return err
	}

	// 打包OenMap安装包并进行分发
	for i := 0; i < len(si.Machines); i++ {
		l.Messagef("Distribute installing the %d machine", i+1)

		var om OMPInfo
		om.Basedir = basedir

		// get the info of the current machine
		var mi *(MachineInfo) = &(si.Machines[i])
		l.Messagef("Get the %dth machine's info", i+1)
		if err = om.OMGetInfo(mi, sm); err != nil {
            l.Errorf("Get the %dth machine's info failed", i+1)
			return err
		}

		// 打包onemap
		// 如果不存在就创建OneMap文件夹
		dstdir := basedir + "/" + ONEMAP_NAME
		l.Message("Make OneMap directory")
		if flag := utl.Exists(dstdir); flag == true {
			// 首先删除原来的
			if err = os.RemoveAll(dstdir); err != nil {
                l.Errorf("Remove the old OneMap package failed")
				return err
			}
			// 再创建新的空文件夹
			if err = os.Mkdir(dstdir, 0755); err != nil {
                l.Errorf("Make OneMap directory failed")
				return err
			}
		}

		l.Message("Package OneMap")
		srcdir := om.Basedir + "/" + ONEMAP_NAME + "_Linux_" + om.Version
		if err = om.OMCopy(srcdir, dstdir); err != nil {
            l.Errorf("Package OneMap failed")
			return err
		}

		// 更新monitoragent模块
		l.Messagef("Update the %d machine's monitoragent module", i+1)
		if err = UpdateMdlAgent(mi, sc); err != nil {
            l.Errorf("Update the %d machine's MonitorAgent module failed", i+1)
			return err
		}

		// 将配置参数写入配置文件中
		l.Message("Update the local config file")
		srcfile := om.Basedir + "/" + ONEMAP_NAME + "/config/SystemConfig/SysConfig.xml"
		if err = RefreshSysConfig(sc, srcfile); err != nil {
            l.Error(err)
			return err
		}

		// 更新脚本文件
		l.Message("Update the local install script")
		srcfile = om.Basedir + "/" + ONEMAP_NAME + "/install.sh"
		if err = UpdateScritp(&om, srcfile); err != nil {
			l.Errorf("Update the local install script failed")
			return err
		}

		// 远程拷贝OneMap package
		srcdir = om.Basedir + "/" + ONEMAP_NAME
		dstdir = om.OMHome

		l.Message("Exec the remote copy")
		if err = om.OMRemoteCopy(srcdir, dstdir); err != nil {
			l.Errorf("Exec retmote copy failed")
			return err
		}

		// remote exec the install bash script
		l.Message("Exec the remote install script")
		if err = om.OMRemoteExec(); err != nil {
			l.Errorf("Exec retmote command failed")
			return err
		}

		l.Message("Distribute installing successfully")
	}

	return nil
}
