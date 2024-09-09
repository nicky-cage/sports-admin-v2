package caches

import (
	common "sports-common"
	models "sports-models"
	"strconv"
	"strings"

	"xorm.io/builder"
)

//
const keyAdminRoles = "admin_roles"

// CachedLocalAdminRoles 本地缓存
var CachedLocalAdminRoles = &map[uint32]models.AdminRole{}

// AdminRoles 后台角色
var AdminRoles = struct {
	Load func(string, ...bool)
	All  func(string) map[uint32]models.AdminRole
	Get  func(string, int) *models.AdminRole
}{
	Load: func(platform string, args ...bool) {
		adminRoles := map[uint32]models.AdminRole{}
		var rs []models.AdminRole
		Menus.Load(platform)

		loadAgain := false
		if len(args) > 0 && args[0] {
			loadAgain = args[0]
		}
		if !loadAgain {
			cache := common.Redis(platform)
			defer common.RedisRestore(platform, cache)
			cachedCount, _ := cache.Exists(CacheAdminPrefix + keyAdminRoles).Result()
			if cachedCount > 0 { // 如果被缓存过, 则不再缓存
				// 加载到内存当中
				adminRoles = map[uint32]models.AdminRole{}
				_ = getCache(platform, keyAdminRoles, &adminRoles)
				CachedLocalAdminRoles = &adminRoles // 本地缓存

				return
			}
		}

		// 列出所有角色
		err := models.AdminRoles.FindAllNoCount(platform, &rs, nil, "id ASC")
		if err != nil {
			return
		}

		for _, r := range rs { // 遍历所有角色
			r.Menus = getRoleMenus(platform, r.MenuIds)
			adminRoles[r.Id] = r
		}

		_ = setCache(platform, keyAdminRoles, adminRoles)
		CachedLocalAdminRoles = &adminRoles // 本地缓存
	},
	All: func(platform string) map[uint32]models.AdminRole {
		adminRoles := map[uint32]models.AdminRole{}
		_ = getCache(platform, keyAdminRoles, &adminRoles)
		return adminRoles
	},
	Get: func(platform string, id int) *models.AdminRole {
		if val, exists := (*CachedLocalAdminRoles)[uint32(id)]; exists {
			return &val
		}

		adminRoles := map[uint32]models.AdminRole{}
		_ = getCache(platform, keyAdminRoles, &adminRoles)
		if v, exists := adminRoles[uint32(id)]; exists {
			return &v
		}
		return nil
	},
}

// 列出所有级联菜单
func getRoleSubMenus(platform string, parentId uint32, ids map[int]int) []models.Menu {
	var rows []models.Menu
	var idArr []interface{}
	for k := range ids {
		idArr = append(idArr, k)
	}

	cond := builder.NewCond().And(builder.Eq{"parent_id": parentId}) // , "state": 2})
	cond = cond.And(builder.In("id", idArr...))
	if err := models.MenusHeap.FindAllNoCount(platform, &rows, cond, "sort DESC"); err == nil {
		for k, v := range rows { // 整理返回数据
			delete(ids, int(v.Id)) // 删除已经获取的id
			if v.Level <= menuMaxLevel {
				rows[k].Children = getRoleSubMenus(platform, v.Id, ids)
			}
		}
	}

	return rows
}

// getRoleMenus 得到角色菜单信息
func getRoleMenus(platform, ids string) []models.Menu {
	var menus []models.Menu
	if ids == "" { // 如果为空直接返回
		return menus
	}

	menuIDStrs := strings.Split(ids, ",")
	menuIds := map[int]int{}
	for _, idStr := range menuIDStrs { // 将id-string格式化为 []int{}
		if id, err := strconv.Atoi(idStr); err == nil && id > 0 {
			menuIds[id] = 0
		}
	}

	menus = getRoleSubMenus(platform, 0, menuIds)

	// sql := fmt.Sprintf("UPDATE admin_roles SET menu_ids = '%s' WHERE id = 78", strings.Trim(allIds, ","))
	// dbSession := common.Mysql(platform)
	// defer dbSession.Close()
	// dbSession.Exec(sql)

	// sortedMenus := Menus.List(platform)
	// for _, m := range sortedMenus {
	// 	m.Children = getRoleSubMenus(m.Id, allMenus, menuIds)
	// 	menus = append(menus, m)
	// }

	return menus
}
