package filters

import (
	"sports-common/consts"

	"github.com/flosch/pongo2"
)

// 日志模块
func userLogModule(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(int); isType {
		if val, exists := consts.UserLogModules[value]; exists {
			return pongo2.AsValue(val), nil
		}
		return pongo2.AsValue("-"), nil
	}
	return pongo2.AsValue(""), nil
}

// 日志类型
func userLogType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(int); isType {
		if val, exists := consts.UserLogTypes[value]; exists {
			return pongo2.AsValue(val), nil
		}
		return pongo2.AsValue("-"), nil
	}
	return pongo2.AsValue(""), nil
}
