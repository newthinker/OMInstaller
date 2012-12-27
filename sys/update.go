package sys

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
    "errors"
)

// Read a whole file into the memory and store it as array of lines
func readLines(path string) (lines []string, err error) {
	var (
		file   *os.File
		part   []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func writeLines(om *OMPInfo, lines []string, path string) (err error) {
	var (
		file *os.File
	)

	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	for _, item := range lines {
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			fmt.Println(err)
			break
		}

        // 查找插入点
        flag := "#!/bin/bash"
        if ins:=strings.Contains(item, flag);ins==true {
            addons,err:=formatAddon(om)
            if err!=nil {
                return err
            }

            for _,addon := range addons {
               _, err = file.WriteString(strings.TrimSpace(addon) + "\n")
               if err != nil {
                fmt.Println(err)
                break
               }
            }
        }
	}

	return
}

// prepare the addon parameters
func formatAddon(om *OMPInfo) (addon []string, err error) {
    if om==nil {
        return addon, errors.New("ERROR: 输入参数为空")
    }

    // public part
    if om.OMHome=="" {
        return addon, errors.New("ERROR: OneMap安装目录不存在，请检查")
    }
    addon = append(addon, "########################################")
    addon = append(addon, "###Input params###")
    addon = append(addon, "ONEMAP_HOME=\""+om.OMHome+"\"")

    if om.Container=="" {
        return addon, errors.New("ERROR: 获取OneMap WEB容器失败")
    }
    addon = append(addon, "CONTAINER_NAME="+om.Container)

    addon = append(addon, "ESRI_GROUP=\""+om.OM_Group+"\"")
    addon = append(addon, "OM_ACCOUNT=\""+om.OM_User+"\"")
    addon = append(addon, "OM_PWD=\""+om.OM_PWD+"\"")
    addon = append(addon, "ULIMIT_NUM=10240")

    // db server
    flag_db := false
    flag_gis := false
    for _, srvtype := range om.Servers {
        if srvtype=="db" {
            flag_db = true
        }

        if srvtype=="gis" {
            flag_gis = true
        }
    }
    if flag_db==true {
        addon = append(addon, "ORCL_ACCOUNT=\""+om.ORCL_User+"\"")
        addon = append(addon, "ORACLE_SYSTEM_ACCOUNT=\""+om.DB_User[0]+"\"")
        addon = append(addon, "ORACLE_SYSTEM_PWD=\""+om.DB_PWD[0]+"\"")
        addon = append(addon, "ORACLE_SID=\""+om.ORCL_SID+"\"")
        addon = append(addon, "MANAGER_USER=\""+om.DB_User[1]+"\"")
        addon = append(addon, "MANAGER_PWD=\""+om.DB_PWD[1]+"\"")
        addon = append(addon, "PORTAL_USER=\""+om.DB_User[2]+"\"")
        addon = append(addon, "PORTAL_PWD=\""+om.DB_PWD[2]+"\"")
        addon = append(addon, "GEOCODING_USER=\""+om.DB_User[3]+"\"")
        addon = append(addon, "GEOCODING_PWD=\""+om.DB_PWD[3]+"\"")
        addon = append(addon, "GEOPORTAL_USER=\""+om.DB_User[4]+"\"")
        addon = append(addon, "GEOPORTAL_PWD=\""+om.DB_PWD[4]+"\"")
        addon = append(addon, "SUB_USER=\""+om.DB_User[5]+"\"")
        addon = append(addon, "SUB_PWD=\""+om.DB_PWD[5]+"\"")
    }

    // gis server
    if flag_gis==true {
        addon = append(addon, "AGS_HOME=\""+om.AGS_Home+"\"")
    }
    addon = append(addon, "########################################")

    return
}

// 更新单机安装脚本，插入需要参数
func UpdateScritp(om *OMPInfo, scriptfile string) error {
    // 检查脚本是否存在
    if (Exists(scriptfile)!=true)  {
        return errors.New("ERROR: 安装脚本不存在")
    }
    // 读取脚本
	lines, err := readLines(scriptfile)
	if err != nil {
	    fmt.Println(err)
		return err
	}
	for _, line := range lines {
		fmt.Println(line)
	}
    // 生成副本并将参数添加到副本中
    bakfile := scriptfile + ".bak"
	err = writeLines(om, lines, bakfile)
    if err!=nil {
	    fmt.Println(err)
        return err
    }
    // 删除原始文件
    if err=os.Remove(scriptfile); err!=nil {
        fmt.Println(err)
        return err
    }
    // 重命名副本
    if err=os.Rename(bakfile, scriptfile); err!=nil {
        fmt.Println(err)
        return err
    }

    return nil
}
