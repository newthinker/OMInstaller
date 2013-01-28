package sys

import (
	"errors"
)

///////////////////////////////
/// post.json 
type SrvBase struct {
	Os        string
	Arch      string
	Ip        string
	User      string
	Pwd       string
	Omhome    string
	Container string
}

type ParamPost struct {
	Paramname  string
	Paramvalue string
	Encrypt    string
	Selects    string
}

type SrvPost struct {
	Srvname string
	Params  []ParamPost
}

type ServerParams struct {
	Server_base   SrvBase
	Server_params []SrvPost
}

// 解析POST的JSON结构
func ParseSysSubmit(jsonstr interface{}, basepath string, sc *SysConfig, sm *ServerMapping) error {
	postmap := jsonstr.(map[string]interface{})

	for _, v := range postmap {
		switch vv := v.(type) {
		case string:
		case int:
			if vv != 0 {
				l.Errorf("Remote return code error, please check")
				return errors.New("Remote return code error, please check")
			}
		case []interface{}:
			// 获取输入参数信息
			sd := &SysDeploy{}
			//var arr_sc = []SysConfig{}

			// 开始解析数据体部分
			/*		for i, s := range vv {
						l.Messagef("Parse the %dth machine's params", i+1)

						//srvparams := s.(map[string]interface{})
						//base := (srvparams["Server_base"]).(map[string]interface{})

						//	switch arrattr := (srvparams["Server_base"]).(type) {
						//	case []interface{}:
						//for _, aa := range arrattr {
						//	sa := aa.(map[string]interface{})

						//	}

						switch arrsrv := (srvparams["Server_params"]).(type) {
						case []interface{}:
							for _, ss := range arrsrv {
								sp := ss.(map[string]interface{})

								server := &ServerInfo{}

								server.Srvname = (sp["Srvname"]).(string)

								switch arrparams := (sp["Params"]).(type) {
								case []interface{}:
									for _, p := range arrparams {
										param := p.(map[string]interface{})

										attr := &AttrInfo{}
										attr.Attrname = (param["Paramname"]).(string)
										attr.Value = (param["Paramvalue"]).(string)

										encrypt := (param["Encrypt"]).(string)
										if encrypt != "" {
											attr.Encrypt = encrypt
										}
										selects := (param["Selects"]).(string)
										if selects != "" {
											//attr.Select = selects
											attr.Value = selects
										}

										server.Attrs = append(server.Attrs, *attr)
									}
								}

								//		machine.Servers = append(machine.Servers, *server)

							}
						}

						//	si.Machines = append(si.Machines, *machine)

					}
			*/

			// update the SysDeploy.xml config file
			conffile := basedir + "/conf/" + SYS_DEPLOY
			if err := RefreshSysDeploy(sd, conffile); err != nil {
				l.Errorf("Update the SysDeploy.xml config file failed")
				return err
			}

			//		err := Distribute(sd, arr_sc)
			//		if err != nil {
			//			l.Errorf("Distribute installing failed")
			//			return err
			//		}
		}
	}

	return nil
}

func SysFormat(status int) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if status == INSTALL {
		result["Servers"] = omsc.LayOut.Servers
	} else if status == UPDATE {

	} else if status == UNINSTALL {

	}

	l.Debug(result)

	return result, nil
}
