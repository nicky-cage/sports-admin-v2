ALTER TABLE notices MODIFY COLUMN `type` tinyint UNSIGNED NOT NULL DEFAULT 1 COMMENT '类型, 1:普通;2:特殊;3:财务;';
ALTER TABLE notices MODIFY COLUMN `is_online` tinyint UNSIGNED NULL DEFAULT 1 COMMENT '是否在线, 1:启用;2:停用;';
ALTER TABLE notices ADD COLUMN `platform_types` varchar(100) NOT NULL DEFAULT '' COMMENT '平台类型, 1:全站;2:体育;3:web;4:h5;';
ALTER TABLE notices ADD COLUMN `jump_url` varchar(200) NOT NULL DEFAULT '' COMMENT '跳转链接';
ALTER TABLE notices ADD COLUMN `vip_ids` varchar(100) NOT NULL DEFAULT '' COMMENT 'VIP等级';
ALTER TABLE notices ADD COLUMN `img_url` varchar(100) NOT NULL DEFAULT '' COMMENT '图标地址';
DROP COLUMN `platform`,
DROP COLUMN `vip`;
CHANGE COLUMN `is_online` `state` tinyint(3) UNSIGNED NULL DEFAULT 1 COMMENT ' 1:启用;2:停用;' AFTER `end_at`;

alter table notices drop column start_at;
alter table notices add column start_at int unsigned not null default 0 comment '开始时间';

alter table notices drop column end_at;
alter table notices add column end_at int unsigned not null default 0 comment '结束时间';
