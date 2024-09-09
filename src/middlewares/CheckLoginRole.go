package middlewares

import (
	"encoding/json"
	"sports-common/log"
	"sports-common/request"
	"sports-common/tools"
	models "sports-models"
	"strings"

	"sports-admin/caches"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// 跳过登录角色验证路由
var skipRoleCheckUrls = []string{
	"/index/main",          // 登录后台主界面
	"/index/right",         // 登录后默认主页
	"/index/profile",       // 后台用户资料
	"/index/profile_save",  // 保存后台用户资料
	"/index/password",      // 后台用户密码
	"/index/password_save", // 保存后台用户密码
	"/index/logout",        // 退出登录
	"/index/google_code",   // google验证
	"/index/google_bind",   // google绑定
	"/index/qr",            // qr-code
	"/upload",              // 图片上传
	"/index/exchange",      // 汇率
	"/index/lang",          // 切换语言
	"/users/game_user_ids", //

	// -- 需要保存到菜单里面 -- 临时添加的将来要加入到菜单列表和权限控制里面
	"/receive_virtuals/save_virtual", //
	"/user_withdraws/save_config",    //
	"/users/import",                  // -- 临时加入导入权限
	"/users/used_ips",
	"/receive_virtuals",                 // 虚拟币
	"/receive_virtuals/create",          // 虚拟币 - 添加
	"/receive_virtuals/update",          // 虚拟币 - 修改
	"/receive_virtuals/delete",          // 虚拟币 - 删除
	"/receive_virtuals/state",           // 虚拟币 - 停启用
	"/receive_virtuals/save",            // 虚拟币 - 保存
	"/recieve_virtuals/set_float_rate",  // 虚拟币 - 浮动汇率
	"/agents/internal_insert",           // 代理内部调用
	"/user_bets/v2",                     // 投注测试
	"/user_bets/sync",                   // 投注测试 - 同步es -> pg
	"/user_bets/verify",                 // 投注测试 - 校对es -> pg
	"/user_bets/count_by_venues",        // 投注测试 - 同步es -> pg
	"/user_bets/count_by_users",         // 投注测试 - 同步es -> pg
	"/v2/user_audits/detail",            // 稽核详情 - 审核详情
	"/v3/user_audits/detail",            // 稽核详情 - 审核详情
	"/v2/user_audits/bets",              // 稽核详情 - 投注列表
	"/v2/user_audits/sync",              // 稽核详情 - 校对 - 同步
	"/user_deposit_coin_matches",        // 自动匹配
	"/user_deposit_coin_matches/create", // 自动匹配 - 添加
	"/user_deposit_coin_matches/save",   // 自动匹配 - 保存
	"/user_withdraw_coins",              // 代币提现
	"/user_withdraw_coin_hrs",           // 代币提现 - 历史记录
	"/admin_login_logs",                 // 后台登录日志
	"/users/tree_users",                 // 会员树形列表
	// "/users/level_up",                // 调整会员等级
	// "/users/level_up_save",           // 调整会员等级

	"/agent_domains",        // 代理域名绑定 - 列表
	"/agent_domains/create", // 代理域名绑定 - 添加
	"/agent_domains/update", // 代理域名绑定 - 修改
	"/agent_domains/save",   // 代理域名绑定 - 保存
	"/agent_domains/state",  // 代理域名绑定 - 状态
	"/agent_domains/delete", // 代理域名绑定 - 删除

	"/admin_tools/down_vips", // vip 降级 - 列表
	"/admin_tools/down_vip",  // vip 降级 - 处理
	"/admin_tools/up_vips",   // vip 降级 - 列表
	"/admin_tools/up_vip",    // vip 降级 - 处理
}

// checkMenu 递归检查菜单权限
func checkMenu(path string, menus []models.Menu, c *gin.Context, admin *models.LoginAdmin) bool { // 检测所有菜单
	checkWithQuery := func(realUrl string) bool { // 处理带?号的情况
		if withUrls := strings.Split(realUrl, "?"); len(withUrls) > 1 && withUrls[0] == path {
			return true
		}
		return false
	}
	checkWithSplit := func(realUrl string) bool { // 处理带 |号, 兼容 /notices/update|/notices/save 这种形式
		menuUrls := strings.Split(realUrl, "|") // 拆分多个路由
		for _, menuURL := range menuUrls {      // 遍历
			if menuURL == path {
				return true
			}
			if strings.Contains(menuURL, "?") { // 兼容 /configs/update?id=1 这种样式
				return checkWithQuery(menuURL)
			}
		}
		return false
	}
	platform := request.GetPlatform(c)
	for _, menu := range menus {
		if menu.Url == path || // 全匹配
			(strings.Contains(menu.Url, "|") && checkWithSplit(menu.Url)) || // 带|号
			(strings.Contains(menu.Url, "?") && checkWithQuery(menu.Url)) || // 带?号
			len(menu.Children) > 0 && checkMenu(path, menu.Children, c, admin) { // 有子菜单

			_, exists := c.Get("__access_log")
			// 表示: 在权限检测通过的情况下,需要写相关的用户日志
			if !exists && path != "/access_logs" && menu.IsEnableLog() { // 如果启用写日志, 则保存相关日志
				info := request.GetFingerPrint(c) // 浏览器指纹信息
				data := map[string]interface{}{
					"menu_id":    menu.Id,
					"menu_name":  menu.Name,
					"admin_id":   (*admin).Id,
					"admin_name": (*admin).Name,
					"created":    tools.NowMicro(),
					"log_type":   menu.LogType,
					"log_level":  menu.LogLevel,
					"path":       path,
					"ip":         c.ClientIP(),
					"method":     c.Request.Method,
					"data": func(r *gin.Context) string {
						postedData := func() map[string]interface{} {
							postData := map[string]interface{}{}
							for k, v := range r.Request.PostForm {
								postData[k] = v
							}
							return postData
						}()
						data := map[string]interface{}{
							"GET":  r.Request.URL.RawQuery,
							"POST": postedData,
							"BODY": request.GetPostedData(r),
						}

						bytes, err := json.Marshal(data)
						if err != nil {
							return err.Error()
						}
						return string(bytes)
					}(c),
					"operation_key": info,
				}
				if _, err := models.AccessLogs.Create(platform, data); err != nil {
					log.Logger.Error(err.Error())
				}
			}

			c.Set("__access_log", true)
			return true
		}
	}

	return false
}

// CheckLoginRole 检测是否有IP授权
func CheckLoginRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Host == "admin.ip.vhost" {
			c.Next()
			return
		}

		ip := c.ClientIP()
		admin, err := models.LoginAdmins.GetLoginByRequest(c)
		if err != nil {
			Render(c, "404.html", err.Error(), pongo2.Context{"ip": ip})
			return
		}
		if admin == nil { // 如果没有查找到此角色菜单权限
			Render(c, "404.html", "没有菜单权限访问当前页面或路由!", pongo2.Context{"ip": ip})
			return
		}

		//userIdStr, err := c.Cookie("user_id")
		//if err != nil {
		//	Render(c, "404.html", "没有菜单权限访问当前页面或路由~", pongo2.Context{"ip": ip})
		//	return
		//}
		//admin := session.Get(c, "admin").(models.LoginAdmin) // 登录角色
		platform := request.GetPlatform(c)
		adminRole := caches.AdminRoles.Get(platform, int(admin.RoleId))
		if adminRole == nil { // 如果没有查找到此角色菜单权限
			Render(c, "404.html", "没有菜单权限访问当前页面或路由~", pongo2.Context{"ip": ip})
			return
		}

		path := c.Request.URL.Path                                      // 请求路由
		if !tools.SliceStringContainsElement(skipRoleCheckUrls, path) { // 如果不存可跳过检测范围
			if !checkMenu(path, adminRole.Menus, c, admin) { // 如果此路由不在权限范围之内
				Render(c, "403.html", "缺少菜单权限访问当前页面或路由-", pongo2.Context{"ip": ip})
				return
			}
		}
		c.Next()
	}
}
