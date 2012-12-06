package sys

import (
    "fmt"
    "os"
    "io"
    "strings"
    "bufio"
    "errors"
)

// 从sql文件中查找目标行的标识
var ARR_FLAG = [...]string {"MAINTENACE_FRAMEWORK_MODULES", "insert", "values"}

// 定义menu对应关系的嵌套结构
//type    relMap    map[int]interface{}

type SubPlatform struct {
    sqlFile         string              // sql文件保存路径 
    menuMap         map[string]string   // id-菜单名map
    relMap          map[string]string   // 子节点与父节点对应map
    selID           []string            // 选择的菜单id
}

// 从sql文件中解析出所有需要配置的列表项
func (sp *SubPlatform) SPGetMenuMap(sqlfile string) error {
    fmt.Printf("Begin to parse the sql file:%s\n", sqlfile)
    // 首先判断文件是否存在
    if (Exists(sqlfile))!=true {
		msg := "ERROR: File isn't existed!"
		return errors.New(msg)
    }

    // 打开并解析sql文件
    file, err := os.Open(sqlfile)
    if err!=nil {
        msg := "ERROR: Open the sql file failed!"
        return errors.New(msg)
    }
    defer file.Close()

    var sqlstate string
    reader := bufio.NewReader(file)
    for {
        var str string
        str, err = reader.ReadString('\n')

        if err==io.EOF {
            sqlstate = ""
            break
        }

        if err!=nil {
            msg := "ERROR: Read the sql file failed!"
            return errors.New(msg)
        }

        str = str[0:len(str)-1]      // 去掉换行符 
        sqlstate += sqlstate + " " + str
        // 如果获取到了一个完整的sql语句 
        if strings.Index(str, ";")>0 {
            if id,parentid,name,nerr:=sp.parseSqlState(str); id=="" || parentid=="" || name=="" || nerr!=nil {
                fmt.Printf("ERROR: Parse sql state(%s) failed!\n", str)
            } else {
                sp.menuMap[id] = name
                sp.relMap[id] = parentid
            }
        }
    }

    fmt.Println("End parsing the sql file!")

    return nil
}

// 解析单条sql语句
func (sp *SubPlatform) parseSqlState(sqlstate string) (string, string, string, error) {
    var id string
    var parentid string
    var name string

    if strings.Index(sqlstate, ";")<0 {
        msg := "ERROR: Sql statement(" + sqlstate + ") is incomplete!"
        return id,parentid,name,errors.New(msg)
    }

    // 开始解析
    for i:=0;i<len(ARR_FLAG);i++ {
        sample := ARR_FLAG[i]
        if sample =="" || strings.Index(sqlstate, sample)<0 {
            return id, parentid, name, nil
        }
    }

    arrValue := strings.Split(sqlstate, "values")
    if len(arrValue)<2 {
        msg := "ERROR: Incompleted statement"
        return id, parentid, name, errors.New(msg)
    }

    values := arrValue[1]

    arrValue = strings.Split(values, ",")

    // 获取id, parentid及name
    id = arrValue[0]
    parentid = arrValue[1]
    name = arrValue[2]

    return id, parentid, name, nil
}

// 更新sql文件
func (sp *SubPlatform) SPUpdateSql(sqlfile string) error {
    if (Exists(sqlfile))!=true {
		msg := "ERROR: File isn't existed!"
		return errors.New(msg)
    }

    // 打开sql文件进行解析
    infile, err := os.Open(sqlfile)
    if err!=nil {
        msg := "ERROR: Open the sql file failed!"
        return errors.New(msg)
    }

    // 同时生成备份文件
    bakfile := sqlfile + ".bak"
    outfile, nerr := os.Create(bakfile)
    if nerr!=nil {
        msg := "ERROR: Create new file failed!"
        return errors.New(msg)
    }

    // 进行更新sql文件
    var sqlstate string
    reader := bufio.NewReader(infile)
    for {
        // 逐行读取
        var str string
        str, err = reader.ReadString('\n')

        if err==io.EOF {
            sqlstate = ""
            break
        }

        if err!=nil {
            msg := "ERROR: Read the sql file failed!"
            outfile.Close()
            return errors.New(msg)
        }
        
        str = str[0:len(str)-1]      // 去掉换行符 
        sqlstate += sqlstate + " " + str
        // 如果获取到了一个完整的sql语句 
        if strings.Index(str, ";")>0 {
            var id string
            if id,_,_,err=sp.parseSqlState(str); err!=nil {
                fmt.Printf("ERROR: Parse sql state(%s) failed!\n", str)
                continue
            }

            // 判断此id是否在所选序列中
            if id!="" {
                flag := false       // 是否在所选序列标识
                for selid := range sp.selID {
                    if id==sp.selID[selid] {
                        flag = true
                        break
                    }
                }

                if flag==false {
                    continue
                }
            }

            // 删除不在所选序列中的menu
            outfile.Write([]byte(sqlstate))
            sqlstate = ""
        }
    }

    // 删除原文件
    infile.Close()
    if err=os.Remove(sqlfile);err!=nil {
        msg := "ERROR: Delete the file(" + sqlfile + ") failed!"
        outfile.Close()
        return errors.New(msg)
    }
    
    //更新sql文件
    outfile.Close()
    if err=os.Rename(bakfile, sqlfile);err!=nil {
        msg := "ERROR: Rename the file(" + bakfile + ") failed!"
        outfile.Close()
        return errors.New(msg)
    }

    return nil
}

