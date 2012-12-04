package sys

import (
    "fmt"
)

// 从sql文件中查找目标行的标识
const ARR_FLAG[3] = {"MAINTENACE_FRAMEWORK_MODULES", "insert", "values"}

type SubPlatform struct {
    sqlFile     string          // sql文件保存路径 
    menuMap     map[int]string  // id-菜单名map
    delID       []int           // 需要删除的菜单id
}

// 从sql文件中解析出所有需要配置的列表项

// 根据用户选择的菜单组解析出需要删除的所有菜单id

// 更新sql文件

