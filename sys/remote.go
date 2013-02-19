package sys

import (
	"errors"
	"fmt"
	"github.com/newthinker/onemap-installer/utl"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// remote copy the OneMap package
func (om *OMPInfo) OMRemoteCopy(srcdir string, dstdir string) error {
	// check srcdir is a file or directory
	if flag := utl.Exists(srcdir); flag != true {
		msg := "Source file or directory " + srcdir + " isn't existed"
		return errors.New(msg)
	}

	fi, _ := os.Stat(srcdir)

	switch curos {
	case "windows":
		// create the remote link
		cmd := exec.Command("cmd", "/C", "net", "use", "\\\\"+om.Ip+"\\admin$",
			om.Pwd, "/user:"+om.Ip+"\\"+om.Root)
		if err := cmd.Run(); err != nil {
			err := errors.New("Create remote link failed")
			l.Error(err)
			return err
		}

		// remote copy
		dstdir = filepath.FromSlash(dstdir)
		temp := strings.Replace(dstdir, ":", "", 1)
		if fi.IsDir() {
			cmd = exec.Command("cmd", "/C", "xcopy", srcdir, "\\\\"+om.Ip+"\\"+temp, "/S")
			if err := cmd.Run(); err != nil {
				err := errors.New("Remote copy directory failed")
				l.Error(err)
				return err
			}
		} else {
			cmd = exec.Command("cmd", "/C", "copy", srcdir, "\\\\"+om.Ip+"\\"+temp, "/S")
			if err := cmd.Run(); err != nil {
				err := errors.New("Remote copy directory failed")
				l.Error(err)
				return err
			}
		}

		// delete the remote link
		cmd = exec.Command("cmd", "/C", "net", "use", "\\\\"+om.Ip+"\\admin$", "/delete")
		if err := cmd.Run(); err != nil {
			err := errors.New("Delete remote link failed")
			l.Error(err)
			return err
		}
	case "linux":
		var cmd *exec.Cmd

		if fi.IsDir() {
			cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "scp", "-r", srcdir, om.Root+"@"+om.Ip+":"+dstdir)

			l.Debugf("sshpass -p %s scp -r %s %s@%s:%s", om.Pwd, srcdir, om.Root, om.Ip, dstdir)
		} else {
			cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "scp", srcdir, om.Root+"@"+om.Ip+":"+dstdir)

			l.Debugf("sshpass -p %s scp %s %s@%s:%s", om.Pwd, srcdir, om.Root, om.Ip, dstdir)
		}
		err := cmd.Run()
		if err != nil {
			err := errors.New("Exec remote copy command failed")
			l.Error(err)
			return err
		}
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

	// service flag
	flag_ma := true  // monitoragent service
	flag_h2 := false // h2memdb service
	flag_mq := false // activemq service
	flag_om := false // onemap service

	var cmd *exec.Cmd

	// exec the remote command line to install the OneMap
	for i := 0; i < len(om.Servers); i++ {
		// exec the server installing command
		switch curos {
		case "windows":
			psexec := filepath.FromSlash(basedir + "/PSTools/PsExec.exe") // path of the psexec
			instbash := filepath.FromSlash(om.OMHome + "/install.bat")    // path of the standalone installation script
			cmd = exec.Command("cmd", "/C", psexec, "\\\\"+om.Ip, "-u", om.Root, "-p", om.Pwd,
				instbash, om.Servers[i])
		case "linux":
			cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/bin/bash", om.OMHome+"/install.sh", om.Servers[i])
			l.Debugf("sshpass -p %s ssh %s@%s /bin/bash %s/install.sh %s", om.Pwd, om.Root, om.Ip,
				om.OMHome, om.Servers[i])
		}
		err := cmd.Run()
		if err != nil {
			msg := "Install " + om.Servers[i] + " server failed"
			return errors.New(msg)
		}

		if (om.Servers[i] == "gis") || (om.Servers[i] == "web") || (om.Servers[i] == "token") {
			flag_om = true
		}
		if om.Servers[i] == "main" {
			flag_h2 = true
			flag_om = true
		}
		if om.Servers[i] == "msg" {
			flag_mq = true
		}
	}

	// install the services
	if flag_ma == true {
		switch curos {
		case "windows":
			psexec := filepath.FromSlash(basedir + "/PSTools/PsExec.exe")               // path of the psexec
			instbash := filepath.FromSlash(om.OMHome + "/bin/service/monitoragent.bat") // path of the monitoragent service installation script
			cmd = exec.Command("cmd", "/C", psexec, "\\\\"+om.Ip,
				"-u", om.Root, "-p", om.Pwd, instbash)
		case "linux":
			cmd = exec.Command("nohup", basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/monitoragent", "start", ">/dev/null", "2>&1", "&")
		}
		err := cmd.Run()
		if err != nil {
			l.Errorf("Start up monitoragent service failed")
		}
	}
	if flag_h2 == true {
		switch curos {
		case "windows":
			psexec := filepath.FromSlash(basedir + "/PSTools/PsExec.exe")          // path of the psexec
			instbash := filepath.FromSlash(om.OMHome + "/bin/service/h2memdb.bat") // path of the h2memdb service installation script
			cmd = exec.Command("cmd", "/C", psexec, "\\\\"+om.Ip,
				"-u", om.Root, "-p", om.Pwd, instbash)
		case "linux":
			cmd = exec.Command("nohup", basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/h2memdb", "start", ">/dev/null", "2>&1", "&")
		}
		err := cmd.Run()
		if err != nil {
			l.Errorf("Start up h2memdb service failed")
		}
	}
	if flag_mq == true {
		switch curos {
		case "windows":
			psexec := filepath.FromSlash(basedir + "/PSTools/PsExec.exe")           // path of the psexec
			instbash := filepath.FromSlash(om.OMHome + "/bin/service/activemq.bat") // path of the activemq service installation script
			cmd = exec.Command("cmd", "/C", psexec, "\\\\"+om.Ip,
				"-u", om.Root, "-p", om.Pwd, instbash)
		case "linux":
			cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/activemq", "start")
		}
		err := cmd.Run()
		if err != nil {
			l.Errorf("Start up activemq service failed")
		}
	}
	if flag_om == true {
		switch curos {
		case "windows":
			psexec := filepath.FromSlash(basedir + "/PSTools/PsExec.exe")         // path of the psexec
			instbash := filepath.FromSlash(om.OMHome + "/bin/service/onemap.bat") // path of the onemap service installation script
			cmd = exec.Command("cmd", "/C", psexec, "\\\\"+om.Ip,
				"-u", om.Root, "-p", om.Pwd, instbash)
		case "linux":
			cmd = exec.Command("nohup", basedir+"/sshpass/bin/sshpass", "-p", om.Pwd, "ssh", om.Root+"@"+om.Ip,
				"/etc/init.d/onemap", "start", ">/dev/null", "2>&1", "&")
		}
		err := cmd.Run()
		if err != nil {
			l.Errorf("Start up onemap service failed")
		}
	}

	return nil
}

// Collect SysConfig from each node
func RemoteCollect(sd *SysDeploy) (los []Layout, err error) {
	for i := 0; i < len(sd.Nodes); i++ {
		var ip string
		var user string
		var pwd string
		var omhome string

		no := sd.Nodes[i]
		for _, a := range no.Attrs {
			if a.Attrname == "ip" {
				ip = a.Attrvalue
			} else if a.Attrname == "user" {
				user = a.Attrvalue
			} else if a.Attrname == "pwd" {
				pwd = a.Attrvalue
			} else if a.Attrname == "omhome" {
				omhome = a.Attrvalue
			}
		}

		// first check the sshpass
		cmd := exec.Command(basedir+"/sshpass/bin/sshpass", "-V")
		err := cmd.Run()
		if err != nil {
			msg := "Sshpass isn't installed"
			l.Errorf(msg)
			return los, errors.New(msg)
		}

		// remote exec the systemconfig.jar process to get the SysConfig.xml
		cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", pwd, "ssh", user+"@"+ip,
			"/bin/bash", omhome+"/config/SystemConfig/SysConfig.sh")
		l.Debugf("sshpass -p %s ssh %s@%s /bin/bash %s/config/SystemConfig/SysConfig.sh",
			pwd, user, ip, omhome)
		err = cmd.Run()
		if err != nil {
			msg := "Remote exec the SystemConfig.jar failed"
			l.Errorf(msg)
			return los, errors.New(msg)
		}

		// remote copy the SysConfig.xml file
		cmd = exec.Command(basedir+"/sshpass/bin/sshpass", "-p", pwd, "scp",
			user+"@"+ip, omhome+"/config/SystemConfig/SysConfig.xml", basedir+"/conf/"+ip+".SysConfig.xml")
		l.Debugf("sshpass -p %s scp %s@%s:%s/config/SystemConfig/SysConfig.xml %s/conf/%s.SysConfig.xml",
			pwd, user, ip, omhome, basedir, ip)
		err = cmd.Run()
		if err != nil {
			msg := "Remote copy the SysConfig.xml file failed"
			l.Errorf(msg)
			return los, errors.New(msg)
		}

		// parse the SysConfig.xml file
		filename := basedir + "/conf/" + ip + ".SysConfig.xml"
		sc, err := OpenSCConfig(filename)
		if err != nil {
			nerr := fmt.Errorf("Open the SysConfig(%s) file failed", filename)
			l.Error(nerr)
			return los, nerr
		}

		los = append(los, sc.LayOut)
	}

	return los, err
}
