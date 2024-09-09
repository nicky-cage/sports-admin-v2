CREATE TABLE `user_levels` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` char(8) NOT NULL DEFAULT '' COMMENT '等级名称',
  `digit` tinyint NOT NULL DEFAULT '0' COMMENT '等级',
  `upgrade_deposit` decimal(16,3) NOT NULL DEFAULT '0.000' COMMENT '升级存款',
  `hold_stream` decimal(16,3) NOT NULL DEFAULT '0.000' COMMENT '保级流水',
  `upgrade_stream` decimal(16,3) NOT NULL DEFAULT '0.000' COMMENT '升级流水',
  `upgrade_bonus` decimal(16,3) NOT NULL DEFAULT '0.000' COMMENT '升级红利',
  `birth_bonus` decimal(16,3) NOT NULL DEFAULT '0.000' COMMENT '生日红利',
  `month_bonus` decimal(16,3) NOT NULL DEFAULT '0.000' COMMENT '月度红利',
  `day_withdraw_count` tinyint unsigned DEFAULT '1' COMMENT '每日提款限制',
  `day_withdraw_total` decimal(16,3) NOT NULL DEFAULT '100.000' COMMENT '每日提款总额',
  `created` int unsigned DEFAULT '0' COMMENT '添加时间',
  `updated` int unsigned DEFAULT '0' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `digit` (`digit`),
  KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO user_levels (name,digit,upgrade_deposit,hold_stream,upgrade_stream,upgrade_bonus,birth_bonus,month_bonus,day_withdraw_count,day_withdraw_total,created,updated) VALUES 
('VIP0',1,1000.000,1000.000,1000.000,2.000,10.000,10.000,2,1000.000,1590720202,1590720202)
,('VIP1',2,2000.000,2000.000,2000.000,2.000,20.000,20.000,4,2000.000,1590720202,1590720202)
,('VIP2',3,3000.000,3000.000,3000.000,2.000,30.000,30.000,6,3000.000,1590720202,1590720202)
,('VIP3',4,4000.000,4000.000,4000.000,5.000,40.000,40.000,8,4000.000,1590720202,1590720202)
,('VIP4',5,5000.000,5000.000,5000.000,5.000,50.000,50.000,10,5000.000,1590720202,1590720202)
,('VIP5',6,6000.000,6000.000,6000.000,5.000,60.000,60.000,12,6000.000,1590720202,1590720202)
,('VIP6',7,7000.000,7000.000,7000.000,10.000,70.000,70.000,14,7000.000,1590720202,1590720202)
,('VIP7',8,8000.000,8000.000,8000.000,10.000,80.000,80.000,16,8000.000,1590720202,1590720202)
,('VIP8',9,9000.000,9000.000,9000.000,10.000,90.000,90.000,18,9000.000,1590720202,1590720202)
,('VIP9',10,10000.000,10000.000,10000.000,20.000,100.000,100.000,20,10000.000,1590720202,1590720202)
,('VIP10',11,11000.000,11000.000,11000.000,20.000,110.000,110.000,22,11000.000,1590720202,1590720202);

/** 添加唯一索引 **/
alter table user_levels add unique index(name);
alter table user_levels add unique index(digit);
