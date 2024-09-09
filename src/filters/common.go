package filters

import (
	"fmt"

	"github.com/flosch/pongo2"
)

// 状态 - 通用
// 默认: 2/1 => 正常/禁用
func stateText(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value := in.Integer()
	text := "<span style='color:green'>正常</span>"
	if value != 2 {
		text = "<span style='color:red'>禁用</span>"
	}
	return pongo2.AsValue(text), nil
}

// PlatformWrap - 将平台信息包装到一起,以供后续调用
func platformWrap(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	platform := param.Interface().(string) // 平台信息
	value := fmt.Sprintf("%s:%v", platform, in)
	return pongo2.AsValue(value), nil
}
