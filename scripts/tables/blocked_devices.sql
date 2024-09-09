create table if not exists blocked_devices ( 
    id int unsigned not null auto_increment,
	device_no varchar(50) not null default '' comment '设备编号',
	user_id int unsigned not null default 0 comment '用户编号',
	user_name varchar(20) not null default '' comment '用户名称',
	remark varchar(200) not null default '' comment '备注',
	created int unsigned not null default 0 comment '添加时间',
	updated int unsigned not null default 0 comment '修改时间',
	admin_id int unsigned not null default 0 comment '修改人员编号',
	admin_name varchar(20) not null default '' comment '修改人员名称',
	disabled_all tinyint not null default 0 comment '禁止所有登录',
    index(device_no),
	primary key(id, created)
) partition by range(created) (
	partition p202009 values less than (1601510400), -- 2020-10-01 之前
	partition p202010 values less than (1604188800), -- 2020年10月
	partition p202011 values less than (1606780800), -- 2020年11月
	partition p202012 values less than (1609459200), -- 2020年12月
	partition p202101 values less than (1612137600), -- 2021年01月
	partition p202102 values less than (1614556800), -- 2021年02月
	partition p202103 values less than (1617235200), -- 2021年03月
	partition p202104 values less than (1619827200), -- 2021年04月
	partition p202105 values less than (1622505600), -- 2021年05月
	partition p202106 values less than (1625097600), -- 2021年06月
	partition p202107 values less than (1627776000), -- 2021年07月
	partition p202108 values less than (1630454400), -- 2021年08月
	partition p202109 values less than (1633046400), -- 2021年09月
	partition p202110 values less than (1635724800), -- 2021年10月
	partition p202111 values less than (1638316800), -- 2021年11月
	partition p202112 values less than (1640995200), -- 2021年12月
	partition p202201 values less than MAXVALUE -- 2022年之后
);

