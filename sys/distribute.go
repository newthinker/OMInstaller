package sys

import (
	"errors"
	"github.com/newthinker/onemap-installer/log"
	"github.com/newthinker/onemap-installer/utl"

//	"os"
)

var (
	l *(log.Logger)
)

func Init(logger *(log.Logger)) error {
	l = logger
	if l == nil {
		err := errors.New("Init successfully")
		l.Error(err)
		return err
	}

	// 获取当前目录
	dir, err := utl.GetLocalDir()
	if err != nil {
		l.Error(err)
		return err
	}
	basedir = dir
	l.Debugf("Current directory is:%s", basedir)

	// open the config files
	sm, err1 := OpenSMConfig(basedir)
	sc, err2 := OpenSCConfig(basedir)
	if err1 != nil || err2 != nil {
		l.Error(errors.New("Parse system config files failed"))
		return errors.New("Parse system config files failed")
	}
	omsc = &sc
	omsm = &sm

	return nil
}

// 进行安装
func Install(sd *SysDeploy, arr_lo []Layout) error {
	// 开始进行安装
	l.Message("Begin installing process")

	for i := 0; i < len(sd.Nodes); i++ { // 逐个节点进行
		l.Messagef("Installing the %d machine", i+1)

		var om OMPInfo
		if flag := utl.Exists(basedir); flag != true {
			err := errors.New("Get working directory failed")
			l.Error(err)
			return err
		}

		// 获取OneMap版本
		filename, err := om.OMGetVersion(basedir)
		if err != nil {
			l.Errorf("Get OneMap version failed")
			return err
		}

		// get the container
		if err = om.OMGetContainer(basedir, filename); err != nil {
			l.Errorf("Get OneMap container failed")
			return err
		}

		// get the info of the current machine
		var mac *node = &(sd.Nodes[i])
		var lo *Layout = &(arr_lo[i])
		l.Messagef("Get the %dth machine's info", i+1)
		if err := om.OMGetInfo(mac, omsm, lo); err != nil {
			l.Errorf("Get the %dth machine's info failed", i+1)
		}

		// package onemap
		if err := om.OMPackage(); err != nil {
			l.Error(err)
			return err
		}

		// 将配置参数写入配置文件中
		l.Message("Update the local config file")
		srcfile := basedir + "/" + ONEMAP_NAME + "/config/SystemConfig/SysConfig.xml"
		omsc.LayOut = *lo
		if err := RefreshSysConfig(omsc, srcfile); err != nil {
			l.Error(err)
			return err
		}

		// 更新脚本文件
		l.Message("Update the local install script")
		srcfile = basedir + "/" + ONEMAP_NAME + "/install.sh"
		if err := UpdateScritp(&om, srcfile); err != nil {
			l.Errorf("Update the local install script failed")
			return err
		}

		// 远程拷贝OneMap package
		srcdir := basedir + "/" + ONEMAP_NAME
		dstdir := om.OMHome

		l.Message("Exec the remote copy")
		if err := om.OMRemoteCopy(srcdir, dstdir); err != nil {
			l.Errorf("Exec retmote copy failed")
			return err
		}

		// remote exec the install bash script
		l.Message("Exec the remote install script")
		if err := om.OMRemoteExec(); err != nil {
			l.Errorf("Exec retmote command failed")
			return err
		}
	}

	l.Message("OneMap installing successfully")

	return nil
}
