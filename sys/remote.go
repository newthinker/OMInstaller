package sys

import (
	"errors"
	"github.com/newthinker/onemap-installer/utl"
	"os"
	"os/exec"
)

// remote copy the OneMap package
func (om *OMPInfo) OMRemoteCopy(srcdir string, dstdir string) error {
	// check whether installed sshpass package
	cmd := exec.Command(basedir+"/sshpass/bin/sshpass", "-V")
	err := cmd.Run()
	if err != nil {
		return errors.New("sshpass isn't installed")
	}

	// check srcdir is a file or directory
	if flag := utl.Exists(srcdir); flag != true {
		msg := "Source file or directory " + srcdir + " isn't existed"
		return errors.New(msg)
	}

	fi, _ := os.Stat(srcdir)
	if fi.IsDir() {
		cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "scp", "-r", srcdir, om.Root+"@"+om.Ip+":"+dstdir)

		l.Debugf("sshpass -p %s scp -r %s %s@%s:%s", om.Pwd, srcdir, om.Root, om.Ip, dstdir)
	} else {
		cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "scp", srcdir, om.Root+"@"+om.Ip+":"+dstdir)

		l.Debugf("sshpass -p %s scp %s %s@%s:%s", om.Pwd, srcdir, om.Root, om.Ip, dstdir)
	}
	err = cmd.Run()
	if err != nil {
		return errors.New("Exec remote copy command failed")
	}

	return nil
}

// exec the remote command
func (om *OMPInfo) OMRemoteExec() error {
	// parse the remote command line
	if len(om.Servers) <= 0 {
		msg := "No install modules"
		return errors.New(msg)
	}

	// check whether installed sshpass package
	cmd := exec.Command(basedir+"/sshpass/bin/sshpass", "-V")
	err := cmd.Run()
	if err != nil {
		msg := "Sshpass isn't installed"
		return errors.New(msg)
	}

	// service flag
	flag_ma := true  // monitoragent service
	flag_h2 := false // h2memdb service
	flag_mq := false // activemq service
	flag_om := false // onemap service

	// exec the remote command line to install the OneMap
	for i := 0; i < len(om.Servers); i++ {
		cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
			"/bin/bash", om.OMHome+"/install.sh", om.Servers[i])
		l.Debugf("sshpass -p %s ssh %s@%s /bin/bash %s/install.sh %s", om.Pwd, om.Root, om.Ip,
			om.OMHome, om.Servers[i])
		err = cmd.Run()
		if err != nil {
			msg := "Install " + om.Servers[i] + " module failed"
			return errors.New(msg)
		}

		if (flag_ma == false) && ((om.Servers[i] == "gis") || (om.Servers[i] == "web") || (om.Servers[i] == "token")) {
			flag_om = true
		}
		if (flag_h2 == false) && (om.Servers[i] == "main") {
			flag_h2 = true
			flag_om = true
		}
		if (flag_mq == false) && (om.Servers[i] == "msg") {
			flag_mq = true
		}

		if flag_ma == true {
			cmd = exec.Command("nohup", basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/monitoragent", "start", ">/dev/null", "2>&1", "&")
			err = cmd.Run()
			if err != nil {
				l.Errorf("Start up service monitoragent failed")
			}

			flag_ma = false // only run one time
		}
		if flag_h2 == true {
			cmd = exec.Command("nohup", basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/h2memdb", "start", ">/dev/null", "2>&1", "&")
			err = cmd.Run()
			if err != nil {
				l.Errorf("Start up service h2memdb failed")
			}

			flag_h2 = false
		}
		if flag_mq == true {
			cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/activemq", "start")
			err = cmd.Run()
			if err != nil {
				l.Errorf("Start up service activemq failed")
			}

			flag_mq = false
		}
		if flag_om == true {
			cmd = exec.Command("nohup", basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/onemap", "start", ">/dev/null", "2>&1", "&")
			err = cmd.Run()
			if err != nil {
				l.Errorf("Start up service onemap failed")
			}

			flag_om = false
		}
	}

	return nil
}
