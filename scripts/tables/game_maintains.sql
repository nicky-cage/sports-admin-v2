-- 场馆维护记录
drop table if exists game_venue_maintains;
create table if not exists game_venue_maintains (
	id int unsigned not null auto_increment,
	venue_id int unsigned not null default 0 comment '场馆编号',
    venue_state tinyint unsigned not null default 0 comment '场馆状态, 0:正常;1:禁用',
    time_start int unsigned not null default 0 comment '开始时间',
    state tinyint unsigned not null default 0 comment '状态, 0:正常;1:维护',
    time_end int unsigned not null default 0 comment '结束时间',
    remark varchar(200) not null default 0 comment '备注',
    admin_id int unsigned not null default 0 comment '管理员编号',
    admin_name varchar(20) not null default '' comment '管理员名称',
    index(venue_id),
	primary key(id)
);

alter table games add column is_online tinyint unsigned not null default 0 comment '是否在线, 0:否;1:是;';
alter table games add column maintain_start_at int unsigned not null default 0 comment '维护开始时间';
alter table games add column maintain_end_at int unsigned not null default 0 comment '维护结束时间';
alter table games add column admin_id int unsigned not null default 0 comment '维护用户编号';
alter table games add column admin_name varchar(50) not null default '' comment '维护人员名称';
alter table games add column remark varchar(200) not null default '' comment '备注';

-- 场馆维护
create or replace view game_venue_maintains_v as 
    select m.id, m.venue_id, m.venue_state, m.time_start, m.time_end, m.state, m.remark, m.admin_id, m.admin_name,
        v.name as venue_name, v.ename as venue_ename 
    from game_venue_maintains as m 
    inner join game_venues as v on m.venue_id = v.id;
   
-- 初始化现数据
-- 此句只能执行一次
insert into game_venue_maintains (venue_id, venue_state) select id, maintain from game_venues;

-- 场馆维护日志
drop table if exists game_venue_maintain_logs;
create table if not exists game_venue_maintain_logs (
	id int unsigned not null auto_increment,
	venue_id int unsigned not null default 0 comment '场馆编号',
    venue_state tinyint unsigned not null default 0 comment '场馆状态, 0:正常;1:禁用',
    time_start int unsigned not null default 0 comment '开始时间',
    time_end int unsigned not null default 0 comment '结束时间',
    state tinyint unsigned not null default 0 comment '到状态, 0:正常;1:维护',
    remark varchar(200) not null default 0 comment '备注',
    admin_id int unsigned not null default 0 comment '管理员编号',
    admin_name varchar(20) not null default '' comment '管理员名称',
    created int unsigned not null default 0 comment '维护时间',
    index(venue_id),
	primary key(id)
);

-- 场馆维护日志
create or replace view game_venue_maintain_logs_v as 
    select m.id, m.venue_id, m.venue_state, m.time_start, m.time_end, m.state, m.remark, m.admin_id, m.admin_name, m.created,
        v.name as venue_name, v.ename as venue_ename 
    from game_venue_maintain_logs as m 
    inner join game_venues as v on m.venue_id = v.id;
