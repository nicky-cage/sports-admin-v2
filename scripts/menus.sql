-- //
-- * 严重警告: 此文件为只读, 请勿修改此文件
-- * 如果要添加、修改菜单/结构, 请修改menu.json文件, 并将菜单结构反映其中
-- * 修改完成 menu.json 文件之后, 执行 php -f ./create_menus.php > menus.sql 生成新的菜单文件即可
-- //

-- 清空当前菜单记录
truncate menus;


-- 一级菜单: 会员管理
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (100000, 0, '会员管理', '#', '1', 'layui-icon-user');
	-- 二级菜单: 会员列表
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (110000, 100000, '会员列表', '/users', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111000, 110000, '会员列表', '/users', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111100, 111000, '会员列表', '/users', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111201, 111000, '添加标签', '/users/add_tags|/users/add_tags_save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111302, 111000, '批量禁用', '/users/disable_all', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111403, 111000, '批量匹配', '/users', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111504, 111000, '修改资料', '/users/update|/users/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111605, 111000, '用户状态', '/users/state', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111706, 111000, '钱包金额', '/user_detail/account_async', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111807, 111000, '添加用户', '/users/create|/users/add', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (111908, 111000, '修改密码', '/users/password|/users/password_save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112009, 111000, '会员详情', '/users/detail', '4', '');
				INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112019, 112009, '中心钱包', '/users/detail', '5', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112020, 112019, '基本信息', '/users/detail', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112021, 112019, '账户信息', '/user_detail/accounts|/user_detail/accounts/transfer_out', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112022, 112019, '输赢信息', '/user_detail/wins', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112023, 112019, '存款信息', '/user_detail/deposits', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112024, 112019, '取款信息', '/user_detail/withdraws', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112025, 112019, '红利信息', '/user_detail/dividends', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112026, 112019, '返水信息', '/user_detail/commissions', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112027, 112019, '平台转账', '/user_detail/transfers', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112028, 112019, '账户调整', '/user_detail/resets|/user_detail/resets/save', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112029, 112019, '调整记录', '/user_detail/changes', '6', '');
				INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112030, 112009, '佣金钱包', '/user_detail/commission_accounts', '5', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112031, 112030, '账户信息', '/user_detail/commission_accounts', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112032, 112030, '提款信息', '/user_detail/commission_withdraws', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112033, 112030, '佣金信息', '/user_detail/commission_records', '6', '');
					INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (112034, 112030, '账户金额', '/user_detail/account', '6', '');
	-- 二级菜单: 会员等级
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (120001, 100000, '会员等级', '/user_levels', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (121001, 120001, '会员等级', '/user_levels', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (121101, 121001, '会员等级列表', '/user_levels', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (121202, 121001, '会员等级修改', '/user_levels/update|/user_levels/save', '4', '');
	-- 二级菜单: 会员标签
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (130002, 100000, '会员标签', '/user_tags', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (131002, 130002, '会员标签', '/user_tags', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (131102, 131002, '会员标签列表', '/user_tags', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (131203, 131002, '会员标签添加', '/user_tags/create|/user_tags/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (131304, 131002, '会员标签修改', '/user_tags/update|/user_tags/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (131405, 131002, '会员标签删除', '/user_tags/delete', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (132003, 130002, '标签分类', '/user_tag_categories', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (132103, 132003, '标签分类列表', '/user_tag_categories', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (132204, 132003, '标签分类添加', '/user_tag_categories/create|/user_tag_categories/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (132305, 132003, '标签分类修改', '/user_tag_categories/update|/user_tag_categories/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (132406, 132003, '标签分类删除', '/user_tag_categories/delete', '4', '');
	-- 二级菜单: 会员绑卡
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (140003, 100000, '会员绑卡', '/user_cards', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (141003, 140003, '会员绑卡', '/user_cards', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (141103, 141003, '会员银行卡列表', '/user_cards/', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (141204, 141003, '会员银行卡添加', '/user_cards/create|/user_cards/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (141305, 141003, '会员银行卡修改', '/user_cards/update|/user_cards/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (141406, 141003, '会员银行卡删除', '/user_cards/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (141507, 141003, '会员银行卡详细', '/user_cards/detail', '4', '');
	-- 二级菜单: 投注管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (150004, 100000, '投注管理', '/user_bets', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (151004, 150004, '投注管理', '/user_bets', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (151104, 151004, '投注管理', '/user_bets', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (151205, 151004, '手动补单', '/user_bets/set_up|/user_bets/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (151306, 151004, '导出数据', '/user_bets', '4', '');
	-- 二级菜单: 记录管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (160005, 100000, '记录管理', '/user_changes', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (161005, 160005, '账户调整', '/user_changes', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (161105, 161005, '同意', '/user_changes/agree', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (161206, 161005, '拒绝', '/user_changes/refuse|/user_changes/save', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (162006, 160005, '平台转账', '/user_transfers', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (162106, 162006, '平台转账', '/user_transfers|/user_transfers/check_transfer_status', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (162207, 162006, '导出数据', '/user_changes', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (163007, 160005, '红利记录', '/user_dividends', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (163107, 163007, '红利记录', '/user_dividends', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (163208, 163007, '导出数据', '/user_changes', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (164008, 160005, '活动记录', '/user_activities', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (164108, 164008, '活动记录', '/user_activities', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (164209, 164008, '导出数据', '/user_changes', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (165009, 160005, '等级记录', '/user_level_changes', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (165109, 165009, '等级记录', '/user_level_changes', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (165210, 165009, '导出数据', '/user_changes', '4', '');
	-- 二级菜单: 会员验证码查询
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (170006, 100000, '会员验证码查询', '/users_code', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (171006, 170006, '会员验证码查询', '/users_code', '3', '');

-- 一级菜单: 财务管理
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (200001, 0, '财务管理', '#', '1', 'layui-icon-rmb');
	-- 二级菜单: 收款银行卡
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (210001, 200001, '收款银行卡', '/receive_bank_cards', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (211001, 210001, '收款银行卡', '/receive_bank_cards', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (211101, 211001, '添加', '/receive_bank_cards/create|/receive_bank_cards/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (211202, 211001, '修改', '/receive_bank_cards/update|/receive_bank_cards/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (211303, 211001, '删除', '/receive_bank_cards/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (211404, 211001, '状态', '/receive_bank_cards/state', '4', '');
	-- 二级菜单: 支付通道
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (220002, 200001, '支付通道', '/payment_channels', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (221002, 220002, '支付通道', '/payment_channels', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (221102, 221002, '添加', '/payment_channels/add|/payment_channels/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (221203, 221002, '修改', '/payment_channels/edit|/payment_channels/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (221304, 221002, '删除', '/payment_channels/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (221405, 221002, '状态', '/payment_channels/state', '4', '');
	-- 二级菜单: 存款管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (230003, 200001, '存款管理', '/user_deposits', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (231003, 230003, '存款中', '/user_deposits', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (231103, 231003, '存款列表', '/user_deposits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (231204, 231003, '手动确认', '/user_deposits/update|/user_deposits/confirm_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (231305, 231003, '获取状态', '/user_deposits/get_status', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (231406, 231003, '导出数据', '/user_deposits', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (232004, 230003, '历史记录', '/user_deposit_hrs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (232104, 232004, '历史记录', '/user_deposit_hrs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (232205, 232004, '日志', '/user_deposit_logs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (232306, 232004, '导出数据', '/user_deposit_hrs', '4', '');
	-- 二级菜单: 提款管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (240004, 200001, '提款管理', '/user_withdraws', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (241004, 240004, '提款中', '/user_withdraws', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (241104, 241004, '提款列表', '/user_withdraws', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (241205, 241004, '成功', '/user_withdraws/success|/user_withdraws/success_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (241306, 241004, '失败', '/user_withdraws/failure|/user_withdraws/failure_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (241407, 241004, '获取状态', '/user_withdraws/get_status', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (241508, 241004, '日志', '/user_withdraw_logs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (241609, 241004, '导出数据', '/user_withdraws', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (242005, 240004, '历史记录', '/user_withdraw_hrs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (242105, 242005, '历史记录', '/user_withdraw_hrs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (242206, 242005, '日志', '/user_withdraw_logs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (242307, 242005, '导出数据', '/user_withdraw_hrs', '4', '');
	-- 二级菜单: 代理提款
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (250005, 200001, '代理提款', '/agent_withdraws', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (251005, 250005, '提款中', '/agent_withdraws', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (251105, 251005, '提款列表', '/agent_withdraws', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (251206, 251005, '成功', '/agent_withdraws/success|/agent_withdraws/success_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (251307, 251005, '失败', '/agent_withdraws/failure|/agent_withdraws/failure_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (251408, 251005, '日志', '/agent_withdraw_logs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (251509, 251005, '导出数据', '/agent_withdraws', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (252006, 250005, '历史记录', '/agent_withdraw_hrs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (252106, 252006, '历史记录', '/agent_withdraw_hrs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (252207, 252006, '日志', '/agent_withdraw_logs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (252308, 252006, '导出数据', '/agent_withdraw_hrs', '4', '');
	-- 二级菜单: 手动上分
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (260006, 200001, '手动上分', '/user_account_sets', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (261006, 260006, '手动上下分', '/user_account_sets', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (261106, 261006, '手动上下分', '/user_account_sets', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (261207, 261006, '手动上分', '/user_account_sets/top_money|/user_account_sets/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (261308, 261006, '手动下分', '/user_account_sets/down_money|/user_account_sets/save', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (262007, 260006, '审核列表', '/user_account_audits', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (262107, 262007, '审核列表', '/user_account_audits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (262208, 262007, '同意', '/user_account_audits/agree|/user_account_audits/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (262309, 262007, '拒绝', '/user_account_audits/refuse|/user_account_audits/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (262410, 262007, '导出数据', '/user_account_audits', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (263008, 260006, '历史记录', '/user_account_hrs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (263108, 263008, '历史记录', '/user_account_hrs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (263209, 263008, '导出数据', '/user_account_hrs', '4', '');

-- 一级菜单: 运营管理
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (300002, 0, '运营管理', '#', '1', 'layui-icon-component');
	-- 二级菜单: 广告管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (310002, 300002, '广告管理', '/ads', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (311002, 310002, 'APP启动页', '/ads', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (311102, 311002, 'APP启动页广告列表', '/ads', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (311203, 311002, 'APP启动页广告添加', '/ads/create|/ads/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (311304, 311002, 'APP启动页广告修改', '/ads/update|/ads/update', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (311405, 311002, 'APP启动页广告删除', '/ads/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (311506, 311002, 'APP启动页广告状态', '/ads/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (312003, 310002, '轮播图', '/ad_carousels', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (312103, 312003, '轮播图列表', '/ad_carousels', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (312204, 312003, '轮播图添加', '/ad_carousels/create|/ad_carousels/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (312305, 312003, '轮播图修改', '/ad_carousels/update|/ad_carousels/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (312406, 312003, '轮播图删除', '/ad_carousels/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (312507, 312003, '轮播图状态', '/ad_carousels/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (313004, 310002, '体育赛事', '/ad_matches', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (313104, 313004, '赛事列表', '/ad_matches', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (313205, 313004, '赛事添加', '/ad_matches/create|/ad_matches/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (313306, 313004, '赛事修改', '/ad_matches/update|/ad_matches/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (313407, 313004, '赛事删除', '/ad_matches/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (313508, 313004, '赛事状态', '/ad_matches/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (314005, 310002, '赞助配置', '/ad_sponsors', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (314105, 314005, '赞助列表', '/ad_sponsors', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (314206, 314005, '赞助添加', '/ad_sponsors/create|/ad_sponsors/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (314307, 314005, '赞助修改', '/ad_sponsors/update|/ad_sponsors/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (314408, 314005, '赞助删除', '/ad_sponsors/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (314509, 314005, '赞助状态', '/ad_sponsors/state', '4', '');
	-- 二级菜单: 内容管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (320003, 300002, '内容管理', '/messages', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (321003, 320003, '系统公告', '/notices', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (321103, 321003, '公告列表', '/notices', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (321204, 321003, '公告添加', '/notices/create|/notices/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (321305, 321003, '公告修改', '/notices/update|/notices/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (321406, 321003, '公告删除', '/notices/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (321507, 321003, '消息状态', '/notices/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (322004, 320003, '站内消息', '/messages', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (322104, 322004, '消息列表', '/messages', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (322205, 322004, '消息添加', '/messages/create|/messages/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (322306, 322004, '消息修改', '/messages/update|/messages/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (322407, 322004, '消息删除', '/messages/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (322508, 322004, '消息置顶', '/messages/top', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (322609, 322004, '消息状态', '/messages/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (323005, 320003, '玩家反馈', '/user_feedback', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (323105, 323005, '用户反馈列表', 'user_feedback', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (323206, 323005, '反馈修改', '/user_feedback/update|/user_feedback/Save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (323307, 323005, '反馈详情', '/user_feedback/delete', '4', '');
	-- 二级菜单: 返水管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (330004, 300002, '返水管理', '/commission_levels', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (331004, 330004, '返水等级', '/commission_levels', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (331104, 331004, '返水等级列表', '/commission_levels', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (331205, 331004, '返水等级修改', '/commission_levels/setup|/commission_levels/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (331306, 331004, '返水等级详情', '/commission_levels/details', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (332005, 330004, '手动返水', '/commission_levels', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (332105, 332005, '手动返水', '/commission_levels', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (332206, 332005, '发放返水', '/commission_levels', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (333006, 330004, '返水记录', '/user_commissions/record', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (333106, 333006, '返水记录', '/user_commissions/record', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (333207, 333006, '导出数据', '/user_commissions/record', '4', '');
	-- 二级菜单: 红利管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (340005, 300002, '红利管理', '/dividend_managements', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (341005, 340005, '发放红利', '/dividend_managements', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (341105, 341005, '发放红利', '/dividend_managements', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (341206, 341005, '下载载模板', '/dividend_managements/excel', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (341307, 341005, '红利提交', '/dividend_managements/submit_do', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (342006, 340005, '审核列表', '/dividend_audits', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (342106, 342006, '审核列表', '/dividend_audits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (342207, 342006, '同意', '/dividend_audits/edit_view|/dividend_audits/agree', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (342308, 342006, '拒绝', '/dividend_audits/edit_view|/dividend_audits/refuse', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (342409, 342006, '批量通过', '/dividend_audits/edit_view|/dividend_audits/batch_agree', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (342510, 342006, '批量拒绝', '/dividend_audits/edit_view|/dividend_audits/batch_refuse', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (343007, 340005, '历史记录', '/dividend_hrs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (343107, 343007, '历史记录', '/dividend_hrs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (343208, 343007, '导出数据', '/dividend_hrs', '4', '');
	-- 二级菜单: 活动管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (350006, 300002, '活动管理', '/activities', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (351006, 350006, '常规活动', '/activities', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (351106, 351006, '活动列表', '/activities', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (351207, 351006, '活动添加', '/activities/add|/activities/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (351308, 351006, '活动修改', '/activities/edit|/activities/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (351409, 351006, '活动删除', '/activities/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (351510, 351006, '活动状态', '/activities/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (352007, 350006, '邀请好友', '/user_invites', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (352107, 352007, '邀请列表', '/user_invites', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (352208, 352007, '规则设置', '/user_invites/rule_setting|/user_invites/save_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (352309, 352007, '开启活动', '/user_invites/enable', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (352410, 352007, '同意', '/user_invites/agree|/user_invites/agree_do', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (352511, 352007, '拒绝', '/user_invites/refuse|/user_invites/refuse_do', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (353008, 350006, '发放活动礼金', '/activities_managements', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (353108, 353008, '发放活动礼金', '/activities_managements', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (353209, 353008, '下载载模板', '/activities_managements/excel', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (353310, 353008, '礼金申请提交', '/activities_managements/submit_do', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (354009, 350006, '审核列表', '/activities_audits', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (354109, 354009, '审核列表', '/activities_audits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (354210, 354009, '同意', '/activities_audits/edit_view|/activities_audits/agree', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (354311, 354009, '拒绝', '/activities_audits/edit_view|/activities_audits/refuse', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (354412, 354009, '批量通过', '/activities_audits/edit_view|/activities_audits/batch_agree', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (354513, 354009, '批量拒绝', '/activities_audits/edit_view|/activities_audits/batch_refuse', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (355010, 350006, '历史记录', '/activities_hrs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (355110, 355010, '历史记录', '/activities_hrs', '4', '');
	-- 二级菜单: 存款优惠
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (360007, 300002, '存款优惠', '/deposit_discounts', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (361007, 360007, '存款优惠', '/deposit_discounts', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (361107, 361007, '存款优惠列表', '/deposit_discounts', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (361208, 361007, '存款优惠状态', '/deposit_discounts/state', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (361309, 361007, '存款优惠修改', '/deposit_discounts/update|/deposit_discounts/save', '4', '');

-- 一级菜单: 风控管理
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (400003, 0, '风控管理', '#', '1', 'layui-icon-auz');
	-- 二级菜单: 提款审核
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (410003, 400003, '提款审核', '/risk_audits', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (411003, 410003, '待审核列表', '/risk_audits', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (411103, 411003, '待审列表', '/risk_audits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (411204, 411003, '系统审核', '/risk_audits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (411305, 411003, '人工审核', '/risk_audits', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (412004, 410003, '审核挂起列表', '/risk_audits/handuplist', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (412104, 412004, '待审列表', '/risk_audits/handuplist', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (412205, 412004, '挂起新增', '/risk_audits/handuplist', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (412306, 412004, '批量导入', '/risk_audits/handuplist', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (413005, 410003, '历史记录列表', '/risk_audits/history', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (413105, 413005, '记录列表', '/risk_audits/history', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (413206, 413005, '记录新增', '/risk_audits/history', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (413307, 413005, '批量导入', '/risk_audits/history', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (413408, 413005, '历史详情', '/risk_audits/history_detail', '4', '');
	-- 二级菜单: 拉黑名单
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (420004, 400003, '拉黑名单', '/black_lists', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (421004, 420004, '设备编号', '/black_lists', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (421104, 421004, '设备编号', '/blocked_devices', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (421205, 421004, '设备添加', '/blocked_devices/create|/blocked_devices/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (421306, 421004, '设备修改', '/blocked_devices/update|/blocked_devices/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (421407, 421004, '设备删除', '/blocked_devices/delete', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (422005, 420004, 'IP地址', '/blocked_ips', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (422105, 422005, 'IP地址列表', '/blocked_ips', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (422206, 422005, 'IP地址添加', '/blocked_ips/create|/blocked_ips/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (422307, 422005, 'IP地址修改', '/blocked_ips/update|/blocked_ips/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (422408, 422005, 'IP地址删除', '/blocked_ips/delete', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (423006, 420004, '邮箱地址', '/blocked_mails', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (423106, 423006, '邮箱列表', '/blocked_mails', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (423207, 423006, '邮箱添加', '/blocked_mails/create|/blocked_mails/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (423308, 423006, '邮箱修改', '/blocked_mails/update|/blocked_mails/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (423409, 423006, '邮箱删除', '/blocked_mails/delete', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (424007, 420004, '手机号码', '/blocked_phones', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (424107, 424007, '手机号码列表', '/blocked_phones', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (424208, 424007, '手机号码添加', '/blocked_phones/create|/blocked_phones/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (424309, 424007, '手机号码修改', '/blocked_phones/update|/blocked_phones/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (424410, 424007, '手机号码删除', '/blocked_phones/delete', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (425008, 420004, '银行卡号', '/blocked_cards', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (425108, 425008, '银行卡号列表', '/blocked_cards', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (425209, 425008, '银行卡号添加', '/blocked_cards/create|/blocked_cards/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (425310, 425008, '银行卡号修改', '/blocked_cards/update|/blocked_cards/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (425411, 425008, '银行卡号删除', '/blocked_cards/delete', '4', '');
	-- 二级菜单: 登录日志
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (430005, 400003, '登录日志', '/user_login_logs', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (431005, 430005, '登录日志', '/user_login_logs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (431105, 431005, '登录日志列表', '/user_login_logs', '4', '');

-- 一级菜单: 报表中心
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (500004, 0, '报表中心', '#', '1', 'layui-icon-table');
	-- 二级菜单: 经营报表
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (510004, 500004, '经营报表', '/report_operations', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (511004, 510004, '经营报表', '/report_operations', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (511104, 511004, '经营报表', '/report_operations', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (511205, 511004, '导出数据', '/report_operations', '4', '');
	-- 二级菜单: 游戏报表
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (520005, 500004, '游戏报表', '/report_games', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (521005, 520005, '游戏报表', '/report_games', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (521105, 521005, '游戏报表', '/report_games', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (521206, 521005, '导出数据', '/report_games', '4', '');
	-- 二级菜单: 代理报表
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (530006, 500004, '代理报表', '/report_agents', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (531006, 530006, '游戏报表', '/report_agents', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (531106, 531006, '游戏报表', '/report_agents', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (531207, 531006, '导出数据', '/report_agents', '4', '');
	-- 二级菜单: 会员报表
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (540007, 500004, '会员报表', '/report_users', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (541007, 540007, '会员报表', '/report_users', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (541107, 541007, '会员报表', '/report_users', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (541208, 541007, '导出数据', '/report_users', '4', '');

-- 一级菜单: 站点管理
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (600005, 0, '站点管理', '#', '1', 'layui-icon-util');
	-- 二级菜单: 站点设置
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (610005, 600005, '站点设置', '/configs/update?id=1', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (611005, 610005, '基本信息', '/configs/update?id=1', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (611105, 611005, '站点设置', '/configs/update?id=1', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (611206, 611005, '站点设置', '/configs/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (611307, 611005, '站点设置', '/configs/update?id=1', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612006, 610005, '帮助中心', '/helps', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612106, 612006, '帮助中心', '/helps', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612207, 612006, '帮助分类添加', '/help_categories/create|/help_categories/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612308, 612006, '帮助分类修改', '/help_categories/update|/help_categories/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612409, 612006, '帮助分类删除', '/help_categories/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612510, 612006, '帮助列表', '/helps', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612611, 612006, '帮助添加', '/helps/create|/helps/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612712, 612006, '帮助修改', '/helps/update|/helps/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612813, 612006, '帮助删除', '/helps/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (612914, 612006, '帮助预览', '/helps/detail', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (613015, 612006, '帮助状态', '/helps/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (613007, 610005, '维护设置', '/site_maintains', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (613107, 613007, '维护设置', '/site_maintains/update', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (613208, 613007, '维护修改', '/site_maintains/update|/site_maintains/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (613309, 613007, '维护日志', '/site_maintains_log', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (614008, 610005, '通用配置', '/parameter_groups', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (614108, 614008, '通用配置', '/parameter_groups', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (614209, 614008, '参数设置', '/parameters', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (614310, 614008, '通用配置创建', '/parameter_groups/created|/parameter_groups/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (614411, 614008, '参数更改', '/parameters/update|/parameters/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (614512, 614008, '参数更改', '/parameters/create|/parameters/save', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (615009, 610005, '存款限制', '/deposit_limits', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (615109, 615009, '存款限制', '/deposit_limits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (615210, 615009, '时间提醒', '/deposit_limits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (615311, 615009, '订单设置', '/deposit_limits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (615412, 615009, '站点设置', '/deposit_limits', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (616010, 610005, 'PC底部信息', '/site_bottoms', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (616110, 616010, '底部信息列表', '/site_bottoms', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (616211, 616010, '底部信息添加', '/site_bottoms/creat', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (616312, 616010, '底部信息修改', '/site_bottoms/update|/site_bottoms/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (616413, 616010, '底部信息删除', '/site_bottoms/delete|/site_bottoms/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (616514, 616010, '底部信息添加', '/site_bottoms/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (617011, 610005, '银行列表', '/banks', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (617111, 617011, '银行列表', '/banks', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (617212, 617011, '银行添加', '/banks/create|/banks/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (617313, 617011, '银行修改', '/banks/update|/banks/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (617414, 617011, '银行删除', '/banks/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (617515, 617011, '银行状态', '/banks/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (618012, 610005, '版本列表', '/versions', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (618112, 618012, '版本列表', '/versions', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (618213, 618012, '版本添加', '/versions/created|/versions/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (618314, 618012, '版本修改', '/versions/updated|/versions/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (618415, 618012, '版本删除', '/versions/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (618516, 618012, '版本状态', '/versions/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (619013, 610005, '常用配置', '/popular', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (619113, 619013, '常用配置', '/popular', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (619214, 619013, '常用配置保存', '/popular/save', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (620014, 610005, '推广素材', '/source_materials', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (620114, 620014, '推广素材', '/source_materials', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (620215, 620014, '推广素材创建', '/source_materials/created|/source_materials/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (620316, 620014, '推广素材修改', '/source_materials/updated|/source_materials/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (620417, 620014, '推广素材删除', '/source_materials/deleted', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (620518, 620014, '推广素材状态', '/source_materials/state', '4', '');
	-- 二级菜单: 场馆管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (620006, 600005, '场馆管理', '/game_venues', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (621006, 620006, '场馆列表', '/game_venues', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (621106, 621006, '场馆列表', '/game_venues', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (621207, 621006, '场馆状态', '/game_venues/create|/game_venues/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (621308, 621006, '场馆修改', '/game_venues/update|/game_venues/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (621409, 621006, '场馆下级状态', '/game_venues/lower_state', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (621510, 621006, '场馆下级状态保存', '/game_venues/state_save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (621611, 621006, '场馆钱包状态', '/game_venues/state', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (621712, 621006, '场馆创建返回视图', '/game_venues/add', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (622007, 620006, '电子游戏列表', '/games', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (622107, 622007, '游戏列表', '/games', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (622208, 622007, '游戏修改', '/games/update|/games/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (622309, 622007, '游戏状态', '/games/state', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (623008, 620006, '场馆游戏', '/game_venues/games', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (623108, 623008, '场馆游戏', '/game_venues/games', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (623209, 623008, '场馆操作日志', '/game_maintains_log', '4', '');

-- 一级菜单: 代理管理
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (700006, 0, '代理管理', '#', '1', 'layui-icon-group');
	-- 二级菜单: 代理列表
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (710006, 700006, '代理列表', '/agents', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (711006, 710006, '代理列表', '/agents', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (711106, 711006, '代理列表', '/agents', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (711207, 711006, '代理查看', '/agents/detail_view', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (711308, 711006, '代理修改', '/agents/update|/agents/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (711409, 711006, '增加会员', '/agents/add|/agents/lower_add', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (711510, 711006, '代理新增', '/agents/insert', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (711611, 711006, '代理详细', '/agents/detail', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (711712, 711006, '代理新增视图', '/agents/create', '4', '');
	-- 二级菜单: 下线会员
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (720007, 700006, '下线会员', '/agents/users', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (721007, 720007, '下线会员', '/agents/users', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (721107, 721007, '下线会员', '/agents/users', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (721208, 721007, '下线查看', '/agents/users/detail', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (721309, 721007, '转代理线', '/agents/users/update|/agents/users_save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (721410, 721007, '查看是否代理', '/agents/user_check', '4', '');
	-- 二级菜单: 佣金管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (730008, 700006, '佣金管理', '/agents/commissions', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (731008, 730008, '发放佣金', '/agents/commissions', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (731108, 731008, '发放佣金', '/agents/commissions/grant', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (731209, 731008, '佣金调整', '/agents/commissions/adjustment|/agents/commissions/save', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (732009, 730008, '佣金记录', '/agents/commissions/record', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (732109, 732009, '佣金记录', '/agents/commissions/record', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (732210, 732009, '导出数据', '/agents/commissions/record', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733010, 730008, '返佣方案', '/agents/commissions/plan', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733110, 733010, '返佣方案', '/agents/commissions/plan', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733211, 733010, '佣金方案更改', '/agents/commissions/plan/updated', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733312, 733010, '佣金方案', '/agents/commissions/plan', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733413, 733010, '佣金方案新增', '/agents/commissions/plan/add', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733514, 733010, '佣金方案保存', '/agents/commissions/plan/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733615, 733010, '佣金方案删除', '/agents/commissions/plan/deleted', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733716, 733010, '佣金方案视图', '/agents/commissions/plan/view', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733817, 733010, '佣金方案展示', '/agents/commissions/plan_list', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (733918, 733010, '佣金方案增加', '/agents/commissions/create', '4', '');
	-- 二级菜单: 提款管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (740009, 700006, '提款管理', '/agents/withdraws', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (741009, 740009, '提款审核', '/agents/withdraws', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (741109, 741009, '提款审核', '/agents/withdraws', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (741210, 741009, '批量通过', '/agents/withdraws', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (741311, 741009, '批量拒绝', '/agents/withdraws', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (742010, 740009, '提款记录', '/agents/withdraws/record', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (742110, 742010, '提款记录', '/agents/withdraws/record', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (742211, 742010, '导出数据', '/agents/withdraws/record', '4', '');
	-- 二级菜单: 记录查询
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (750010, 700006, '记录查询', '/agents/logs', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (751010, 750010, '红利', '/agents/logs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (751110, 751010, '红利列表', '/agents/logs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (751211, 751010, '导出数据', '/agents/logs', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (752011, 750010, '返水', '/agents/logs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (752111, 752011, '返水列表', '/agents/logs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (752212, 752011, '导出数据', '/agents/logs', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (753012, 750010, '游戏', '/agents/logs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (753112, 753012, '游戏列表', '/agents/logs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (753213, 753012, '导出数据', '/agents/logs', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (754013, 750010, '输赢调整', '/agents/logs_adjust', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (754113, 754013, '调整列表', '/agents/logs_adjust', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (754214, 754013, '导出数据', '/agents/logs_adjust', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (755014, 750010, '转代', '/agents/logs_transfer', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (755114, 755014, '转代列表', '/agents/logs_transfer', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (755215, 755014, '导出数据', '/agents/logs_transfer', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (756015, 750010, '存款', '/agents/logs_deposits', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (756115, 756015, '存款列表', '/agents/logs_deposits', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (756216, 756015, '导出数据', '/agents/logs_deposits', '4', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (757016, 750010, '登录', '/agents/logs_login', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (757116, 757016, '登录列表', '/agents/logs_login', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (757217, 757016, '导出数据', '/agents/logs_login', '4', '');

-- 一级菜单: 系统管理
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (800007, 0, '系统管理', '#', '1', 'layui-icon-windows');
	-- 二级菜单: 系统账号
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (810007, 800007, '系统账号', '/admins', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (811007, 810007, '系统账号', '/admins', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (811107, 811007, '系统用户列表', '/admins', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (811208, 811007, '系统用户添加', '/admins/create|/admins/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (811309, 811007, '系统用户修改', '/admins/update|/admins/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (811410, 811007, '系统用户删除', '/admins/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (811511, 811007, '系统用户状态', '/admins/state', '4', '');
	-- 二级菜单: 菜单管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (820008, 800007, '菜单管理', '/menus', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (821008, 820008, '菜单管理', '/menus', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (821108, 821008, '菜单管理列表', '/menus', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (821209, 821008, '菜单管理添加', '/menus/create|/menus/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (821310, 821008, '菜单管理修改', '/menus/update|/menus/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (821411, 821008, '菜单管理删除', '/menus/delete', '4', '');
	-- 二级菜单: 角色管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (830009, 800007, '角色管理', '/admin_roles', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (831009, 830009, '角色管理', '/admin_roles', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (831109, 831009, '角色管理列表', '/admin_roles', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (831210, 831009, '角色管理添加', '/admin_roles/create|/admin_roles/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (831311, 831009, '角色管理修改', '/admin_roles/update|/admin_roles/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (831412, 831009, '角色管理删除', '/admin_roles/delete', '4', '');
	-- 二级菜单: 访问权限
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (840010, 800007, '访问权限', '/permission_ips', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (841010, 840010, '访问权限', '/permission_ips', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (841110, 841010, '访问权限列表', '/permission_ips', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (841211, 841010, '访问权限添加', '/permission_ips/create|/permission_ips/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (841312, 841010, '访问权限修改', '/permission_ips/update|/permission_ips/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (841413, 841010, '访问权限删除', '/permission_ips/delete', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (841514, 841010, '访问权限状态', '/permission_ips/state', '4', '');
