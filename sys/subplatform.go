package sys

import (
	"bufio"
	"errors"
	"github.com/newthinker/onemap-installer/utl"
	"io"
	"os"
	"strings"
)

// 从sql文件中查找目标行的标识
var ARR_FLAG = [...]string{"MAINTENACE_FRAMEWORK_MODULES", "insert", "values"}

type SubPlatform struct {
	SqlFile string            // sql文件保存路径 
	MenuMap map[string]string // id-菜单名map
	RelMap  map[string]string // 子节点与父节点对应map
	SelID   []string          // 选择的菜单id
}

// 从sql文件中解析出所有需要配置的列表项
func (sp *SubPlatform) SPParseSQLFile() error {
	sqlfile := sp.SqlFile

	l.Messagef("Begin to parse the sql file:%s", sqlfile)
	// 首先判断文件是否存在
	if (utl.Exists(sqlfile)) != true {
		msg := "File isn't existed!"
		l.Errorf(msg)
		return errors.New(msg)
	}

	// 打开并解析sql文件
	file, err := os.Open(sqlfile)
	if err != nil {
		msg := "Open the sql file failed!"
		l.Errorf(msg)
		return errors.New(msg)
	}
	defer file.Close()

	var sqlstate string
	reader := bufio.NewReader(file)
	for {
		var str string
		str, err = reader.ReadString('\n')

		if err == io.EOF {
			sqlstate = ""
			break
		}

		if err != nil {
			msg := "Read the sql file failed!"
			l.Errorf(msg)
			return errors.New(msg)
		}

		//str = str[0 : len(str)-2] // 去掉换行符 
		sqlstate = sqlstate + " " + str
		// 如果获取到了一个完整的sql语句 
		if strings.Index(str, ";") > 0 {
			id, parentid, name := sp.parseSqlState(sqlstate)

			if id == "" || parentid == "" || name == "" {
				sqlstate = ""
				continue
			} else {
				sp.MenuMap[id] = name
				sp.RelMap[id] = parentid

				sqlstate = ""
			}
		}
	}

	l.Message("End parsing the sql file")

	return nil
}

// 解析单条sql语句
func (sp *SubPlatform) parseSqlState(sqlstate string) (string, string, string) {
	var id string
	var parentid string
	var name string

	if strings.Index(sqlstate, ";") < 0 {
		l.Warningf("Sql statement(%s) is incomplete", sqlstate)
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
		l.Warning("Incompleted statement")
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

	return id, parentid, name
}

// 更新sql文件
func (sp *SubPlatform) SPUpdateSql() error {
	sqlfile := sp.SqlFile

	if (utl.Exists(sqlfile)) != true {
		msg := "File isn't existed!"
		l.Errorf(msg)
		return errors.New(msg)
	}

	// open and parse the sql file
	infile, err := os.Open(sqlfile)
	if err != nil {
		msg := "Open the sql file failed!"
		l.Errorf(msg)
		return errors.New(msg)
	}

	// check and create the bak file at the same time
	bakfile := sqlfile + ".bak"
	if (utl.Exists(bakfile)) == true {
		if err = os.Remove(bakfile); err != nil {
			err = errors.New("Remove bak sql file failed")
			l.Error(err)
			return err
		}
	}
	outfile, nerr := os.Create(bakfile)
	if nerr != nil {
		msg := "Create new file failed!"
		l.Errorf(msg)
		return errors.New(msg)
	}

	// update the sql file
	var sqlstate string
	reader := bufio.NewReader(infile)
	for {
		// read line by line
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
			msg := "Read the sql file failed!"
			infile.Close()
			outfile.Close()
			l.Errorf(msg)
			return errors.New(msg)
		}

		sqlstate = sqlstate + " " + str
		// when get a whole line sql statement 
		if strings.Index(str, ";") > 0 {
			var id string
			id, _, _ = sp.parseSqlState(sqlstate)

			// check this id whether in selected ids
			if id != "" && id != "0" {
				flag := false
				for selid := range sp.SelID {
					if id == sp.SelID[selid] {
						flag = true
						break
					}
				}

				// if not in then not update into the bak file
				if flag == false {
					sqlstate = ""
					continue
				}
			}

			// output the bak file
			outfile.Write([]byte(sqlstate))
			sqlstate = ""
		}
	}

	infile.Close()
	outfile.Close()

	return nil
}
