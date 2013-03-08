package sys

import (
	"bufio"
	"errors"
	"github.com/newthinker/onemap-installer/utl"
	"io"
	"os"
	"strings"
)

// search flags in the sql sentences
var ARR_FLAG = [...]string{"MAINTENACE_FRAMEWORK_MODULES", "insert", "values"}

type SubPlatform struct {
	SqlFile string            // the sql filename 
	MenuMap map[string]string // menu map
	RelMap  map[string]string // parent-sun nodes map
	SelID   []string          // selected nodes' ids
}

// Parse the sql file with config menus
func (sp *SubPlatform) SPParseSQLFile() error {
	sqlfile := sp.SqlFile

	l.Messagef("Begin to parse the sql file:%s", sqlfile)
	// first check whether the sql file is there
	if (utl.Exists(sqlfile)) != true {
		msg := "File isn't existed!"
		l.Errorf(msg)
		return errors.New(msg)
	}

	// open and parse the sql file
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

		//str = str[0 : len(str)-2] // remove line breaks 
		sqlstate = sqlstate + " " + str
		// get a complete sql sentences 
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

// Parse single sql statement
func (sp *SubPlatform) parseSqlState(sqlstate string) (string, string, string) {
	var id string
	var parentid string
	var name string

	if strings.Index(sqlstate, ";") < 0 {
		l.Warningf("Sql statement(%s) is incomplete", sqlstate)
		return id, parentid, name
	}

	// parse
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

	// get id, parentid and name
	id = arrValue[0]
	parentid = arrValue[1]
	name = arrValue[2]

	return id, parentid, name
}

// Update the sql file
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
