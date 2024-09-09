-- 维护
drop table if exists site_maintains;
create table if not exists site_maintains (
	id int unsigned not null auto_increment,
	platform_id int unsigned not null default 0 comment '平台编号',
    time_start int unsigned not null default 0 comment '开始时间',
    time_end int unsigned not null default 0 comment '结束时间',
    state tinyint unsigned not null default 0 comment '状态, 0:正常;1:维护',
    remark varchar(200) not null default 0 comment '备注',
    admin_id int unsigned not null default 0 comment '管理员编号',
    admin_name varchar(20) not null default '' comment '管理员名称',
    index(platform_id),
	primary key(id)
);

-- 此句只能执行一次
insert into site_maintains (platform_id) values 
(1),
(2),
(3),
(4),
(5),
(6);

-- 维护日志
drop table if exists site_maintain_logs;
create table if not exists site_maintain_logs (
	id int unsigned not null auto_increment,
	platform_id int unsigned not null default 0 comment '平台编号',
    time_start int unsigned not null default 0 comment '开始时间',
    time_end int unsigned not null default 0 comment '结束时间',
    from_state tinyint unsigned not null default 0 comment '从状态, 0:正常;1:维护',
    to_state tinyint unsigned not null default 0 comment '到状态, 0:正常;1:维护',
    remark varchar(200) not null default 0 comment '备注',
    admin_id int unsigned not null default 0 comment '管理员编号',
    admin_name varchar(20) not null default '' comment '管理员名称',
    created int unsigned not null default 0 comment '维护时间',
    index(platform_id),
	primary key(id)
);

-- 修改表结构
alter table site_maintain_logs 
    change to_state state tinyint unsigned not null default 0 comment '状态, 1:正常;0:维护;',
    drop column from_state;
