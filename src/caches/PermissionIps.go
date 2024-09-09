package caches

import (
	"fmt"
	"os"
	models "sports-models"
)

const keyPermissionIPs = "permission_ips"

// PermissionIps 后台授权IP
var PermissionIps = struct {
	Load          func(string)
	All           func(string) map[uint32]models.PermissionIp
	HasPermission func(string, string) bool //此IP是否有访问权限
}{
	Load: func(platform string) { // 将授权的ip加载到内存当中
		permissionIPs := map[uint32]models.PermissionIp{}
		var rs []models.PermissionIp
		err := models.PermissionIps.FindAllNoCount(platform, &rs, nil, "id ASC")
		if err != nil {
			fmt.Println("加载可信IP出错: ", err)
			os.Exit(1)
			return
		}

		for _, r := range rs {
			permissionIPs[r.Id] = r
		}

		_ = setCache(platform, keyPermissionIPs, permissionIPs)
	},
	All: func(platform string) map[uint32]models.PermissionIp {
		permissionIPs := map[uint32]models.PermissionIp{}
		_ = getCache(platform, keyPermissionIPs, &permissionIPs)
		return permissionIPs
	},
	HasPermission: func(platform, ip string) bool {
		permissionIPs := map[uint32]models.PermissionIp{}
		_ = getCache(platform, keyPermissionIPs, &permissionIPs)

		for _, v := range permissionIPs {
			if v.Ip == ip && v.Enabled() {
				return true
			}
		}
		return false
	},
}
