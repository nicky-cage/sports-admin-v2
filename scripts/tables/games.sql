ALTER TABLE games
    MODIFY COLUMN name varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT '' COMMENT '游戏场馆名称' AFTER id,
    ADD COLUMN ename varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NULL DEFAULT '' COMMENT '游戏场馆英文名称' AFTER name;

-- 修改英文名称
update games set ename = code;

-- 插入初始数据
insert into game_electrons (game_type, game_code, web_code, h5_code, cn_name, en_name, created, updated) values 
(1, 'SPORT', 'SPORT_WEB', 'SPORT_H5', '贪吃蛇大战', 'Snake', unix_timestamp(), unix_timestamp()),
(1, 'SPORT', 'SPORT_WEB', 'SPORT_H5', '三国志曹参传', 'Snake', unix_timestamp(), unix_timestamp());


-- 游戏场馆表
drop table if exists game_venues;
create table if not exists game_venues select * from games;
-- 设定主键
alter table game_venues
    change id id int unsigned not null auto_increment primary key;

-- 删除原来游戏表
drop table if exists games; 

-- 新建游戏表
create table if not exists games select * from game_electrons;

-- 设定主键
alter table games change id id int unsigned not null auto_increment primary key;

-- 增加支持平台字段
alter table games add column platform_types char(16) not null default '' comment '支持平台';

-- 增加展示类型
alter table games add column show_types char(16) not null default '' comment '展示类型';

-- 添加游戏编号
alter table games add column venue_id int unsigned not null default 0 comment '场馆编号';


-- 删除某些字段
alter table games drop column display_type;
alter table games drop column platform_type;
alter table games drop column show_types;

-- 修改字段类型
alter table games add column display_types char(16) not null default '' comment '展示类型';

-- 修改状态
alter table games add column state tinyint unsigned not null default 0 comment '状态, 0:禁用;1:正常;';

alter table games drop column is_new;
alter table games drop column is_hot;
alter table games drop column is_online;

-- 修改场馆类型
alter table games change game_type venue_type tinyint unsigned not null default 0 comment '场馆类型';

-- 修改场馆类型
alter table game_venues change `type` venue_type tinyint unsigned not null default 0 comment '场馆类型';
