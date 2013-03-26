package sys

import (
	"errors"
	"fmt"
	"github.com/newthinker/onemap-installer/log"
	"github.com/newthinker/onemap-installer/utl"
	"os"
	"path/filepath"
	"runtime"
)

// process status
const (
	BEGIN = iota // 0, process is running
	END   = 100  // 100, process is end
	BREAK = -1   // -1, process is abnormal

	PREPARE = 20 // rate of the preparing
	PROCESS = 70 // rate of the main process
	CLEAN   = 10 // rate of cleaning up work

	// prepare
	GET_JSON       = 5
	GET_WORKINGDIR = 5
	PARSE_JSON     = 5
	MAIN_PROCESS   = 5

	// process
	CHECK_WORKINGDIR  = 5
	GET_VERSION       = 5
	GET_CONTAINER     = 5
	GET_INFO          = 10
	PACKAGE           = 5
	REFRESH_SYSCONFIG = 10
	// install
	UPDATE_SCRIPT = 5
	REMOTE_COPY   = 15
	REMOTE_EXEC   = 10
	/// update 
	// UPDATE_SCRIPT = 5
	MOUNT_SSHFS = 5
	REMOTE_DIFF = 10
	// REMOTE_EXEC = 10

	/// uninstall 
	REMOTE_STANDALONE = 45

	// sysdeploy
	REFRESH_SYSDEPLOY = 10

	///////////////////////////////////////////////////////
	MAX_POOL_SIZE = 4 // message queue's max size
)

var (
	l       *(log.Logger)
	mc      chan Result // message chan
	SubFlag bool        // whether install subplatform module
	CurOS   string      // current machine's os	["windows"|"linux"]
	Deploy  int         // deploy status
)

// return result
type Result struct {
	Ret    int
	Reason string
	Data   interface{}
}

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
	filename := filepath.FromSlash(basedir + "/conf/" + SERVER_MAPPING)
	l.Debugf("SrvMapping file:%s", filename)
	sm, err := OpenSMConfig(filename)
	if err != nil {
		l.Error(errors.New("Parse SrvMapping config files failed"))
		return errors.New("Parse SrvMapping config files failed")
	}
	filename = filepath.FromSlash(basedir + "/conf/" + SYS_CONFIG)
	l.Debugf("SysConfig file:%s", filename)
	sc, err := OpenSCConfig(filename)
	if err != nil {
		l.Error(errors.New("Parse SysConfig config files failed"))
		return errors.New("Parse SysConfig config files failed")
	}
	filename = filepath.FromSlash(basedir + "/conf/" + SYS_DEPLOY)
	l.Debugf("SysDeploy file:%s", filename)
	sd, err := OpenSDConfig(filename)
	if err != nil {
		l.Error(errors.New("Parse SysDeploy config files failed"))
		return errors.New("Parse SysDeploy config files failed")
	}

	omsc = &sc
	omsm = &sm
	omsd = &sd

	// message queue
	mc = make(chan Result, MAX_POOL_SIZE)

	// default installing subplatform module
	SubFlag = true

	// os type
	CurOS = runtime.GOOS

	// deploy status, default 
	Deploy = MAINTAIN

	return nil
}

