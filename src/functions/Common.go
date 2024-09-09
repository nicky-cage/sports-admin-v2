package functions

import (
	"os"
	"sports-admin/caches"
	models "sports-models"
	"strings"

	"github.com/flosch/pongo2"
)

// GetEnv
func getEnv(args ...*pongo2.Value) *pongo2.Value {
	return pongo2.AsValue(os.Getenv(args[0].String()))
}

// 判断是否有权限访问此url
func checkURLPermission(menus []models.Menu, url string, level int, hasPermission *bool) {
	if *hasPermission {
		return
	}
	for _, menu := range menus {
		if menu.Url == url && int(menu.Level) == level {
			*hasPermission = true
			return
		}
		if strings.Index(menu.Url, "|") > 0 {
			tArr := strings.Split(menu.Url, "|")
			if (tArr[0] == url || tArr[1] == url) && int(menu.Level) == level {
				*hasPermission = true
				return
			}
		}
		if len(menu.Children) > 0 {
			checkURLPermission(menu.Children, url, level, hasPermission)
		}
	}
}

// 此url是否已经授权
// platformV: 平台
// roleIDV: 角色id
// urlV: 路由url
// levelV: 路由level
func isGranted(args ...*pongo2.Value) *pongo2.Value {
	//t1 := tools.TimeDebugBegin("is-granted")
	//defer tools.TimeDebugAt(t1, "is-granted-end")
	platform := args[0].String()
	roleID := args[1].Integer()

	role := caches.AdminRoles.Get(platform, roleID)
	if role == nil {
		return pongo2.AsValue(false)
	}

	//return pongo2.AsValue(true)

	url := args[2].String()
	level := args[3].Integer()
	hasPermission := false
	checkURLPermission(role.Menus, url, level, &hasPermission)

	return pongo2.AsValue(hasPermission)
}
