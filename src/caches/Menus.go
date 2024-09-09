package caches

import (
	"fmt"
	common "sports-common"
	models "sports-models"

	"xorm.io/builder"
)

// 关联菜单
var keyRelatedMenus = "menus_related"

// 非关联菜单
var keynoAssocMenus = "menus_no_assoc"

// 关联后台菜单
var keyLayMenus = "menus_layui_tree"

// 菜单的最大层级
var menuMaxLevel uint8 = 6

//type MenuLogs struct {
//	MenuUrl    string //菜单url.path
//	MenuName   string //菜单名称
//	MenuId     uint32 //菜单名称Id
//	MenuNumber uint8  //菜单编号
//	EnableLog  bool   //是否开启日志
//}

//var MenusUrlPathToMenuNameMap map[string]MenuLogs

// Menus 菜单
var Menus = struct {
	Load           func(string)                           //加载初始菜单
	Get            func(string, int) *models.Menu         //得到单个菜单
	All            func(string) map[uint32]models.Menu    // 获取非关联的菜单
	List           func(string) []models.Menu             //获取关联菜单
	LayMenus       func(platform string) []models.LayMenu //
	LayMenusByJson func(platform string) string           //
}{
	Load: func(platform string) {

		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		_, _ = dbSession.Exec("DROP TABLE IF EXISTS menus_in_heap") // 每次启动要清除之前的缓存
		rows, err := dbSession.QueryString("SELECT id FROM menus_in_heap LIMIT 1")
		if err != nil { // 如果有错误, 则新建表
			_, _ = dbSession.Exec("CREATE TABLE IF NOT EXISTS menus_in_heap ENGINE=MEMORY SELECT * FROM menus")
			_, _ = dbSession.Exec("ALTER TABLE menus_in_heap ADD PRIMARY KEY(id)")
		} else if len(rows) == 0 { // 如果数据为0则新建数据
			_, _ = dbSession.Exec("INSERT INTO menus_in_heap SELECT * FROM menus")
		}

		cache := common.Redis(platform)
		defer common.RedisRestore(platform, cache)

		// cotNoAssoc, _ := cache.Exists(CacheAdminPrefix + keynoAssocMenus).Result()
		// cotRelated, _ := cache.Exists(CacheAdminPrefix + keynoAssocMenus).Result()
		// cotLayMenus, _ := cache.Exists(CacheAdminPrefix + keynoAssocMenus).Result()
		// if cotNoAssoc > 0 && cotRelated > 0 && cotLayMenus > 0 {
		// 	return
		// }

		var relatedMenus []models.Menu           // 关联菜单
		noAssocMenus := map[uint32]models.Menu{} // 非关联菜单
		var layMenus []models.LayMenu

		// 获取一级菜单
		var allMenus []models.Menu
		cond := builder.NewCond().And(builder.Eq{"parent_id": 0}).And(builder.Eq{"state": 2}) //
		if err := models.MenusHeap.FindAllNoCount(platform, &allMenus, cond, "sort DESC"); err == nil {
			for _, v := range allMenus {
				noAssocMenus[v.Id] = v                                       //非关联菜单 - 累加
				rMenus, lMenus := getSubMenus(platform, v.Id, &noAssocMenus) // 递归显示子菜单
				v.Children = rMenus                                          // 关联菜单 - 子菜单
				relatedMenus = append(relatedMenus, v)                       // 关联菜单
				layMenus = append(layMenus, models.LayMenu{
					ID:       v.Id,
					Title:    v.Name,
					Field:    fmt.Sprintf("m%d_%d", 1, v.Id),
					Spread:   true,
					Children: lMenus,
				})
			}
		}

		_ = setCache(platform, keynoAssocMenus, noAssocMenus) // 非关联菜单
		_ = setCache(platform, keyRelatedMenus, relatedMenus) // 关联菜单
		_ = setCache(platform, keyLayMenus, layMenus)
	},
	Get: func(platform string, id int) *models.Menu {
		noAssocMenus := map[uint32]models.Menu{} // 非关联菜单
		_ = getCache(platform, keynoAssocMenus, &noAssocMenus)
		if v, exists := noAssocMenus[uint32(id)]; exists {
			return &v
		}
		return nil
	},
	All: func(platform string) map[uint32]models.Menu { // 非关联菜单
		noAssocMenus := map[uint32]models.Menu{} // 非关联菜单
		_ = getCache(platform, keynoAssocMenus, &noAssocMenus)
		return noAssocMenus
	},
	List: func(platform string) []models.Menu { // 得到关联菜单
		var relatedMenus []models.Menu // 关联菜单
		_ = getCache(platform, keyRelatedMenus, &relatedMenus)
		return relatedMenus
	},
	LayMenus: func(platform string) []models.LayMenu {
		var menus []models.LayMenu
		_ = getCache(platform, keyLayMenus, &menus)
		return menus
	},
	LayMenusByJson: func(platform string) string {
		return getCacheByJson(platform, keyLayMenus)
	},
}

// 列出所有级联菜单
func getSubMenus(platform string, parentId uint32, noKeyArr *map[uint32]models.Menu) ([]models.Menu, []models.LayMenu) {
	var rows []models.Menu                                           // 关联菜单
	var lRows []models.LayMenu                                       // 非关联菜单
	cond := builder.NewCond().And(builder.Eq{"parent_id": parentId}) // 指定上级id
	if err := models.MenusHeap.FindAllNoCount(platform, &rows, cond, "sort DESC"); err == nil {
		for k, r := range rows {
			(*noKeyArr)[r.Id] = r // 写入非关联菜单
			lRows = append(lRows, models.LayMenu{
				ID:     r.Id,
				Title:  r.Name,
				Field:  fmt.Sprintf("m%d_%d", r.Level, r.Id),
				Spread: true,
			})
			if r.Level <= menuMaxLevel {
				rMenus, lMenus := getSubMenus(platform, r.Id, noKeyArr)
				rows[k].Children = rMenus
				lRows[k].Children = lMenus
			}
		}
	}
	return rows, lRows
}
