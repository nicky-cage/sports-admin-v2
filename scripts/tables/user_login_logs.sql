drop table if exists user_login_logs;
create table if not exists user_login_logs (
	id int unsigned not null auto_increment,
	user_id int unsigned not null comment '用户编号',
	user_name varchar(50) not null default '' comment '用户名称',
	device_no varchar(50) not null default '' comment '设备编号',
	login_ip varchar(30) not null default '' comment '登录号ip',
	last_ip varchar(30) not null default '' comment '上次登录IP',
	created int unsigned not null default 0 comment '登录时间',
	login_area varchar(30) not null default '' comment '登录地区',
	log_type tinyint not null default 0 comment '类型',
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

insert into user_login_logs (user_id, user_name, device_no, login_ip, last_ip, created, login_area, log_type) values 
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp(), '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 30, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 50, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 60, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 70, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 80, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 90, '本地 局域网', 1);

alter table user_login_logs change log_type login_type tinyint not null default 0 comment '登录类型';
alter table user_login_logs add column remark varchar(200) not null default '' comment '备注';

insert into user_login_logs (user_id, user_name, device_no, login_ip, last_ip, created, login_area, login_type) values 
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp(), '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 30, '本地 局域网', 0),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 2),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 50, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 60, '本地 局域网', 2),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 60, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 60, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 60, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 60, '本地 局域网', 2),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 70, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 80, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 90, '本地 局域网', 1);

insert into user_login_logs (user_id, user_name, device_no, login_ip, last_ip, created, login_area, login_type) values 
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp(), '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 30, '本地 局域网', 0),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 2),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 2),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 2),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 0),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 40, '本地 局域网', 2),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 50, '本地 局域网', 1),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 60, '本地 局域网', 2),
(1, 'tempname', '01x:392k:d12h:1239:3dar', '127.0.0.', '127.0.0.1', unix_timestamp() + 60, '本地 局域网', 1);

-- 创建视图
create or replace view user_login_logs_v as 
    select l.id, l.user_id, l.user_name, l.device_no, l.last_ip, l.login_ip, l.created, l.login_area, l.remark, l.login_type, u.top_id, u.top_name
    from user_login_logs as l inner join users as u
    on l.user_id = u.id;

