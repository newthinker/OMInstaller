package sys

import (
	"errors"
	"fmt"
)

///////////////////////////////
// 解析POST的JSON结构
func ParseSysSubmit(jsonstr interface{}) (SysDeploy, []Layout, error) {
	postmap := jsonstr.(map[string]interface{})
	fmt.Println(jsonstr)

	// 获取输入参数信息
	sd := &SysDeploy{}
	var arr_lo = []Layout{}
    var status int

	for _, v := range postmap {
		switch vv := v.(type) {
            case int:   // status flag [0:maintain, 1: install, 2: update, 3: uninstall]
            status = vv
			if status < 0 || status > 3 {
				err := errors.New("Invalid status code, please check")
				l.Error(err)
				return *sd, arr_lo, err
			}
		case []interface{}:
			// 开始解析数据体部分
			for i, s := range vv {
				l.Messagef("Parse the %dth node's params", i+1)

				srvparams := s.(map[string]interface{})

				curnode := &node{}

				// base info
				switch bases := (srvparams["Base"]).(type) {
				case []interface{}:     /// Test the bases is null(install mode with no base section)
					for _, ba := range bases {
						base := ba.(map[string]interface{})

						attr := &attr{}
						attr.Attrname = (base["Attrname"]).(string)
						attr.Attrvalue = (base["Attrvalue"]).(string)

						curnode.Attrs = append(curnode.Attrs, *attr)
					}
				}

				// params info
				lo := &Layout{}
				switch arrsrv := (srvparams["Params"]).(type) {
				case []interface{}:
					for _, ss := range arrsrv {
						sp := ss.(map[string]interface{})

						server := &srv{}
						server.Srvname = (sp["Srvname"]).(string)
						curnode.Srvs = append(curnode.Srvs, *server)

						srvinfo := &ServerInfo{}
						srvinfo.Srvname = (sp["Srvname"]).(string)

						switch arrparams := (sp["Attrs"]).(type) {
						case []interface{}:
							for _, p := range arrparams {
								param := p.(map[string]interface{})
								attr := &AttrInfo{}
                                value, ok := (param["Attrname"]).(string)
                                if ok != true || value == "" {
                                    continue
                                }
                                attr.Attrname = value

                                value, ok = (param["Attrvalue"]).(string)
                                if ok != true || value == "" {
                                    continue
                                }
                                attr.Attrvalue = value

								encrypt, ok := (param["Encrypt"]).(string)
								if ok == true && encrypt != "" {
									attr.Encrypt = encrypt
								}

								selects, ok := (param["Selects"]).(string)
								if ok == true && selects != "" {
									attr.Attrvalue = selects
								}

								srvinfo.Attrs = append(srvinfo.Attrs, *attr)
							}
						}

						lo.Servers = append(lo.Servers, *srvinfo)
					}
				}

				arr_lo = append(arr_lo, *lo)

				sd.Nodes = append(sd.Nodes, *curnode)
			}
		}
	}

	return *sd, arr_lo, nil
}

/// Format the exchange data
func SysFormat(status int) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	if status == INSTALL {
		result["Servers"] = omsc.LayOut.Servers
	} else if status == UPDATE {
        // first get the SysDeploy and Layout array struct
        los, err := RemoteCollect(omsd)
        if err!=nil {
            l.Error(err)
            return result, err
        }

        fmt.Println(los)

        /// then format exchange data and post
        /// err = UpdateFormat(sd *SysDeploy, los []Layout)
	} else if status == UNINSTALL {

	}

	l.Debug(result)

	return result, nil
}
