package router

import (
	"os"
	"pongo2gin"
	"sports-common/response"
	"sports-common/tools"
	"sports-common/ws"

	"sports-admin/controllers"
	"sports-admin/controllers/admin_login_logs"
	"sports-admin/controllers/admin_tools"
	"sports-admin/controllers/user_relations"
	"sports-admin/filters"
	"sports-admin/functions"
	"sports-admin/libs"
	"sports-admin/middlewares"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

// 默认模板目录
// 判断当前目录下是否存在模板文件目录
var templateDir = "../templates/"

// Init 路由初始化
func Init(router *gin.Engine) {
	router.GET("/add_ip", controllers.Index.AddIP)                 // 增加ip
	router.POST("/add_ip_save", controllers.Index.AddIPSave)       // 增加ip
	router.POST("/send_mail_code", controllers.Index.SendMailCode) // 增加ip - 发送验证码
	router.POST("/index/captcha", controllers.Index.Captcha)       // 生成图形验证码

	router.Use(middlewares.DomainCheck())       // 检测非法域名访问
	router.Use(middlewares.CheckPermissionIp()) // 检测ip是否有使用权限
	// router.Use(middlewares.Cors()) //跨域

	// websocket - 实时返回服务器信息
	ws.WebSocket.Start()
	ws.WebSocket.AddRouters(map[int]func(*melody.Session, interface{}) interface{}{
		801001: controllers.Index.Overload, // 系统信息检测
	})
	ws.WebSocket.SetEndPoint(ws.EndPointTypes.Backend)
	ws.WebSocket.SetConnectID(libs.WebSocketConfig.GetConnectID)
	ws.WebSocket.InitNodes() // 加载各个节点信息
	router.GET("/ws", ws.WebSocket.Handle)

	filters.InitFilters() // 设置过滤函数
	response.LoadTemplateFuncs(functions.All())

	// 如果当前目录下有模板目录, 则先取当前目录之下
	if s, err := os.Stat("./templates"); err == nil && s.IsDir() {
		templateDir = "./templates"
	}
	renderOptions := pongo2gin.RenderOptions{
		TemplateDir: templateDir,                // 模板文件目录, 可以修改为自己的模板文件夹位置, 请使用相对路径
		ContentType: "text/html; charset=utf-8", // 输出文件类型, 一般不要改动
	}
	router.HTMLRender = pongo2gin.New(renderOptions)

	noAuth := router.Group("/")
	{
		noAuth.GET("/", controllers.Index.Index)                         // 默认首页
		noAuth.GET("/index/test", controllers.Index.Test)                // 测试页面
		noAuth.POST("/index/login", controllers.Index.Login)             // 测试页面
		noAuth.GET("/examples", controllers.Examples.List)               // 示例
		noAuth.GET("/payments/notify", controllers.UserWithdraws.Notify) // 支付平台出款回调
	}

	auth := router.Group("/")
	{
		auth.Use(middlewares.CheckLogin())     // 检测用户是否登录
		auth.Use(middlewares.CheckLoginRole()) // 检测用户是否登录
		auth.POST("/upload", tools.UploadFile) // 上传文件-通用

		auth.GET("/index/main", controllers.Index.Main)                   // 后台界面
		auth.GET("/index/right", controllers.Index.Right)                 // 后台首页
		auth.GET("/index/profile", controllers.Index.Profile)             // 用户资料
		auth.POST("/index/logout", controllers.Index.Logout)              // 后台界面 - 退出登录
		auth.POST("/index/profile_save", controllers.Index.ProfileSave)   // 用户资料
		auth.GET("/index/password", controllers.Index.Password)           // 修改密码
		auth.GET("/index/qr", controllers.Index.QRCode)                   // 二维码
		auth.POST("/index/password_save", controllers.Index.PasswordSave) // 修改密码
		auth.GET("/index/google_code", controllers.Index.GoogleCode)      // 谷歌验证
		auth.POST("/index/google_bind", controllers.Index.GoogleBind)     // 绑定谷歌验证
		auth.GET("/index/exchange", controllers.Index.Exchange)           // 汇率获取

		// -- 临时添加
		auth.GET("/receive_virtuals", controllers.ReceiveVirtuals.List)                         // 虚拟币
		auth.GET("/receive_virtuals/create", controllers.ReceiveVirtuals.Create)                // 虚拟币 - 添加
		auth.GET("/receive_virtuals/update", controllers.ReceiveVirtuals.Update)                // 虚拟币 - 修改
		auth.GET("/receive_virtuals/delete", controllers.ReceiveVirtuals.Delete)                // 虚拟币 - 删除
		auth.GET("/receive_virtuals/state", controllers.ReceiveVirtuals.State)                  // 虚拟币 - 停启用
		auth.POST("/receive_virtuals/save", controllers.ReceiveVirtuals.Save)                   // 虚拟币 - 保存
		auth.POST("/receive_virtuals/set_float_rate", controllers.ReceiveVirtuals.SetFloatRate) // 虚拟币 - 浮动汇率

		auth.GET("/user_detail/index", controllers.UserDetailAccounts.Index)                        // 用户详情 - 基本信息
		auth.GET("/user_detail/accounts", controllers.UserDetailAccounts.Accounts)                  // 用户详情 - 账户信息
		auth.POST("/user_detail/accounts/transfer_out", controllers.UserDetailAccounts.TransferOut) // 用户详情 - 账户信息
		auth.GET("/user_detail/accounts_recovery", controllers.UserDetailAccounts.Recovery)         // 用户详情
		auth.GET("/user_detail/account_async", controllers.UserDetailAccounts.AccountAsync)         // 用户详情 - 账户信息
		auth.GET("/user_detail/wins", controllers.UserDetailAccounts.Wins)                          // 用户详情 - 输赢信息
		auth.GET("/user_detail/deposits", controllers.UserDetailAccounts.Deposits)                  // 用户详情 - 存款信息
		auth.GET("/user_detail/withdraws", controllers.UserDetailAccounts.Withdraws)                // 用户详情 - 提款信息
		auth.GET("/user_detail/dividends", controllers.UserDetailAccounts.Dividends)                // 用户详情 - 红利信息
		auth.GET("/user_detail/commissions", controllers.UserDetailAccounts.Commissions)            // 用户详情 - 佣金信息
		auth.GET("/user_detail/transfers", controllers.UserDetailAccounts.Transfers)                // 用户详情 - 转账信息
		auth.GET("/user_detail/resets", controllers.UserDetailResets.Index)                         // 用户详情 - 账户调整
		auth.POST("/user_detail/resets/save", controllers.UserDetailResets.SaveMore)                // 用户详情 -
		// auth.POST("/user_detail/resets/savemore", controllers.UserDetailResets.SaveMore)                  // 用户详情 -
		auth.GET("/user_detail/changes", controllers.UserDetailAccounts.Changes)                          // 用户详情 - 账变记录
		auth.GET("/user_detail/commission_accounts", controllers.UserDetailCommissions.Accounts)          // 用户详情 - 账户信息
		auth.GET("/user_detail/commission_account_async", controllers.UserDetailCommissions.AccountAsync) // 用户详情 - 同步账户信息
		auth.GET("/user_detail/commission_withdraws", controllers.UserDetailCommissions.Withdraws)        // 用户详情 - 提款信息
		auth.GET("/user_detail/commission_records", controllers.UserDetailCommissions.Records)            // 用户详情 - 佣金信息
		auth.GET("/user_detail/account", controllers.UserDetailCommissions.Detail)                        // 用户详情 - 佣金信息
		auth.GET("/user_levels", controllers.UserLevels.List)                                             // 会员管理-会员等级
		auth.GET("/user_levels/update", controllers.UserLevels.Update)                                    // 会员管理-会员等级-设置
		auth.POST("/user_levels/save", controllers.UserLevels.Save)                                       // 会员管理-会员等级-设置
		auth.GET("/user_login_logs", controllers.UserLoginLogs.List)                                      // 默认首页
		auth.GET("/user/real_name/history", controllers.RealNameHistories.List)                           // 用户详情 - 银行卡姓名修改历史
		auth.GET("/sites", controllers.Sites.Index)                                                       // 站点配置
		auth.POST("/users/show_im", controllers.Users.ShowIM)                                             // 显示用户信息
		auth.GET("/users/export", controllers.Users.Export)                                               // 导出会员信息
		auth.POST("/users/export", controllers.Users.Export)                                              // 导出会员信息
		auth.POST("/users/import", controllers.Users.ImportExcel)                                         // 导入会员信息
		auth.GET("/users/game_user_ids", controllers.Users.GameUserIds)                                   // 游戏用户列表ID
		// auth.GET("/users", controllers.Users.Index)                     //默认首页

		// 用户验证码查询
		auth.GET("/users_code", controllers.UserCodes.List)
		auth.GET("/user_access_logs", controllers.UserAccessLogs.List) // 后台日志

		auth.GET("/configs/update", controllers.Configs.Update) // 站点配置 - 界面
		auth.POST("/configs/save", controllers.Configs.Save)    // 站点配置 - 保存信息

		auth.GET("/popular", controllers.ConfigPopulars.List)       // 站点配置 - 界面
		auth.POST("/popular/save", controllers.ConfigPopulars.Save) // 站点配置 - 界面
		auth.GET("/notices", controllers.Notices.List)              // 系统公告
		auth.GET("/notices/create", controllers.Notices.Create)     // 系统公告 - 添加
		auth.GET("/notices/update", controllers.Notices.Update)     // 系统公告 - 修改
		auth.POST("/notices/save", controllers.Notices.Save)        // 系统公告 - 保存
		auth.GET("/notices/delete", controllers.Notices.Delete)     // 系统公告 - 删除
		auth.GET("/notices/state", controllers.Notices.State)       // 系统公告 - 停启用
		//
		auth.GET("/user_feedback", controllers.UserFeedBack.List)           // 系统公告 - 停启用
		auth.GET("/user_feedback/update", controllers.UserFeedBack.Update)  // 系统公告 - 停启用
		auth.POST("/user_feedback/Save", controllers.UserFeedBack.Save)     // 系统公告 - 停启用
		auth.GET("/user_feedback/delete", controllers.UserFeedBack.Delete)  // 系统公告 - 停启用
		auth.GET("/messages", controllers.Messages.List)                    // 站内消息
		auth.GET("/messages/create", controllers.Messages.Create)           // 站内消息 - 添加
		auth.GET("/messages/update", controllers.Messages.Update)           // 站内消息 - 修改
		auth.POST("/messages/save", controllers.Messages.Save)              // 站内消息 - 保存
		auth.GET("/messages/delete", controllers.Messages.Delete)           // 站内消息 - 删除
		auth.GET("/messages/state", controllers.Messages.State)             // 站内消息 - 停启用
		auth.GET("/messages/top", controllers.Messages.Top)                 // 站内消息 - 置顶
		auth.GET("/message_agent", controllers.MessageAgent.List)           // 代理站内信
		auth.GET("/message_agent/created", controllers.MessageAgent.Create) // 代理站内信
		auth.GET("/message_agent/updated", controllers.MessageAgent.Update) // 代理站内信
		auth.POST("/message_agent/save", controllers.MessageAgent.Save)     // 代理站内信
		auth.GET("/message_agent/deleted", controllers.MessageAgent.Delete) // 代理站内信

		// 社区爆料
		auth.GET("/recommend", controllers.Recommends.List)                     // 社区爆料
		auth.GET("/recommend/created", controllers.Recommends.Create)           // 社区爆料方案创建
		auth.GET("/recommend/updated", controllers.Recommends.Update)           // 社区爆料
		auth.POST("/recommend/save", controllers.Recommends.Save)               // 社区爆料
		auth.GET("/recommend/delete", controllers.Recommends.Delete)            // 社区爆料
		auth.GET("/recommend_nickname/created", controllers.Recommends.Created) // 社区爆料
		auth.POST("/recommend_nickname/save", controllers.Recommends.Saved)     // 社区爆料

		// 体育资讯
		auth.GET("/sport_news", controllers.SportNews.List)           // 体育资讯
		auth.POST("/sport_news/save", controllers.SportNews.Save)     // 体育资讯
		auth.GET("/sport_news/created", controllers.SportNews.Create) // 体育资讯
		auth.GET("/sport_news/updated", controllers.SportNews.Update) // 体育资讯
		auth.GET("/sport_news/deleted", controllers.SportNews.Delete) // 体育资讯

		// 赛事匹配
		auth.GET("/match_game", controllers.MatchGames.Created)         // 赛事未匹配
		auth.POST("/match_game/risk", controllers.MatchGames.Save)      // 赛事未匹配
		auth.GET("/match_games", controllers.MatchGames.List)           // 赛事已匹配
		auth.GET("/match_games/deleted", controllers.MatchGames.Delete) // 赛事已匹配

		auth.GET("/ads", controllers.Ads.List)          // 广告管理-app启动
		auth.GET("/ads/create", controllers.Ads.Create) // app启动 - 添加
		auth.GET("/ads/update", controllers.Ads.Update) // app启动 - 修改
		auth.POST("/ads/save", controllers.Ads.Save)    // app启动 - 保存
		auth.GET("/ads/delete", controllers.Ads.Delete) // app启动 - 删除
		auth.GET("/ads/state", controllers.Ads.State)   // app启动 - 停启用

		auth.GET("/ad_carousels", controllers.AdCarousels.List)          // 广告管理-轮播图
		auth.GET("/ad_carousels/create", controllers.AdCarousels.Create) // 轮播图 - 添加
		auth.GET("/ad_carousels/update", controllers.AdCarousels.Update) // 轮播图 - 修改
		auth.POST("/ad_carousels/save", controllers.AdCarousels.Save)    // 轮播图 - 保存
		auth.GET("/ad_carousels/delete", controllers.AdCarousels.Delete) // 轮播图 - 删除
		auth.GET("/ad_carousels/state", controllers.AdCarousels.State)   // 轮播图 - 停启用

		auth.GET("/ad_matches", controllers.AdMatches.List)          // 广告管理-体育赛事
		auth.GET("/ad_matches/create", controllers.AdMatches.Create) // 体育赛事 - 添加
		auth.GET("/ad_matches/update", controllers.AdMatches.Update) // 体育赛事 - 修改
		auth.POST("/ad_matches/save", controllers.AdMatches.Save)    // 体育赛事 - 保存
		auth.GET("/ad_matches/delete", controllers.AdMatches.Delete) // 体育赛事 - 删除
		auth.GET("/ad_matches/state", controllers.AdMatches.State)   // 体育赛事 - 停启用

		auth.GET("/ad_sponsors", controllers.AdSponsors.List)          // 广告管理-赞助配置
		auth.GET("/ad_sponsors/create", controllers.AdSponsors.Create) // 赞助配置 - 添加
		auth.GET("/ad_sponsors/update", controllers.AdSponsors.Update) // 赞助配置 - 修改
		auth.POST("/ad_sponsors/save", controllers.AdSponsors.Save)    // 赞助配置 - 保存
		auth.GET("/ad_sponsors/delete", controllers.AdSponsors.Delete) // 赞助配置 - 删除
		auth.GET("/ad_sponsors/state", controllers.AdSponsors.State)   // 赞助配置 - 停启用

		auth.GET("/agents", controllers.Agents.ListsV2)                   // 代理管理
		auth.GET("/agents/lists", controllers.Agents.List)                // 代理联表列表
		auth.GET("/agents/detail_view", controllers.Agents.DetailView)    // 代理管理
		auth.GET("/agents/update", controllers.Agents.Updated)            // 代理管理
		auth.POST("/agents/save", controllers.Agents.Save)                // 代理管理
		auth.POST("/agents/lower_add", controllers.Agents.LowerAdd)       // 代理管理
		auth.POST("/agents/insert", controllers.Agents.Insert)            // 新增代理
		auth.GET("/agents/create", controllers.Agents.Create)             // 代理管理
		auth.GET("/agents/add", controllers.Agents.Add)                   // 代理管理
		auth.GET("/agents/detail", controllers.Agents.Detail)             // 代理管理
		auth.GET("/agents/users", controllers.AgentUsers.Lists)           // 代理管理
		auth.GET("/agents/users/detail", controllers.AgentUsers.Details)  // 代理管理
		auth.GET("/agents/users/update", controllers.AgentUsers.Update)   // 代理管理
		auth.POST("/agents/users_save", controllers.AgentUsers.Save)      // 代理管理
		auth.GET("/agents/user_check", controllers.AgentUsers.CheckAgent) // 代理管理

		auth.GET("/agents/commissions", controllers.AgentCommissions.Lists)                      // 代理佣金列表
		auth.GET("/agents/commissions/list", controllers.AgentCommissions.List)                  // 代理佣金联表
		auth.GET("/agents/commissions/plan", controllers.AgentCommissionsPlan.Index)             // 代理佣金方案
		auth.POST("/agents/commissions/plan/add", controllers.AgentCommissionsPlan.Add)          // 代理佣金方案
		auth.POST("/agents/commissions/plan/save", controllers.AgentCommissionsPlan.Saved)       // 代理佣金方案
		auth.GET("/agents/commissions/plan/updated", controllers.AgentCommissionsPlan.Updated)   // 代理佣金方案
		auth.GET("/agents/commissions/plan/deleted", controllers.AgentCommissionsPlan.Deleted)   // 代理佣金方案
		auth.GET("/agents/commissions/plan/view", controllers.AgentCommissionsPlan.View)         // 代理佣金方案
		auth.GET("/agents/commissions/plan_list", controllers.AgentCommissionsPlan.Lists)        // 代理佣金方案
		auth.GET("/agents/commissions/create", controllers.AgentCommissionsPlan.Create)          // 代理佣金方案
		auth.GET("/agents/commissions/create_mode", controllers.AgentCommissionsPlan.CreateMode) // 代理佣金方案
		auth.GET("/agents/commissions/record", controllers.AgentCommissionsLogs.List)            // 代理佣金记录
		auth.GET("/agents/commissions/adjustment", controllers.AgentCommissions.Create)          // 代理佣金调整
		auth.POST("/agents/commissions/save", controllers.AgentCommissions.Save)                 // 代理佣金调整保存
		auth.POST("/agents/commissions/grant", controllers.AgentCommissions.Grant)               // 代理佣金发送
		auth.GET("/agents/withdraws", controllers.AgentWithdraws.List)                           // 代理提款审核
		auth.POST("/agents/withdraws/saves", controllers.AgentWithdraws.Saves)                   // 代理提款审核
		auth.GET("/agents/withdraws/record", controllers.AgentWithdrawsRecords.List)             // 代理提款审核记录

		auth.GET("/agent_domains", controllers.AgentDomains.List)          // 代理绑定域名- 列表
		auth.GET("/agent_domains/create", controllers.AgentDomains.Create) // 代理绑定域名- 添加
		auth.GET("/agent_domains/update", controllers.AgentDomains.Update) // 代理绑定域名- 修改
		auth.GET("/agent_domains/delete", controllers.AgentDomains.Delete) // 代理绑定域名- 删除
		auth.GET("/agent_domains/state", controllers.AgentDomains.State)   // 代理绑定域名- 状态修改
		auth.POST("/agent_domains/save", controllers.AgentDomains.Save)    // 代理绑定域名- 保存

		auth.GET("/agents/logs", controllers.AgentLogs.Games)                        // 代理游戏记录
		auth.GET("/agents/logs_adjust", controllers.AgentLogs.Adjust)                // 代理输赢调整记录
		auth.GET("/agents/logs_deposits", controllers.AgentLogs.Deposits)            // 代理存款记录
		auth.GET("/agents/logs_login", controllers.AgentLogs.List)                   // 代理登录记录
		auth.GET("/agents/logs_transfer", controllers.AgentLogs.Transfer)            // 代理管理
		auth.GET("/deposit_discounts", controllers.DepositDiscounts.List)            // 存款优惠
		auth.GET("/deposit_discounts/edit", controllers.DepositDiscounts.Edit)       // 存款优惠-设置
		auth.POST("/deposit_discounts/save_do", controllers.DepositDiscounts.SaveDo) // 存款优惠-设置-保存
		auth.GET("/deposit_discounts/state", controllers.DepositDiscounts.State)     // 存款优惠 - 停启用

		auth.GET("/black_lists", controllers.BlockedDevices.List)              // 黑名单 - 设备编号 - 列表
		auth.GET("/blocked_devices/create", controllers.BlockedDevices.Create) // 黑名单 - 设备编号 - 添加
		auth.GET("/blocked_devices/update", controllers.BlockedDevices.Update) // 黑名单 - 设备编号 - 修改
		auth.POST("/blocked_devices/save", controllers.BlockedDevices.Save)    // 黑名单 - 设备编号 - 保存
		auth.GET("/blocked_devices/delete", controllers.BlockedDevices.Delete) // 黑名单 - 设备编号 - 删除

		auth.GET("/blocked_ips", controllers.BlockedIps.List)          // 黑名单 - IP - 列表
		auth.GET("/blocked_ips/create", controllers.BlockedIps.Create) // 黑名单 - IP - 添加
		auth.GET("/blocked_ips/update", controllers.BlockedIps.Update) // 黑名单 - IP - 修改
		auth.POST("/blocked_ips/save", controllers.BlockedIps.Save)    // 黑名单 - IP - 保存
		auth.GET("/blocked_ips/delete", controllers.BlockedIps.Delete) // 黑名单 - IP - 删除

		auth.GET("/blocked_mails", controllers.BlockedMails.List)          // 黑名单 - 电子邮件 - 列表
		auth.GET("/blocked_mails/create", controllers.BlockedMails.Create) // 黑名单 - 电子邮件 - 添加
		auth.GET("/blocked_mails/update", controllers.BlockedMails.Update) // 黑名单 - 电子邮件 - 修改
		auth.POST("/blocked_mails/save", controllers.BlockedMails.Save)    // 黑名单 - 电子邮件 - 保存
		auth.GET("/blocked_mails/delete", controllers.BlockedMails.Delete) // 黑名单 - 电子邮件 - 删除

		auth.GET("/blocked_phones", controllers.BlockedPhones.List)          // 黑名单 - 手机号码 - 列表
		auth.GET("/blocked_phones/create", controllers.BlockedPhones.Create) // 黑名单 - 手机号码 - 添加
		auth.GET("/blocked_phones/update", controllers.BlockedPhones.Update) // 黑名单 - 手机号码 - 修改
		auth.POST("/blocked_phones/save", controllers.BlockedPhones.Save)    // 黑名单 - 手机号码 - 保存
		auth.GET("/blocked_phones/delete", controllers.BlockedPhones.Delete) // 黑名单 - 手机号码 - 删除

		auth.GET("/blocked_cards", controllers.BlockedCards.List)          // 黑名单 - 银行卡号 - 列表
		auth.GET("/blocked_cards/create", controllers.BlockedCards.Create) // 黑名单 - 银行卡号 - 添加
		auth.GET("/blocked_cards/update", controllers.BlockedCards.Update) // 黑名单 - 银行卡号 - 修改
		auth.POST("/blocked_cards/save", controllers.BlockedCards.Save)    // 黑名单 - 银行卡号 - 保存
		auth.GET("/blocked_cards/delete", controllers.BlockedCards.Delete) // 黑名单 - 银行卡号 - 删除

		auth.GET("/risk_audits", controllers.RiskAudits.List)                     // 提款审核列表
		auth.GET("/risk_audits/export", controllers.RiskAudits.Export)            // 提款审核列表 - 导出
		auth.POST("/risk_audits/export", controllers.RiskAudits.Export)           // 提款审核列表 - 导出
		auth.GET("/risk_audits/created", controllers.RiskAudits.Create)           // 提款审核列表
		auth.GET("/risk_audits/state", controllers.RiskAudits.State)              // 提款审核列表
		auth.POST("/risk_audits/created_save", controllers.RiskAudits.CreateSave) // 提款审核列表
		// auth.GET("/risk_audits/list", controllers.RiskAudits.List)                //提款审核
		auth.GET("/risk_audits/detail", controllers.RiskAudits.Detail)             // 提款审核详情
		auth.POST("/risk_audits/save", controllers.RiskAudits.Save)                // 提款审核更改
		auth.GET("/risk_audits/detail_view", controllers.RiskAudits.View)          // 提款审核详情视图
		auth.GET("/risk_audits/sys_audits", controllers.RiskAudits.SysAudits)      // 系统审核
		auth.GET("/risk_audits/sys_detail", controllers.RiskAudits.SysDetail)      // 系统审核详情
		auth.GET("/risk_audits/refuse", controllers.RiskAudits.Refuse)             // 提款拒绝
		auth.GET("/risk_audits/hand_up", controllers.RiskAudits.HandUp)            // 提款挂起
		auth.GET("/risk_audits/receive", controllers.RiskAudits.Receive)           // 人工审核
		auth.GET("/risk_audits/hand_up/list", controllers.RiskAuditsList.Lists)    // 审核挂起列表展示
		auth.POST("/risk_audits/receive_save", controllers.RiskAudits.ReceiveSave) // 审核人工领取
		auth.POST("/risk_audits/saves", controllers.RiskAudits.Saves)              // 审核，拒绝，通过，挂起
		auth.GET("/risk_audits/history", controllers.RiskAuditsLog.List)           // 提款审核历史记录
		auth.GET("/risk_audits/history_detail", controllers.RiskAuditsLog.Update)  // 提款审核历史记录

		auth.GET("/users", controllers.Users.List)                                 // 用户管理 - 列表
		auth.GET("/users/create", controllers.Users.Create)                        // 用户管理 - 添加
		auth.POST("/users/add", controllers.Users.Add)                             // 用户管理 - 添加
		auth.GET("/users/update", controllers.Users.Update)                        // 用户管理 - 修改
		auth.GET("/users/update/address", controllers.Users.Address)               // 用户管理 - 修改
		auth.GET("/users/state", controllers.Users.State)                          // 用户管理 - 状态
		auth.POST("/users/save", controllers.Users.Save)                           // 用户管理 - 保存
		auth.GET("/users/detail", controllers.Users.Detail)                        // 用户管理 - 详情
		auth.GET("/users/password", controllers.Users.Password)                    // 用户管理 - 修改密码
		auth.POST("/users/password_save", controllers.Users.PasswordSave)          // 用户管理 - 保存密码
		auth.POST("/users/disable_all", controllers.Users.DisableAll)              // 用户管理 - 批量禁用
		auth.GET("/users/add_tags", controllers.Users.AddTags)                     // 用户管理 - 批量添加标签
		auth.POST("/users/add_tags_save", controllers.Users.AddTagsSave)           // 用户管理 - 批量添加标签 - 保存
		auth.POST("/users/withdraw_password", controllers.Users.WithdrawPassword)  // 用户管理 - 列表
		auth.GET("/user_bets", controllers.UserBets.List)                          // 用户管理 - 注单
		auth.GET("/user_bets/v1", controllers.UserBets.ListV1)                     // 用户管理 - 注单
		auth.GET("/user_bets/sync", controllers.UserBets.Sync)                     // 用户管理 - 注单 - 同步
		auth.GET("/user_bets/verify", controllers.UserBets.Verify)                 // 用户管理 - 注单 - 校对
		auth.GET("/user_bets/count_by_venues", controllers.UserBets.CountByVenues) // 用户管理 - 注单 - 同步
		auth.GET("/user_bets/count_by_users", controllers.UserBets.CountByUsers)   // 用户管理 - 注单 - 同步
		auth.GET("/user_bets/set_up", controllers.UserBets.SetUp)                  // 用户管理 - 手动补单页面
		auth.POST("/user_bets/save_do", controllers.UserBets.SaveDo)               // 用户管理 - 手动补单
		auth.GET("/user_bets/detail", controllers.UserBets.Detail)                 // 用户管理 - 注单详情
		auth.GET("/users/used_ips", controllers.Users.Ips)                         // 用户管理 - ip记录
		auth.GET("/users/level_up", controllers.Users.LevelUp)                     // 用户管理 - 修改vip等级
		auth.POST("/users/level_up_save", controllers.Users.LevelUpSave)           // 用户管理 - 修改vip等级
		auth.GET("/users/tree_users", user_relations.UserRelations.TreeUsers)      // 用户管理 - 用户节点

		// auth.GET("/banks", controllers.Banks.Index) //银行列表
		auth.GET("/banks", controllers.Banks.List)          // 银行列表 - 列表
		auth.GET("/banks/create", controllers.Banks.Create) // 银行列表 - 添加
		auth.GET("/banks/update", controllers.Banks.Update) // 银行列表 - 修改
		auth.POST("/banks/save", controllers.Banks.Save)    // 银行列表 - 保存
		auth.GET("/banks/state", controllers.Banks.State)   // 银行列表 - 状态

		auth.GET("/help_categories/create", controllers.HelpCategories.Create) // 帮助分类
		auth.POST("/help_categories/save", controllers.HelpCategories.Save)    // 帮助分类

		auth.GET("/helps", controllers.Helps.List)          // 帮助列表
		auth.GET("/helps/create", controllers.Helps.Create) // 帮助列表 - 添加
		auth.GET("/helps/update", controllers.Helps.Update) // 帮助列表 - 修改
		auth.POST("/helps/save", controllers.Helps.Save)    // 帮助列表 - 保存
		auth.GET("/helps/delete", controllers.Helps.Delete) // 帮助列表 - 删除
		auth.GET("/helps/state", controllers.Helps.State)   // 帮助列表 - 状态
		auth.GET("/helps/check", controllers.Helps.Check)   // 帮助列表 - 搜索
		auth.GET("/helps/detail", controllers.Helps.Detail) // 帮助列表 - 预览

		auth.GET("/user_cards", controllers.UserCards.List)                         // 会员绑卡 - 默认首页
		auth.GET("/user_cards/detail", controllers.UserCards.Detail)                // 会员绑卡 - 默认首页
		auth.GET("/user_cards/create", controllers.UserCards.Create)                // 会员绑卡 - 新增页面
		auth.GET("/user_cards/update", controllers.UserCards.Update)                // 会员绑卡 - 编辑页面
		auth.POST("/user_cards/save", controllers.UserCards.Save)                   // 会员绑卡 - 保存
		auth.GET("/user_cards/delete", controllers.UserCards.Delete)                // 会员绑卡 - 删除
		auth.GET("/user_virtual_coins", controllers.UserVirtualCoins.List)          // 会员虚拟货币 - 默认首页
		auth.GET("/user_virtual_coins/create", controllers.UserVirtualCoins.Create) // 会员虚拟货币 - 新增页面
		auth.GET("/user_virtual_coins/update", controllers.UserVirtualCoins.Update) // 会员虚拟货币 - 编辑页面
		auth.GET("/user_virtual_coins/delete", controllers.UserVirtualCoins.Delete) // 会员虚拟货币 - 删除
		auth.POST("/user_virtual_coins/save", controllers.UserVirtualCoins.Save)    // 会员虚拟货币 - 保存

		auth.GET("/user_audits", controllers.UserAudits.List)               // 会员稽核 - 默认首页
		auth.GET("/user_audits/detail", controllers.UserAudits.Detail)      // 会员稽核 - 详情页面
		auth.GET("/v2/user_audits/detail", controllers.UserAudits.DetailV2) // 会员稽核 - 详情页面
		auth.GET("/v3/user_audits/detail", controllers.UserAudits.DetailV3) // 会员稽核 - 详情页面
		auth.GET("/v2/user_audits/bets", controllers.UserAudits.BetsV2)     // 会员稽核 - 详情页面
		auth.GET("/v2/user_audits/sync", controllers.UserAudits.Sync)       // 会员稽核 - 同步 - 校对
		auth.POST("/user_audits/update", controllers.UserAudits.Update)     // 会员稽核 - 修改
		auth.POST("/user_audits/delete", controllers.UserAudits.Delete)     // 会员稽核 - 删除

		auth.GET("/user_tag_categories", controllers.UserTagCategories.List)          // 会员标签分类 - 默认首页
		auth.GET("/user_tag_categories/create", controllers.UserTagCategories.Create) // 会员标签分类 - 新增页面
		auth.GET("/user_tag_categories/update", controllers.UserTagCategories.Update) // 会员标签分类 - 编辑页面
		auth.POST("/user_tag_categories/save", controllers.UserTagCategories.Save)    // 会员标签分类 - 保存
		auth.GET("/user_tag_categories/delete", controllers.UserTagCategories.Delete) // 会员标签分类 - 删除

		auth.GET("/user_notes/create", controllers.UserNotes.Create) // 会员标签分类 - 添加页面
		auth.GET("/user_notes/update", controllers.UserNotes.Update) // 会员标签分类 - 编辑页面
		auth.POST("/user_notes/save", controllers.UserNotes.Save)    // 会员标签分类 - 保存

		auth.GET("/user_tags", controllers.UserTags.List)          // 会员标签 - 默认首页
		auth.GET("/user_tags/create", controllers.UserTags.Create) // 会员标签 - 新增页面
		auth.GET("/user_tags/update", controllers.UserTags.Update) // 会员标签 - 编辑页面
		auth.POST("/user_tags/save", controllers.UserTags.Save)    // 会员标签 - 保存
		auth.GET("/user_tags/delete", controllers.UserTags.Delete) // 会员标签 - 删除

		// 会员账户调整记录
		auth.GET("/user_changes", controllers.UserChanges.List)          // 默认首页
		auth.GET("/user_changes/agree", controllers.UserChanges.Agree)   // 同意
		auth.GET("/user_changes/refuse", controllers.UserChanges.Refuse) // 拒绝
		auth.POST("/user_changes/save", controllers.UserChanges.Save)    // 拒绝
		// 会员平台转账/user_cards
		auth.GET("/user_transfers", controllers.UserTransfers.List)                                       // 默认首页
		auth.POST("/user_transfers/check_transfer_status", controllers.UserTransfers.CheckTransferStatus) // 转账状态
		// 会员红利记录
		auth.GET("/user_dividends", controllers.UserDevidends.List)     // 会员红利
		auth.GET("/activity_applies", controllers.ActivityApplies.List) // 会员特别

		// 会员活动记录
		auth.GET("/user_activities", controllers.UserActivities.Index)         // 默认首页
		auth.GET("/user_activities/cancel", controllers.UserActivities.Cancel) // 取消活动
		auth.GET("/user_activities/create", controllers.UserActivities.Create) // 活动创建
		auth.POST("/user_activities/save", controllers.UserActivities.Save)    // 活动保存
		auth.GET("/user_activities/agree", controllers.UserActivities.Agree)   // 派发

		// 会员等级记录
		auth.GET("/user_level_changes", controllers.UserLevelChanges.List) // 默认首页

		// 收款银行卡
		auth.GET("/receive_bank_cards", controllers.ReceiveBankCards.List)
		auth.GET("/receive_bank_cards/create", controllers.ReceiveBankCards.Create) // 收款银行卡 - 添加
		auth.GET("/receive_bank_cards/update", controllers.ReceiveBankCards.Update) // 收款银行卡 - 修改
		auth.POST("/receive_bank_cards/save", controllers.ReceiveBankCards.Save)    // 收款银行卡 - 保存
		auth.GET("/receive_bank_cards/delete", controllers.ReceiveBankCards.Delete) // 收款银行卡 - 删除
		auth.GET("/receive_bank_cards/state", controllers.ReceiveBankCards.State)   // 收款银行卡 - 停启用
		// 支付通道
		auth.GET("/payment_channels", controllers.PaymentChannels.List)
		auth.GET("/payment_channels/add", controllers.PaymentChannels.Add)         // 支付通道 - 添加
		auth.GET("/payment_channels/edit", controllers.PaymentChannels.Edit)       // 支付通道 - 修改
		auth.POST("/payment_channels/save_do", controllers.PaymentChannels.SaveDo) // 支付通道 - 保存
		auth.GET("/payment_channels/delete", controllers.PaymentChannels.Delete)   // 支付通道 - 删除
		auth.GET("/payment_channels/state", controllers.PaymentChannels.State)     // 支付通道 - 停启用

		// 离线存款管理
		auth.GET("/user_deposits", controllers.UserDeposits.List)                       // 默认首页-存款中
		auth.POST("/user_deposits/order_info", controllers.UserDeposits.OrderInfo)      // 存款 - 在线信息
		auth.GET("/user_deposits/add_silp", controllers.UserDeposits.AddSlip)           // 添加存款单页面
		auth.GET("/user_deposits/user_info", controllers.UserDeposits.UserInfo)         // 存款单获取用户信息
		auth.POST("/user_deposits/add_slip_save", controllers.UserDeposits.AddSlipSave) // 添加存款单-保存
		auth.GET("/user_deposits/update", controllers.UserDeposits.Update)              // 手动确认页面
		auth.POST("/user_deposits/confirm_do", controllers.UserDeposits.ConfirmDo)      // 手动确认
		auth.POST("/user_deposits/get_status", controllers.UserDeposits.GetStatus)      // 获取状态
		auth.GET("/user_deposits/export", controllers.UserDeposits.Export)              // 获取状态
		auth.POST("/user_deposits/export", controllers.UserDeposits.Export)             // 获取状态
		// 离线存款管理-审核列表
		auth.GET("/user_deposit_audits", controllers.UserDepositAudits.List)           // 审核列表
		auth.POST("/user_deposit_audits/agree", controllers.UserDepositAudits.Agree)   // 同意
		auth.POST("/user_deposit_audits/refuse", controllers.UserDepositAudits.Refuse) // 拒绝
		// 离线存款管理-历史记录
		auth.GET("/user_deposit_hrs", controllers.UserDepositHrs.List)                  // 历史记录
		auth.GET("/user_deposit_hrs/mistake", controllers.UserDepositHrs.Mistake)       // 失误反转
		auth.POST("/user_deposit_hrs/fix", controllers.UserDepositHrs.Fix)              // 手动补单
		auth.POST("/user_deposit_hrs/mistake_do", controllers.UserDepositHrs.MistakeDo) // 失误反转保存
		auth.GET("/user_deposit_logs", controllers.UserDepositLogs.List)                // 日志记录
		auth.GET("/user_deposit_hrs/export", controllers.UserDepositHrs.Export)         // 历史记录
		auth.POST("/user_deposit_hrs/export", controllers.UserDepositHrs.Export)        // 历史记录

		// -- 代币存款
		auth.GET("/user_deposit_coins", controllers.UserDepositCoins.List)                            // 代币存款 - 列表
		auth.GET("/user_deposit_coin_hrs", controllers.UserDepositCoinHrs.List)                       // 代币存款 - 历史记录
		auth.GET("/user_deposit_coin_matches", controllers.UserDepositCoinMatches.List)               // 代币存款匹配 - 列表
		auth.GET("/user_deposit_coin_matches/create", controllers.UserDepositCoinMatches.Create)      // 代币存款匹配 - 添加
		auth.POST("/user_deposit_coin_matches/save", controllers.UserDepositCoinMatches.SaveRecharge) // 代币存款匹配 - 保存

		// 在线存款管理
		auth.GET("/user_deposit_onlines", controllers.UserDepositOnlines.List)                       // 默认首页-存款中
		auth.POST("/user_deposit_onlines/order_info", controllers.UserDepositOnlines.OrderInfo)      // 存款 - 在线信息
		auth.GET("/user_deposit_onlines/add_silp", controllers.UserDepositOnlines.AddSlip)           // 添加存款单页面
		auth.GET("/user_deposit_onlines/user_info", controllers.UserDepositOnlines.UserInfo)         // 存款单获取用户信息
		auth.POST("/user_deposit_onlines/add_slip_save", controllers.UserDepositOnlines.AddSlipSave) // 添加存款单-保存
		auth.GET("/user_deposit_onlines/update", controllers.UserDepositOnlines.Update)              // 手动确认页面
		auth.POST("/user_deposit_onlines/confirm_do", controllers.UserDepositOnlines.ConfirmDo)      // 手动确认
		auth.POST("/user_deposit_onlines/get_status", controllers.UserDepositOnlines.GetStatus)      // 获取状态
		auth.GET("/user_deposit_onlines/export", controllers.UserDepositOnlines.Export)              // 获取状态
		auth.POST("/user_deposit_onlines/export", controllers.UserDepositOnlines.Export)             // 获取状态
		// 在线存款管理-审核列表
		auth.GET("/user_deposit_online_audits", controllers.UserDepositOnlineAudits.List)           // 审核列表
		auth.POST("/user_deposit_online_audits/agree", controllers.UserDepositOnlineAudits.Agree)   // 同意
		auth.POST("/user_deposit_Online_audits/refuse", controllers.UserDepositOnlineAudits.Refuse) // 拒绝
		// 在线存款管理-历史记录
		auth.GET("/user_deposit_online_hrs", controllers.UserDepositOnlineHrs.List)                  // 历史记录
		auth.GET("/user_deposit_online_hrs/mistake", controllers.UserDepositOnlineHrs.Mistake)       // 失误反转
		auth.POST("/user_deposit_online_hrs/fix", controllers.UserDepositOnlineHrs.Fix)              // 手动补单
		auth.POST("/user_deposit_online_hrs/mistake_do", controllers.UserDepositOnlineHrs.MistakeDo) // 失误反转保存
		auth.GET("/user_deposit_online_logs", controllers.UserDepositOnlineLogs.List)                // 日志记录
		auth.GET("/user_deposit_online_hrs/export", controllers.UserDepositOnlineHrs.Export)         // 历史记录
		auth.POST("/user_deposit_online_hrs/export", controllers.UserDepositOnlineHrs.Export)        // 历史记录

		// 提款管理
		auth.GET("/user_withdraws", controllers.UserWithdraws.List)                    // 默认首页-提款中
		auth.GET("/user_withdraws/success", controllers.UserWithdraws.Success)         // 成功
		auth.GET("/user_withdraws/failure", controllers.UserWithdraws.Failure)         // 失败
		auth.POST("/user_withdraws/success_do", controllers.UserWithdraws.SuccessDo)   // 成功保存
		auth.POST("/user_withdraws/failure_do", controllers.UserWithdraws.FailureDo)   // 失败保存
		auth.POST("/user_withdraws/get_status", controllers.UserWithdraws.GetStatus)   // 获取状态
		auth.POST("/user_withdraws/save_config", controllers.UserWithdraws.SaveConfig) // 成功
		auth.GET("/user_withdraw_logs", controllers.UserWithdrawLogs.List)             // 日志记录
		// 提款管理-历史记录
		auth.GET("/user_withdraw_hrs", controllers.UserWithdrawHrs.List)           // 历史记录
		auth.GET("/user_withdraw_hrs/export", controllers.UserWithdrawHrs.Export)  // 历史记录 - 导出
		auth.POST("/user_withdraw_hrs/export", controllers.UserWithdrawHrs.Export) // 历史记录 - 导出
		// 代币提款
		auth.GET("/user_withdraw_coins", controllers.UserWithdrawCoins.List)      // 代币提款 - 处理中
		auth.GET("/user_withdraw_coin_hrs", controllers.UserWithdrawCoinHrs.List) // 代币提款 - 历史记录

		// 稽核
		auth.GET("/audit", controllers.Audits.Deposit)               // 存款稽核
		auth.GET("/audit_deposit", controllers.Audits.Deposit)       // 存款稽核
		auth.GET("/audit_dividend", controllers.Audits.Dividend)     // 优惠稽核
		auth.GET("/audit_deposit/create", controllers.Audits.Create) // 存款稽核调整
		auth.POST("/audit_deposit/save", controllers.Audits.Save)    //
		auth.GET("/audit_logs", controllers.AuditLogs.List)          //

		// 代理提款-提款中
		auth.GET("/agent_withdraws", controllers.AgentWithdraws.Index)                 // 默认首页
		auth.GET("/agent_withdraws/success", controllers.AgentWithdraws.Success)       // 成功
		auth.GET("/agent_withdraws/failure", controllers.AgentWithdraws.Failure)       // 失败
		auth.POST("/agent_withdraws/success_do", controllers.AgentWithdraws.SuccessDo) // 成功保存
		auth.POST("/agent_withdraws/failure_do", controllers.AgentWithdraws.FailureDo) // 失败保存
		auth.GET("/agent_withdraw_logs", controllers.AgentWithdrawLogs.List)           // 日志
		// 代理提款-历史记录
		auth.GET("/agent_withdraw_hrs", controllers.AgentWithdrawHrs.List)           // 历史记录
		auth.GET("/agent_withdraw_hrs/export", controllers.AgentWithdrawHrs.Export)  // 代理提款历史记录 - 导出
		auth.POST("/agent_withdraw_hrs/export", controllers.AgentWithdrawHrs.Export) // 代理提款历史记录 - 导出

		// 手动上下分
		auth.GET("/user_account_sets", controllers.UserAccountSets.List)                 // 默认首页
		auth.GET("/user_account_sets/top_money", controllers.UserAccountSets.TopMoney)   // 上分界面
		auth.GET("/user_account_sets/down_money", controllers.UserAccountSets.DownMoney) // 下分界面
		auth.POST("/user_account_sets/save", controllers.UserAccountSets.Save)           // 上分下分保存
		auth.GET("/user_account_sets/export", controllers.UserAccountSets.Export)        // 导出
		auth.POST("/user_account_sets/export", controllers.UserAccountSets.Export)       // 导出
		// 手动上下分-审核列表
		auth.GET("/user_account_audits", controllers.UserAccountAudits.List)            // 审核列表
		auth.GET("/user_account_audits/agree", controllers.UserAccountAudits.Agree)     // 同意
		auth.GET("/user_account_audits/refuse", controllers.UserAccountAudits.Refuse)   // 拒绝
		auth.POST("/user_account_audits/save_do", controllers.UserAccountAudits.SaveDo) // 同意拒绝保存
		// 手动上下分-历史记录
		auth.GET("/user_account_hrs", controllers.UserAccountHrs.List) // 历史记录

		// 返水管理
		auth.GET("/commission_levels", controllers.CommissionLevls.List)            // 默认首页-返水等级
		auth.GET("/commission_levels/setup", controllers.CommissionLevls.Setup)     // 设置
		auth.POST("/commission_levels/save_do", controllers.CommissionLevls.SaveDo) // 同意拒绝保存
		auth.GET("/commission_levels/details", controllers.CommissionLevls.Details) // 详情
		auth.GET("/user_commissions/manual", controllers.UserCommissions.Manual)    // 返水管理-手动返水
		auth.POST("/user_commissions/issue", controllers.UserCommissions.Issue)     // 手动返水-发放
		auth.GET("/user_commissions/record", controllers.UserCommissions.Record)    // 返水管理-返水记录

		// 红利管理
		auth.GET("/dividend_managements", controllers.DividendManagements.Index)              // 默认首页-红利发放
		auth.GET("/dividend_managements/excel", controllers.DividendManagements.FileDownload) // 红利模板下载
		auth.POST("/dividend_managements/submit_do", controllers.DividendManagements.SubmitDo)
		auth.GET("/dividend_audits", controllers.DividendAudits.List)                      // 审核列表
		auth.GET("/dividend_audits/edit_view", controllers.DividendAudits.EditView)        // 操作页面
		auth.POST("/dividend_audits/agree", controllers.DividendAudits.Agree)              // 同意
		auth.POST("/dividend_audits/refuse", controllers.DividendAudits.Refuse)            // 拒绝
		auth.POST("/dividend_audits/batch_agree", controllers.DividendAudits.BatchAgree)   // 批量同意
		auth.POST("/dividend_audits/batch_refuse", controllers.DividendAudits.BatchRefuse) // 批量拒绝
		auth.GET("/dividend_hrs", controllers.DividendHrs.List)                            // 历史记录

		// 活动管理
		auth.GET("/activities", controllers.Activities.List)                           // 默认首页-常规活动
		auth.GET("/activities/add", controllers.Activities.Add)                        // 添加
		auth.GET("/activities/edit", controllers.Activities.Edit)                      // 编辑
		auth.POST("/activities/save_do", controllers.Activities.SaveDo)                // 保存
		auth.GET("/activities/delete", controllers.Activities.Delete)                  // 赞助配置 - 删除
		auth.GET("/activities/state", controllers.Activities.State)                    // 赞助配置 - 停启用
		auth.GET("/activities/turntable", controllers.Activities.Turntable)            // 转盘活动配置
		auth.POST("/activities/turntable_do", controllers.Activities.TurntableDo)      // 转盘活动配置
		auth.GET("/activities/arrived", controllers.Arrived.Update)                    // 转盘活动配置
		auth.POST("/activities/arrived_do", controllers.Arrived.Save)                  // 转盘活动配置
		auth.GET("/activities/red_envelope_rain", controllers.RedEnvelopes.List)       // 转盘活动配置
		auth.POST("/activities/red_envelope_rain_save", controllers.RedEnvelopes.Save) // 转盘活动配置

		// 活动管理-邀请好友
		auth.GET("/user_invites", controllers.UserInvites.List)
		auth.GET("/user_invites/rule_setting", controllers.UserInvites.RuleSetting) // 规则设置
		auth.POST("/user_invites/save_do", controllers.UserInvites.SaveDo)          // 规则设置-保存
		auth.POST("/user_invites/enable", controllers.UserInvites.Enable)           // 规则设置-开启活动
		auth.GET("/user_invites/agree", controllers.UserInvites.Agree)              // 同意页面
		auth.GET("/user_invites/refuse", controllers.UserInvites.Refuse)            // 拒绝页面
		auth.POST("/user_invites/agree_do", controllers.UserInvites.AgreeDo)        // 同意
		auth.POST("/user_invites/refuse_do", controllers.UserInvites.RefuseDo)      // 拒绝
		// 活动礼金-申请
		auth.GET("/activities_managements", controllers.ActivitiesManagements.Index)              // 默认首页-红利发放
		auth.GET("/activities_managements/excel", controllers.ActivitiesManagements.FileDownload) // 红利模板下载
		auth.POST("/activities_managements/submit_do", controllers.ActivitiesManagements.SubmitDo)
		// 活动礼金-审核
		auth.GET("/activities_audits", controllers.ActivitiesAudits.List)                      // 审核列表
		auth.GET("/activities_audits/edit_view", controllers.ActivitiesAudits.EditView)        // 操作页面
		auth.POST("/activities_audits/agree", controllers.ActivitiesAudits.Agree)              // 同意
		auth.POST("/activities_audits/refuse", controllers.ActivitiesAudits.Refuse)            // 拒绝
		auth.POST("/activities_audits/batch_agree", controllers.ActivitiesAudits.BatchAgree)   // 批量同意
		auth.POST("/activities_audits/batch_refuse", controllers.ActivitiesAudits.BatchRefuse) // 批量拒绝
		// 活动礼金-记录
		auth.GET("/activities_hrs", controllers.ActivitiesHrs.List) // 历史记录

		// 报表管理-经营报表
		auth.GET("/report_operations", controllers.ReportOperations.List)      // 报表管理-经营报表
		auth.GET("/report_games", controllers.ReportGames.List)                // 报表管理-游戏报表
		auth.GET("/report_agents", controllers.ReportAgents.List)              // 报表管理-代理报表
		auth.GET("/report_users", controllers.ReportUsers.List)                // 报表管理-会员报表
		auth.GET("/report_users/detail", controllers.ReportUsers.Detail)       // 报表管理-会员报表
		auth.GET("/game_venues/wallet", controllers.ReportUsers.State)         // 场馆钱包状态,先放着
		auth.GET("/commission_team", controllers.CommissionTeam.List)          // 报表管理-团队报表
		auth.GET("/commission_team_detail", controllers.CommissionTeam.Detail) // 报表管理-团队详情
		// 场馆管理-场馆列表
		auth.GET("/game_venues", controllers.GameVenues.List)                  // 场馆管理-场馆列表
		auth.GET("/game_venues/create", controllers.GameVenues.Create)         // 新增页面
		auth.GET("/game_venues/state", controllers.GameVenues.Stated)          // 场馆状态
		auth.GET("/game_venues/lower_state", controllers.GameVenues.State)     // 场馆状态
		auth.POST("/game_venues/state_save", controllers.GameVenues.StateSave) // 新增和编辑保存
		auth.GET("/game_venues/add", controllers.GameVenues.Add)               // 场馆状态
		auth.GET("/game_venues/games", controllers.GameVenues.Games)           // 场馆状态
		auth.GET("/game_venues/update", controllers.GameVenues.Update)         // 编辑页面
		auth.POST("/game_venues/save", controllers.GameVenues.Saved)           // 新增和编辑保存
		auth.POST("/game_venues/create", controllers.GameVenues.Save)          // 新增和编辑保存

		// 场馆管理-游戏列表
		auth.GET("/games", controllers.Games.List) // 场馆管理-游戏列表a
		// auth.GET("/games/create", controllers.Games.Create) //新增页面
		auth.GET("/games/update", controllers.Games.Update) // 编辑页面
		auth.POST("/games/save", controllers.Games.Save)    // 新增和编辑保存
		auth.GET("/games/state", controllers.Games.State)   // 状态

		// 场馆管理-维护设置
		auth.GET("/game_maintains", controllers.GameMaintains.List)          // 场馆管理-维护设置
		auth.GET("/game_maintains/create", controllers.GameMaintains.Create) // 新增页面
		auth.GET("/game_maintains/update", controllers.GameMaintains.Update) // 维护设置
		auth.POST("/game_maintains/save", controllers.GameMaintains.Save)    // 新增和编辑保存
		auth.GET("/game_maintains_log", controllers.GameMaintainLogs.List)   // 操作日志

		// 版本控制
		auth.GET("/versions", controllers.Versions.List)          // 版本列表
		auth.POST("/version/save", controllers.Versions.Save)     // 版本保存
		auth.GET("/version/updated", controllers.Versions.Update) // 版本更新
		auth.GET("/version/created", controllers.Versions.Create) // 版本创建
		auth.GET("/version/delete", controllers.Versions.Delete)  // 版本创建
		// 站点设置 - 维护设置
		auth.GET("/site_maintains", controllers.SiteMaintains.List)          // 站点设置 - 维护设置
		auth.GET("/site_maintains/update", controllers.SiteMaintains.Update) // 维护设置
		auth.POST("/site_maintains/save", controllers.SiteMaintains.Save)    // 新增和编辑保存
		auth.GET("/site_maintains_log", controllers.SiteMaintainLogs.List)   // 操作日志

		// 参数分组
		auth.GET("/parameter_groups", controllers.ParameterGroups.List)           // 参数分组
		auth.GET("/parameter_groups/created", controllers.ParameterGroups.Create) // 分组
		auth.POST("/parameter_groups/save", controllers.ParameterGroups.Save)     // 分组

		// 参数分组
		auth.GET("/payment_groups", controllers.PaymentGroups.List)          // 参数分组
		auth.GET("/payment_groups/create", controllers.PaymentGroups.Create) // 支付分组 - 添加
		auth.GET("/payment_groups/update", controllers.PaymentGroups.Update) // 支付分组 - 修改
		auth.POST("/payment_groups/save", controllers.PaymentGroups.Save)    // 支付分组 - 保存
		auth.GET("/payment_groups/delete", controllers.PaymentGroups.Delete) // 支付分组 - 删除

		// 参数管理
		auth.GET("/parameters", controllers.Parameters.List)          // 参数管理
		auth.GET("/parameters/update", controllers.Parameters.Update) // 参数管理
		auth.POST("/parameters/save", controllers.Parameters.Save)    // 参数管理
		auth.GET("/parameters/create", controllers.Parameters.Create) // 参数管理
		// 素材管理
		auth.GET("/source_materials/created", controllers.SourceMaterials.Create) // 素材管理
		auth.POST("/source_materials/save", controllers.SourceMaterials.Save)
		auth.GET("/source_materials/updated", controllers.SourceMaterials.Update)
		auth.GET("/source_materials/deleted", controllers.SourceMaterials.Delete)
		auth.GET("/source_materials", controllers.SourceMaterials.List)
		auth.GET("/source_materials/state", controllers.SourceMaterials.State)
		// 站点管理-存款限制
		auth.GET("/deposit_limits", controllers.DepositLimits.Index)                   // 站点管理-存款限制
		auth.GET("/deposit_limits/remind", controllers.DepositLimits.Remind)           // 时间提醒设置
		auth.POST("/deposit_limits/remind_save", controllers.DepositLimits.RemindSave) // 时间提醒设置 - 保存
		auth.POST("/deposit_limits/allow", controllers.DepositLimits.Allow)            // 是否允许多笔未支付订单

		// 站点管理-底部信息
		auth.GET("/site_bottoms", controllers.SiteBottoms.List)          // 站点管理-底部信息
		auth.GET("/site_bottoms/create", controllers.SiteBottoms.Create) // 新增
		auth.GET("/site_bottoms/update", controllers.SiteBottoms.Update) // 修改
		auth.POST("/site_bottoms/save", controllers.SiteBottoms.Save)    // 编辑
		auth.GET("/site_bottoms/delete", controllers.SiteBottoms.Delete) // 删除
		auth.GET("/site_bottoms/state", controllers.SiteBottoms.State)   // 状态

		// 系统管理 - 系统账号管理
		auth.GET("/admins", controllers.Admins.List)                        // 管理员列表
		auth.GET("/admins/kick", controllers.Admins.Kick)                   // 管理员 - T下线
		auth.GET("/admins/create", controllers.Admins.Create)               // 管理员 - 新增
		auth.GET("/admins/update", controllers.Admins.Update)               // 管理员 - 新增
		auth.POST("/admins/save", controllers.Admins.Save)                  // 管理员 - 保存
		auth.GET("/admins/state", controllers.Admins.State)                 // 管理员 - 状态
		auth.GET("/admins/google_code", controllers.Admins.GoogleCode)      // 管理员 - 状态
		auth.GET("/admins/handover", controllers.Admins.Handover)           // 一键交接
		auth.POST("/admins/handover_save", controllers.Admins.HandoverSave) // 一键交接
		// auth.GET("/admins/delete", controllers.Admins.Delete)               //管理员 - 删除

		// 系统管理-角色管理
		auth.GET("/admin_roles", controllers.AdminRoles.List)               // 系统管理-角色管理
		auth.GET("/admin_roles/create", controllers.AdminRoles.Create)      // 系统管理-角色管理
		auth.GET("/admin_roles/update", controllers.AdminRoles.Update)      // 系统管理-角色管理
		auth.POST("/admin_roles/save", controllers.AdminRoles.Save)         // 系统管理-角色管理
		auth.GET("/admin_roles/configs", controllers.AdminRoles.Configs)    // 系统管理-角色管理
		auth.GET("/admin_roles/detail", controllers.AdminRoles.Detail)      // 系统管理 - 角色管理 - 详情
		auth.GET("/admin_roles/sub_menus", controllers.AdminRoles.SubMenus) // 系统管理 - 角色管理 - 二级菜单
		// auth.GET("/admin_roles/delete", controllers.AdminRoles.Delete)      //系统管理-角色管理

		// 系统管理-菜单管理
		auth.GET("/menus", controllers.Menus.List)          // 系统管理 - 菜单列表
		auth.GET("/menus/create", controllers.Menus.Create) // 系统管理 - 管理添加
		auth.GET("/menus/update", controllers.Menus.Update) // 系统管理 - 管理修改
		auth.POST("/menus/save", controllers.Menus.Save)    // 系统管理 - 管理保存

		auth.GET("/access_logs", controllers.AccessLogs.List)               // 后台日志
		auth.GET("/admin_login_logs", admin_login_logs.AdminLoginLogs.List) // 后台日志

		// 系统管理-访问授权
		auth.GET("/admin_authorizes", controllers.AdminAuthorizes.Index)     // 系统管理-访问授权
		auth.GET("/admin_authorizes/add", controllers.AdminAuthorizes.Add)   // 编辑
		auth.GET("/admin_authorizes/edit", controllers.AdminAuthorizes.Edit) // 编辑

		// 系统管理 - 访问授权
		auth.GET("/permission_ips", controllers.PermissionIps.List)
		auth.GET("/permission_ips/create", controllers.PermissionIps.Create) // 增加
		auth.GET("/permission_ips/update", controllers.PermissionIps.Update) // 修改
		auth.POST("/permission_ips/save", controllers.PermissionIps.Save)    // 保存
		auth.GET("/permission_ips/delete", controllers.PermissionIps.Delete) // 删除
		auth.GET("/permission_ips/state", controllers.PermissionIps.State)   // 删除

		// 平台管理 - 平台列表
		auth.GET("/platforms", controllers.Platforms.List)
		auth.GET("/platforms/create", controllers.Platforms.Create) // 增加
		auth.GET("/platforms/update", controllers.Platforms.Update) // 修改
		auth.POST("/platforms/save", controllers.Platforms.Save)    // 保存
		auth.GET("/platforms/state", controllers.Platforms.State)   // 状态

		// 平台箮理 - 站点列表
		auth.GET("/platform_sites", controllers.PlatformSites.List)
		auth.GET("/platform_sites/create", controllers.PlatformSites.Create)           // 增加
		auth.GET("/platform_sites/update", controllers.PlatformSites.Update)           // 修改
		auth.GET("/platform_sites/state", controllers.PlatformSites.State)             // 状态
		auth.POST("/platform_sites/save", controllers.PlatformSites.Save)              // 保存
		auth.GET("/platform_sites/config", controllers.PlatformSites.Config)           // 保存
		auth.POST("/platform_sites/config_save", controllers.PlatformSites.ConfigSave) // 保存

		// 平台管理 - 站点配置
		auth.GET("/platform_site_configs", controllers.PlatformSiteConfigs.List)
		auth.GET("/platform_site_configs/create", controllers.PlatformSiteConfigs.Create) // 增加
		auth.GET("/platform_site_configs/update", controllers.PlatformSiteConfigs.Update) // 修改
		auth.GET("/platform_site_configs/state", controllers.PlatformSiteConfigs.State)   // 删除
		auth.POST("/platform_site_configs/save", controllers.PlatformSiteConfigs.Save)    // 保存

		// 赛事直播
		auth.GET("/sport_live", controllers.SportLives.List)                               // 直播赛事
		auth.GET("/sport_live/stop", controllers.SportLives.Stop)                          // 禁言
		auth.GET("/sport_live/manage", controllers.SportLives.Manage)                      // 直播管理
		auth.GET("/sport_live/manage/words", controllers.SportLives.ManageWords)           // 过滤词
		auth.POST("/sport_live/manage/words_save", controllers.SportLives.ManageWordsSave) // 过滤词修改

		auth.GET("/admin_tools/down_vips", admin_tools.DownVips) // vip 降级处理
		auth.GET("/admin_tools/down_vip", admin_tools.DownVip)   // vip 降级处理
		auth.GET("/admin_tools/up_vips", admin_tools.UpVips)     // vip 降级处理
		auth.GET("/admin_tools/up_vip", admin_tools.UpVip)       // vip 降级处理

	}
}