// the main process
func Process(sd SysDeploy, arr_lo []Layout) error {
	flag := true    // process status
	rate := PREPARE // main process's initial rate

	for i := 0; i < len(sd.Nodes); i++ { // from one to one
		var mac *node = &(sd.Nodes[i])
		var lo *Layout = &(arr_lo[i])
		var num int = len(sd.Nodes)

		ip, err := GetNodeIP(mac)
		if err != nil || ip == "" {
			l.Errorf("Get current node's IP failed")
			continue
		}

		msg := fmt.Sprintf("Begin to process the %dth machine and ip is:%s", i+1, ip)
		l.Message(msg)
		go FormatResult(rate, utl.NowString()+msg, nil)

		rate += CHECK_WORKINGDIR / num
		go FormatResult(rate, utl.NowString()+"Check the working directory", nil)
		var om OMPInfo
		if flag := utl.Exists(basedir); flag != true {
			err := errors.New("The working directory isn't existed")
			l.Error(err)
			flag = false
			go FormatResult(rate, utl.NowString()+"Check the working directory failed", nil)

			// Goto process the next machine
			rate = PREPARE + PROCESS/num
			mac.ResetSysDeploy(INSTALL)
			continue
		}

		// get OneMap's version
		rate += GET_VERSION / num
		go FormatResult(rate, utl.NowString()+"Get OneMap's version info", nil)
		filename, err := om.OMGetVersion(basedir)
		if err != nil {
			l.Errorf("Get OneMap version failed")
			flag = false
			go FormatResult(rate, utl.NowString()+"Get OneMap's version info failed", nil)

			//			goto Unexpected
			rate = PREPARE + PROCESS/num
			mac.ResetSysDeploy(INSTALL)
			continue
		}

		// get the container
		rate += GET_CONTAINER / num
		go FormatResult(rate, utl.NowString()+"Get OneMap's container", nil)
		if err = om.OMGetContainer(basedir, filename); err != nil {
			l.Errorf("Get OneMap container failed")
			flag = false
			go FormatResult(rate, utl.NowString()+"Get OneMap's container failed", nil)

			//goto Unexpected
			rate = PREPARE + PROCESS/num
			mac.ResetSysDeploy(INSTALL)
			continue
		}

		// get the info of the current machine
		rate += GET_INFO / num
		go FormatResult(rate, utl.NowString()+"Get OneMap package's info", nil)
		l.Messagef("Get the %dth machine's info", i+1)
		if err := om.OMGetInfo(mac, omsm, lo); err != nil {
			l.Errorf("Get the %dth machine's info failed", i+1)
			flag = false
			go FormatResult(rate, utl.NowString()+"Get OneMap package's info failed", nil)

			rate = PREPARE + PROCESS/num
			mac.ResetSysDeploy(INSTALL)
			continue
		}

		status := om.Deploy // current node's deploy status
		// package onemap
		rate += PACKAGE / num
		go FormatResult(rate, utl.NowString()+"Package the OneMap", nil)
		if err := om.OMPackage(); err != nil {
			l.Error(err)
			flag = false
			go FormatResult(rate, utl.NowString()+"Package the OneMap failed", nil)

			rate = PREPARE + PROCESS/num
			mac.ResetSysDeploy(INSTALL)
			continue
		}

		// refresh the SysConfig.xml file
		rate += REFRESH_SYSCONFIG / num
		go FormatResult(rate, utl.NowString()+"Refresh the SysConfig file", nil)
		l.Message("Update the local config file")
		srcfile := filepath.FromSlash(basedir + "/" + om.MacName + "_" + ONEMAP_NAME + "/config/SystemConfig/SysConfig.xml")
		omsc.LayOut = *lo
		if err := RefreshSysConfig(omsc, srcfile); err != nil {
			l.Error(err)
			flag = false
			go FormatResult(rate, utl.NowString()+"Refresh the SysConfig file failed", nil)

			rate = PREPARE + PROCESS/num
			mac.ResetSysDeploy(INSTALL)
			continue
		}

		l.Debugf("status:%d\n", status)
		if status == MAINTAIN { // do nothing
			continue
		} else if status == INSTALL { // install process
			// update the installing script
			rate += UPDATE_SCRIPT / num
			go FormatResult(rate, utl.NowString()+"Update standalone install script", nil)
			l.Message("Update the local install script")
			switch CurOS {
			case "windows":
				srcfile = filepath.FromSlash(basedir + "/" + om.MacName + "_" + ONEMAP_NAME + "/install.bat")
			case "linux":
				srcfile = filepath.FromSlash(basedir + "/" + om.MacName + "_" + ONEMAP_NAME + "/install.sh")
			}
			if err := UpdateScritp(&om, srcfile); err != nil {
				l.Errorf("Update the local install script failed")
				flag = false
				go FormatResult(rate, utl.NowString()+"Update standalone install script failed", nil)

				rate = PREPARE + PROCESS/num
				mac.ResetSysDeploy(INSTALL)
				continue
			}

			// remote copy the OneMap package
			srcdir := filepath.FromSlash(basedir + "/" + om.MacName + "_" + ONEMAP_NAME) /// 目前不考虑异构平台
			dstdir := filepath.FromSlash(om.OMHome)

			rate += REMOTE_COPY / num
			go FormatResult(rate, utl.NowString()+"Remote copy the OneMap package", nil)
			l.Message("Exec the remote copy")
			if err := om.OMRemoteCopy(srcdir, dstdir); err != nil {
				l.Errorf("Exec retmote copy failed")
				flag = false
				go FormatResult(rate, utl.NowString()+"Remote copy the OneMap package failed", nil)

				rate = PREPARE + PROCESS/num
				mac.ResetSysDeploy(INSTALL)
				continue
			}

			// remote exec the install bash script
			rate += REMOTE_EXEC / num
			go FormatResult(rate, utl.NowString()+"Remote exec the standalone install script", nil)
			l.Message("Exec the remote install script")
			if err := om.OMRemoteExec(); err != nil {
				l.Errorf("Exec retmote command failed")
				flag = false
				go FormatResult(rate, utl.NowString()+"Remote exec the standalone install script failed", nil)

				rate = PREPARE + PROCESS/num
				mac.ResetSysDeploy(INSTALL)
				continue
			}

			flag = true
		} else if status == UPDATE { /// update 
			// package onemap
			rate += PACKAGE / num
			go FormatResult(rate, utl.NowString()+"Package the OneMap", nil)
			if err := om.OMPackage(); err != nil {
				l.Error(err)
				flag = false

				rate = PREPARE + PROCESS/num
				mac.ResetSysDeploy(INSTALL)
				continue
			}

			// refresh the SysConfig.xml file
			rate += REFRESH_SYSCONFIG / num
			go FormatResult(rate, utl.NowString()+"Refresh the SysConfig file", nil)
			l.Message("Update the local config file")
			srcfile := filepath.FromSlash(basedir + "/" + om.MacName + "_" + ONEMAP_NAME + "/config/SystemConfig/SysConfig.xml")
			omsc.LayOut = *lo
			if err := RefreshSysConfig(omsc, srcfile); err != nil {
				l.Error(err)
				flag = false

				rate = PREPARE + PROCESS/num
				mac.ResetSysDeploy(INSTALL)
				continue
			}

			/// remote mount sshfs

			/// exec the remote script diff the two directory,
			/// parse the result and merge the two directory

			/// remote exec the sysconfig process and restart the services

		} else if status == UNINSTALL { /// uninstall
			/// remote exec the standalone uninstall script(include uninstall the services
			/// and delete the directory)
		}
		///////////////////////////////////////////////////////////////////////////////////
		rate += CLEAN / num
		go FormatResult(rate, utl.NowString()+"Remove the temp directory", nil)
		// delete the temp directory
		if err := os.RemoveAll(filepath.FromSlash(basedir + "/" + om.MacName + "_" + ONEMAP_NAME)); err != nil {
			l.Warningf("Remove the temp directory failed and please do it manually")
		}
	}

	// refresh the SysDeploy.xml config file
	filename := filepath.FromSlash(basedir + "/conf/" + SYS_DEPLOY)
	if err := RefreshSysDeploy(&sd, filename); err != nil {
		l.Errorf("Save the SysDeploy config file failed")
		return err
	}

	if flag == false {
		if Deploy == INSTALL {
			err := errors.New("OneMap install failed")
			l.Error(err)
			return err
		} else if Deploy == UPDATE {
			err := errors.New("OneMap update failed")
			l.Error(err)
			return err
		} else if Deploy == UNINSTALL {
			err := errors.New("OneMap uninstall failed")
			l.Error(err)
			return err
		}
	} else {
		if status == INSTALL {
			l.Message("OneMap install successfully")
		} else if status == UPDATE {
			l.Message("OneMap update successfully")
		} else if status == UNINSTALL {
			l.Message("OneMap uninstall successfully")
		}
	}

	return nil
}
