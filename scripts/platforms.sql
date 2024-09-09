-- 一级站点:  平台管理
INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (900008, 0, '平台管理', '#', '1', 'layui-icon-windows');
	-- 二级站点: 平台列表
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (910007, 900008, '平台列表', '/platforms', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (911007, 910007, '平台列表', '/platforms', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (911107, 911007, '系统用户列表', '/platforms', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (911208, 911007, '系统用户添加', '/platforms/create|/platforms/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (911309, 911007, '系统用户修改', '/platforms/update|/platforms/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (911511, 911007, '系统用户状态', '/platforms/state', '4', '');
	-- 二级站点: 站点管理
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (920008, 900008, '网站管理', '/platform_sites', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (921008, 920008, '网站管理', '/platform_sites', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (921108, 921008, '网站管理列表', '/platform_sites', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (921209, 921008, '网站管理添加', '/platform_sites/create|/platform_sites/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (921310, 921008, '网站管理修改', '/platform_sites/update|/platform_sites/save', '4', '');
	-- 二级站点: 站点配置
	INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (930009, 900008, '参数配置', '/platform_site_configs', '2', '');
		INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (931009, 930009, '参数配置', '/platform_site_configs', '3', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (931109, 931009, '参数配置列表', '/platform_site_configs', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (931210, 931009, '参数配置添加', '/platform_site_configs/create|/platform_site_configs/save', '4', '');
			INSERT INTO menus (id, parent_id, name, url, level, icon) VALUES (931311, 931009, '参数配置修改', '/platform_site_configs/update|/platform_site_configs/save', '4', '');
