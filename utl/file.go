package utl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	//	"path"
	"path/filepath"
	"strings"
)

// 判断文件或者路径是否存在
func Exists(path string) bool {
	path = filepath.FromSlash(path)
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
	srcfile = filepath.FromSlash(srcfile)
	dstfile = filepath.FromSlash(dstfile)

	// first check the srcfile whether exist
	si, serr := os.Stat(srcfile)
	if os.IsNotExist(serr) {
		return os.ErrNotExist
	}

	// check the dstfile whether existed
	di, derr := os.Stat(dstfile)

	// check the srcfile is file or directory
	if si.IsDir() {
		if !os.IsNotExist(derr) { // dst is existed
			if di.IsDir() { // if dst is dir then add the last
				filename := srcfile[strings.LastIndex(srcfile, string(filepath.Separator))+1 : len(srcfile)]
				if filename == "" {
					msg := fmt.Sprintf("Invalid file path(%s)", srcfile)
					return errors.New(msg)
				}
				//				fmt.Println(filename)
				dstfile = filepath.FromSlash(dstfile + "/" + filename)
			} else { // if dst is file then return error
				msg := fmt.Sprintf("Cann't copy a directory(%s) to a file(%s)", srcfile, dstfile)
				return errors.New(msg)
			}
		}
		serr = CopyDir(srcfile, dstfile)
	} else {
		if !os.IsNotExist(derr) { // dst is existed
			if di.IsDir() { // dst is directory and then add the filename
				/// not compatible at windows platform of path package
				filename := srcfile[strings.LastIndex(srcfile, string(filepath.Separator))+1 : len(srcfile)]
				if filename == "" {
					msg := fmt.Sprintf("Invalid file path(%s)", srcfile)
					return errors.New(msg)
				}
				dstfile = filepath.FromSlash(dstfile + "/" + filename)
			} else { // dst is file and delete it first
				if serr = os.Remove(dstfile); serr != nil {
					msg := fmt.Sprintf("Remove the existed file failed(%s)", dstfile)
					return errors.New(msg)
				}
			}
		}
		serr = CopyFile(srcfile, dstfile)
	}
	// exec the copy comand
	if serr != nil {
		fmt.Println(serr)
		return serr
	}

	return nil
}

// Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	//	fmt.Printf("Source file is:%s\n", source)
	//	fmt.Printf("Dest file is:%s\n", dest)
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
		}

	}

	return
}

// Recursively copies a directory tree, attempting to preserve permissions. 
// Source directory must exist, destination directory must *not* exist. 
func CopyDir(source string, dest string) (err error) {
	//	fmt.Printf("Source directory is:%s\n", source)
	//	fmt.Printf("Dest directory is:%s\n", dest)

	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New("Source is not a directory")
	}

	// ensure dest dir does not already exist

	_, err = os.Open(dest)
	if !os.IsNotExist(err) {
		return errors.New("Destination already exists")
	}

	// create dest dir
	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {

		sfp := source + string(filepath.Separator) + entry.Name()
		dfp := dest + string(filepath.Separator) + entry.Name()
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {
			// perform copy         
			err = CopyFile(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		}

	}
	return
}

// list files in the path recursion
func GetAllfiles(path string) {
	path = filepath.FromSlash(path)
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
	path = filepath.FromSlash(path)
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

	return pn, nil
}

// Get local directory
func GetLocalDir() (string, error) {
	basedir, err := filepath.Abs("./")
	if err != nil || basedir == "" {
		return basedir, err
	}

	return basedir, nil
}

// Recursion remove the directory (For windows' long directory)
func RemoveDir(dir string) error {
	if Exists(dir) != true { // the directory isn't existed
		return nil
	}

	entries, err := ioutil.ReadDir(dir)

	for _, entry := range entries {
		sfp := dir + string(filepath.Separator) + entry.Name()
		if entry.IsDir() {
			if err = RemoveDir(sfp); err != nil {
				fmt.Printf("Remove the directory failed:%s\n", sfp)
				return errors.New("Remove directory failed")
			}
		} else {
			// perform remove file         
			err = os.Remove(sfp)
			if err != nil {
				fmt.Printf("Remove the file failed:%s", sfp)
				return errors.New("Remove file failed")
			}
		}
	}

	// remove the empty dir
	err = os.Remove(dir)
	if err != nil {
		fmt.Printf("Remove current directory failed:%s\n", dir)
		return err
	}

	return nil
}
