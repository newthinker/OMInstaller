package sys

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// 从sql文件中查找目标行的标识
var ARR_FLAG = [...]string{"MAINTENACE_FRAMEWORK_MODULES", "insert", "values"}

// 定义menu对应关系的嵌套结构
//type    relMap    map[int]interface{}

type SubPlatform struct {
	SqlFile string            // sql文件保存路径 
	MenuMap map[string]string // id-菜单名map
	RelMap  map[string]string // 子节点与父节点对应map
	SelID   []string          // 选择的菜单id
}

// 从sql文件中解析出所有需要配置的列表项
func (sp *SubPlatform) SPParseSQLFile() error {
	sqlfile := sp.SqlFile

	fmt.Printf("Begin to parse the sql file:%s\n", sqlfile)
	// 首先判断文件是否存在
	if (Exists(sqlfile)) != true {
		msg := "ERROR: File isn't existed!"
		return errors.New(msg)
	}

	// 打开并解析sql文件
	file, err := os.Open(sqlfile)
	if err != nil {
		msg := "ERROR: Open the sql file failed!"
		return errors.New(msg)
	}
	defer file.Close()

	var sqlstate string
	reader := bufio.NewReader(file)
	for {
		var str string
		str, err = reader.ReadString('\n')

		//        fmt.Println(str)

		if err == io.EOF {
			sqlstate = ""
			break
		}

		if err != nil {
			msg := "ERROR: Read the sql file failed!"
			return errors.New(msg)
		}

		//str = str[0 : len(str)-2] // 去掉换行符 
		sqlstate = sqlstate + " " + str
		// 如果获取到了一个完整的sql语句 
		if strings.Index(str, ";") > 0 {
			id, parentid, name := sp.parseSqlState(sqlstate)

			if id == "" || parentid == "" || name == "" {
				//				fmt.Printf("WARN: Parse sql state(%s) failed!\n", str)

				sqlstate = ""
				continue
			} else {
				sp.MenuMap[id] = name
				sp.RelMap[id] = parentid

				sqlstate = ""
			}
		}
	}

	fmt.Println("End parsing the sql file!")

	return nil
}

// 解析单条sql语句
func (sp *SubPlatform) parseSqlState(sqlstate string) (string, string, string) {
	var id string
	var parentid string
	var name string

	if strings.Index(sqlstate, ";") < 0 {
		fmt.Println("WARN: Sql statement(" + sqlstate + ") is incomplete!")
		return id, parentid, name
	}

	// 开始解析
	for i := 0; i < len(ARR_FLAG); i++ {
		sample := ARR_FLAG[i]
		if sample == "" || strings.Index(sqlstate, sample) < 0 {
			return id, parentid, name
		}
	}

	arrValue := strings.Split(sqlstate, "values")
	if len(arrValue) < 2 {
		fmt.Println("WARN: Incompleted statement")
		return id, parentid, name
	}

	values := arrValue[1]
	values = strings.Trim(values, " ")
	values = strings.TrimLeft(values, "(")
	arrValue = strings.Split(values, ",")

	// 获取id, parentid及name
	id = arrValue[0]
	parentid = arrValue[1]
	name = arrValue[2]

	//fmt.Println("id:", id, ";Parentid:", parentid, ";name:", name)

	return id, parentid, name
}

// 更新sql文件
func (sp *SubPlatform) SPUpdateSql() error {
	sqlfile := sp.SqlFile

	if (Exists(sqlfile)) != true {
		msg := "ERROR: File isn't existed!"
		return errors.New(msg)
	}

	// 打开sql文件进行解析
	infile, err := os.Open(sqlfile)
	if err != nil {
		msg := "ERROR: Open the sql file failed!"
		return errors.New(msg)
	}

	// 同时生成备份文件
	bakfile := sqlfile + ".bak"
	outfile, nerr := os.Create(bakfile)
	if nerr != nil {
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

		if err == io.EOF {
			if sqlstate != "" {
				outfile.Write([]byte(sqlstate))
				sqlstate = ""
			}
			break
		}

		if err != nil {
			msg := "ERROR: Read the sql file failed!"
			infile.Close()
			outfile.Close()
			return errors.New(msg)
		}

		//		str = str[0 : len(str)-1] // 去掉换行符 
		sqlstate = sqlstate + " " + str
		// 如果获取到了一个完整的sql语句 
		if strings.Index(str, ";") > 0 {
			var id string
			id, _, _ = sp.parseSqlState(sqlstate)

			// 判断此id是否在所选序列中
			if id != "" && id != "0" {
				flag := false // 是否在所选序列标识
				for selid := range sp.SelID {
					if id == sp.SelID[selid] {
						flag = true
						break
					}
				}

				// 如果不在所选列表中，不写入到新文件
				if flag == false {
					sqlstate = ""
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
	if err = os.Remove(sqlfile); err != nil {
		msg := "ERROR: Delete the file(" + sqlfile + ") failed!"
		outfile.Close()
		return errors.New(msg)
	}

	//更新sql文件
	outfile.Close()
	if err = os.Rename(bakfile, sqlfile); err != nil {
		msg := "ERROR: Rename the file(" + bakfile + ") failed!"
		outfile.Close()
		return errors.New(msg)
	}

	return nil
}
