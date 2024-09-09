package functions

import (
	"github.com/flosch/pongo2"
)

// 模板函数
var customFunctions = map[string]func(...*pongo2.Value) *pongo2.Value{
	"get_env":    getEnv,
	"is_granted": isGranted,
}

// All 注册所有模板函数
func All() map[string]func(...*pongo2.Value) *pongo2.Value {
	return customFunctions
}
