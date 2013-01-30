package sys

import (
	"errors"
	"github.com/newthinker/onemap-installer/log"
	"github.com/newthinker/onemap-installer/utl"
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
    filename := basedir + "/conf/" + SERVER_MAPPING
	sm, err := OpenSMConfig(filename)
	if err != nil {
		l.Error(errors.New("Parse SrvMapping config files failed"))
		return errors.New("Parse SrvMapping config files failed")
	}
    filename = basedir + "/conf" + SYS_CONFIG
	sc, err := OpenSCConfig(filename)
	if err != nil {
		l.Error(errors.New("Parse SysConfig config files failed"))
		return errors.New("Parse SysConfig config files failed")
	}
	filename = basedir + "/conf/" + SYS_DEPLOY
	sd, err := OpenSDConfig(filename)
	if err != nil {
		l.Error(errors.New("Parse SysDeploy config files failed"))
		return errors.New("Parse SysDeploy config files failed")		
	}
	
	omsc = &sc
	omsm = &sm
	omsd = &sd

	return nil
}

// the main process
func Process(sd SysDeploy, arr_lo []Layout) error {
	l.Message("Begin installing process")

	for i := 0; i < len(sd.Nodes); i++ { // from one to one
		l.Messagef("Begin process the %d machine", i+1)

		var mac *node = &(sd.Nodes[i])
		var lo *Layout = &(arr_lo[i])

		flag := true

	Unexpected:
		if flag == false {
			mac.ResetSysDeploy(INSTALL)
			continue
		}

		var om OMPInfo
		if flag := utl.Exists(basedir); flag != true {
			err := errors.New("Get working directory failed")
			l.Error(err)
			flag = false
			goto Unexpected
		}

		// get OneMap's version
		filename, err := om.OMGetVersion(basedir)
		if err != nil {
			l.Errorf("Get OneMap version failed")
			flag = false
			goto Unexpected
		}

		// get the container
		if err = om.OMGetContainer(basedir, filename); err != nil {
			l.Errorf("Get OneMap container failed")
			flag = false
			goto Unexpected
		}

		// get the info of the current machine
		l.Messagef("Get the %dth machine's info", i+1)
		if err := om.OMGetInfo(mac, omsm, lo); err != nil {
			l.Errorf("Get the %dth machine's info failed", i+1)
			flag = false
			goto Unexpected
		}

		// package onemap
		if err := om.OMPackage(); err != nil {
			l.Error(err)
			flag = false
			goto Unexpected
		}

		// refresh the SysConfig.xml file
		l.Message("Update the local config file")
		srcfile := basedir + "/" + ONEMAP_NAME + "/config/SystemConfig/SysConfig.xml"
		omsc.LayOut = *lo
		if err := RefreshSysConfig(omsc, srcfile); err != nil {
			l.Error(err)
			flag = false
			goto Unexpected
		}

        status := om.Deploy
        if status==MAINTAIN {       // do nothing
            continue
        } else if status==INSTALL { // install process
            // update the installing script
            l.Message("Update the local install script")
            srcfile = basedir + "/" + ONEMAP_NAME + "/install.sh"
            if err := UpdateScritp(&om, srcfile); err != nil {
                l.Errorf("Update the local install script failed")
                flag = false
                goto Unexpected
            }

            // remote copy the OneMap package
            srcdir := basedir + "/" + ONEMAP_NAME
            dstdir := om.OMHome

            l.Message("Exec the remote copy")
            if err := om.OMRemoteCopy(srcdir, dstdir); err != nil {
                l.Errorf("Exec retmote copy failed")
                flag = false
                goto Unexpected
            }

            // remote exec the install bash script
            l.Message("Exec the remote install script")
            if err := om.OMRemoteExec(); err != nil {
                l.Errorf("Exec retmote command failed")
                flag = false
                goto Unexpected
            }
        } else if status==UPDATE {  /// update 
            /// remote mount sshfs

            /// exec the remote script diff the two directory,
            /// parse the result and merge the two directory

            /// remote exec the sysconfig process and restart the services

        } else if status==UNINSTALL {   /// uninstall
            /// remote exec the standalone uninstall script(include uninstall the services
            /// and delete the directory)
        }
	}

	l.Message("OneMap installing successfully")

	// refresh the SysDeploy.xml config file
	filename := basedir + "/conf/" + SYS_DEPLOY
	if err := RefreshSysDeploy(&sd, filename); err != nil {
		l.Errorf("Save the SysDeploy config file failed")
		return err
	}

	return nil
}
