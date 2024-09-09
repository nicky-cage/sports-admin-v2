package filters

import (
	"sports-admin/caches"

	"github.com/flosch/pongo2"
)

// 平台名称
func getPlatformName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(uint32); isType {
		if v := caches.Platforms.Get(int(value)); v != nil {
			return pongo2.AsValue(v.Name), nil
		}
	}
	return pongo2.AsValue(""), nil
}

// 站点名称
func getPlatformSiteName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(uint32); isType {
		if v := caches.PlatformSites.Get(int(value)); v != nil {
			return pongo2.AsValue(v.Name), nil
		}
	}
	return pongo2.AsValue(""), nil
}
