package admin_roles

import (
	"fmt"
	"sports-admin/caches"
	"sports-common/request"
	"sports-common/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ths *AdminRoles) SubMenus(c *gin.Context) {
	menuID, err := strconv.Atoi(c.DefaultQuery("id", ""))
	if err != nil {
		response.Err(c, fmt.Sprintf("菜单编号有误: %d", menuID))
		return
	}

	platform := request.GetPlatform(c)
	allMenus := caches.Menus.LayMenus(platform)
	for _, m1 := range allMenus {
		for _, m2 := range m1.Children {
			if int(m2.ID) == menuID {
				response.Result(c, m2.Children)
				return
			}
		}
	}

	response.Err(c, "无法获取子菜单信息")
}
