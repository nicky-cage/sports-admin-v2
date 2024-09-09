--场馆维护日志
drop table if exists game_maintains_logs;
create table if not exists game_maintains_logs (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `game_id` int unsigned NOT NULL DEFAULT '0' COMMENT '游戏场馆编号',
  `game_state` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '游戏状态, 0:正常;1:禁用',
  `time_start` int unsigned NOT NULL DEFAULT '0' COMMENT '开始时间',
  `time_end` int unsigned NOT NULL DEFAULT '0' COMMENT '结束时间',
  `from_state` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '从状态, 0:正常;1:维护',
  `to_state` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '到状态, 0:正常;1:维护',
  `remark` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '备注',
  `admin_id` int unsigned NOT NULL DEFAULT '0' COMMENT '管理员编号',
  `admin_name` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '管理员名称',
  `created` int unsigned NOT NULL DEFAULT '0' COMMENT '维护时间',
  `game_name` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '游戏中文名',
  `game_ename` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '游戏英文名',
  PRIMARY KEY (`id`),
  KEY `game_id` (`game_id`)
);