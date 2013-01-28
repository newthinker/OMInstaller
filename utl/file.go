package utl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// 判断文件或者路径是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return false
}

// Copy file or dirctory
func Copy(srcfile string, dstfile string) error {
	// first check the srcfile whether exist
	fi, serr := os.Stat(srcfile)
	if os.IsNotExist(serr) {
		return os.ErrNotExist
	}

	// check dstfile's parent path whether existed
	dir := filepath.Dir(dstfile)
	_, derr := os.Stat(dir)
	if os.IsNotExist(derr) {

		if serr = os.MkdirAll(dir, 0755); serr != nil {
			return serr
		}
	}

	// check the srcfile is file or directory
	if fi.IsDir() {
		cmd := exec.Command("cp", "-r", srcfile, dstfile)
		serr = cmd.Run()
	} else {
		cmd := exec.Command("cp", srcfile, dstfile)
		serr = cmd.Run()
	}
	// exec the copy comand
	if serr != nil {
		return serr
	}

	return nil
}

// list files in the path recursion
func GetAllfiles(path string) {
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		println(path)
		return nil
	})

	if err != nil {
		fmt.Printf("filepath.Walk() return %v", err)
	}
}

// list sub directory in current path
func GetSubDir(path string) ([]string, error) {
	pn := []string{}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Open input path(%s) failed", path)
		return pn, err
	}

	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		fmt.Printf("Read input path(%s) failed", path)
		return pn, err
	}

	for _, fileinfo := range list {
		if fileinfo == nil {
			continue
		}
		if fileinfo.IsDir() {
			var pathname string = fileinfo.Name()
			if pathname != "" {
				pn = append(pn, pathname)
			}
		}
	}

	return pn, err
}

//Get local directory
func GetLocalDir() (string, error) {
	basedir, err := filepath.Abs("./")
	if err != nil || basedir == "" {
		return basedir, err
	}

	return basedir, nil
}