insert into blocked_devices (device_no, user_id, user_name, remark, created, updated, admin_id, admin_name, disabled_all) values 
('abc:123:456:789:111', 1, 'tempname', '搙羊毛用户00', unix_timestamp(), unix_timestamp() + 100, 1, 'admin', 1),
('abc:123:456:789:112', 1, 'tempname', '搙羊毛用户01', unix_timestamp(), unix_timestamp() + 200, 1, 'admin', 1),
('abc:123:456:789:113', 1, 'tempname', '搙羊毛用户02', unix_timestamp(), unix_timestamp() + 300, 1, 'admin', 0),
('abc:123:456:789:114', 1, 'tempname', '搙羊毛用户03', unix_timestamp(), unix_timestamp() + 400, 1, 'admin', 1),
('abc:123:456:789:115', 1, 'tempname', '搙羊毛用户04', unix_timestamp(), unix_timestamp() + 500, 1, 'admin', 0),
('abc:123:456:789:116', 1, 'tempname', '搙羊毛用户05', unix_timestamp(), unix_timestamp() + 600, 1, 'admin', 1),
('abc:123:456:789:117', 1, 'tempname', '搙羊毛用户06', unix_timestamp(), unix_timestamp() + 700, 1, 'admin', 0),
('abc:123:456:789:118', 1, 'tempname', '搙羊毛用户07', unix_timestamp(), unix_timestamp() + 800, 1, 'admin', 0),
('abc:123:456:789:119', 1, 'tempname', '搙羊毛用户08', unix_timestamp(), unix_timestamp() + 900, 1, 'admin', 1),
('abc:123:456:789:121', 1, 'tempname', '搙羊毛用户09', unix_timestamp(), unix_timestamp() + 110, 1, 'admin', 1),
('abc:123:456:789:122', 1, 'tempname', '搙羊毛用户00', unix_timestamp(), unix_timestamp() + 120, 1, 'admin', 0),
('abc:123:456:789:123', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 130, 1, 'admin', 1),
('abc:123:456:789:124', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 140, 1, 'admin', 0),
('abc:123:456:789:125', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 150, 1, 'admin', 0),
('abc:123:456:789:126', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 160, 1, 'admin', 0),
('abc:123:456:789:127', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 170, 1, 'admin', 0),
('abc:123:456:789:128', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 180, 1, 'admin', 0),
('abc:123:456:789:129', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 190, 1, 'admin', 1),
('abc:123:456:789:131', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 101, 1, 'admin', 1),
('abc:123:456:789:131', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 102, 1, 'admin', 1),
('abc:123:456:789:132', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 103, 1, 'admin', 0),
('abc:123:456:789:133', 1, 'tempname', '搙羊毛用户10', unix_timestamp(), unix_timestamp() + 104, 1, 'admin', 0),
('abc:123:456:789:134', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 105, 1, 'admin', 0),
('abc:123:456:789:135', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 106, 1, 'admin', 0),
('abc:123:456:789:136', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 107, 1, 'admin', 0),
('abc:123:456:789:137', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 108, 1, 'admin', 0),
('abc:123:456:789:518', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 109, 1, 'admin', 0),
('abc:123:456:789:611', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 100, 1, 'admin', 0),
('abc:123:456:789:612', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 901, 1, 'admin', 1),
('abc:123:456:789:613', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 902, 1, 'admin', 1),
('abc:123:456:789:614', 1, 'tempname', '搙羊毛用户20', unix_timestamp(), unix_timestamp() + 903, 1, 'admin', 1),
('abc:123:456:789:615', 1, 'tempname', '搙羊毛用户30', unix_timestamp(), unix_timestamp() + 904, 1, 'admin', 1),
('abc:123:456:789:616', 1, 'tempname', '搙羊毛用户30', unix_timestamp(), unix_timestamp() + 905, 1, 'admin', 1),
('abc:123:456:789:617', 1, 'tempname', '搙羊毛用户30', unix_timestamp(), unix_timestamp() + 906, 1, 'admin', 1),
('abc:123:456:789:618', 1, 'tempname', '搙羊毛用户30', unix_timestamp(), unix_timestamp() + 907, 1, 'admin', 1),
('abc:123:456:789:811', 1, 'tempname', '搙羊毛用户30', unix_timestamp(), unix_timestamp() + 908, 1, 'admin', 1),
('abc:123:456:789:812', 1, 'tempname', '搙羊毛用户30', unix_timestamp(), unix_timestamp() + 909, 1, 'admin', 1),
('abc:123:456:789:813', 1, 'tempname', '搙羊毛用户30', unix_timestamp(), unix_timestamp() + 910, 1, 'admin', 1),
('abc:123:456:789:816', 1, 'tempname', '搙羊毛用户30', unix_timestamp(), unix_timestamp() + 912, 1, 'admin', 1),
('abc:123:456:789:817', 1, 'tempname', '搙羊毛用户60', unix_timestamp(), unix_timestamp() + 913, 1, 'admin', 1),
('abc:123:456:789:818', 1, 'tempname', '搙羊毛用户60', unix_timestamp(), unix_timestamp() + 801, 1, 'admin', 1),
('abc:123:456:789:819', 1, 'tempname', '搙羊毛用户60', unix_timestamp(), unix_timestamp() + 802, 1, 'admin', 1),
('abc:123:456:789:891', 1, 'tempname', '搙羊毛用户60', unix_timestamp(), unix_timestamp() + 803, 1, 'admin', 1),
('abc:123:456:789:892', 1, 'tempname', '搙羊毛用户60', unix_timestamp(), unix_timestamp() + 804, 1, 'admin', 1),
('abc:123:456:789:912', 1, 'tempname', '搙羊毛用户60', unix_timestamp(), unix_timestamp() + 805, 1, 'admin', 1),
('abc:123:456:789:913', 1, 'tempname', '搙羊毛用户60', unix_timestamp(), unix_timestamp() + 806, 1, 'admin', 1);

