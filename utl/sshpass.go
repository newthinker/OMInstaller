// sshpass.go
package utl

import (
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

// install the sshpass package
func InstallSshpass(base string) error {
	sshpass := filepath.FromSlash(base + "/sshpass/bin/sshpass")
	if flag := Exists(sshpass); flag != true {
		sshpass := filepath.FromSlash(base + "/sshpass/Install.sh")
		if flag = Exists(sshpass); flag != true {
			fmt.Println("No sshpass software package")
			return errors.New("No sshpass software package")
		}

		// exec the install script
		cmd := exec.Command("/bin/sh", filepath.FromSlash(base+"/sshpass/Install.sh"), filepath.FromSlash(base+"/sshpass"))
		err := cmd.Run()
		if err != nil {
			fmt.Println("Complier sshpass failed")
			return errors.New("Complier sshpass failed")
		}

		// whether install successfully
		if flag := Exists(filepath.FromSlash(base + "/sshpass/bin/sshpass")); flag != true {
			fmt.Println("Install sshpass failed")
			return errors.New("Install sshpass failed")
		}
	} else {
		fmt.Println("Sshpass is installed and go on")
	}

	return nil
}

// check the sshpass whether installed successfully
func CheckSshpass(base string) error {
	cmd := exec.Command(filepath.FromSlash(base+"/sshpass/bin/sshpass"), "-V")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Sshpass isn't installed")
		return err
	}

	return nil
}
