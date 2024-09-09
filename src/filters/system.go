package filters

import (
	"sports-common/consts"

	"github.com/flosch/pongo2"
)

// 底部信息配置 - 内容类型
func menuLevelType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.MenuLevels, "-")
}
